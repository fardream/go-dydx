package dydx

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/fardream/go-dydx/starkex"
)

const (
	onboardingAction    = "dYdX Onboarding"
	keyDerivationAction = "dYdX STARK Key"
)

const eip712StructName = "dYdX"

// ByteSigner can sign raw bytes according to Ethereum convention.
type ByteSigner interface {
	// EthSignRawData signs the raw data
	EthSignRawData(data []byte) ([]byte, error)
}

// ecdsaPrivateKeySigner implements the ByteSigner interface for ecdsa.PrivateKey
type ecdsaPrivateKeySigner ecdsa.PrivateKey

func NewEcdsaPrivateKeySigner(key *ecdsa.PrivateKey) *ecdsaPrivateKeySigner {
	return (*ecdsaPrivateKeySigner)(key)
}

func (s *ecdsaPrivateKeySigner) EthSignRawData(data []byte) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(data), (*ecdsa.PrivateKey)(s))
}

func getOnboardingTypedData(isMainnet bool, action string) apitypes.TypedData {
	result := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
		},
		PrimaryType: eip712StructName,
		Domain: apitypes.TypedDataDomain{
			Name:    "dYdX",
			Version: "1.0",
		},
		Message: map[string]any{
			"action": action,
		},
	}

	result.PrimaryType = eip712StructName

	if isMainnet {
		result.Types[eip712StructName] = []apitypes.Type{
			{Name: "action", Type: "string"},
			{Name: "onlySignOn", Type: "string"},
		}
		result.Domain.ChainId = math.NewHexOrDecimal256(NetworkIdMainnet)
		result.Message["onlySignOn"] = "https://trade.dydx.exchange"
	} else {
		result.Types[eip712StructName] = []apitypes.Type{
			{Name: "action", Type: "string"},
		}
		result.Domain.ChainId = math.NewHexOrDecimal256(NetworkIdRopsten)
	}

	return result
}

// prepareTypedDataForSign does part of the process in geth's SignerApi.SignTypedData.
// https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.18/signer/core#SignerAPI.SignTypedData
// See the source code of tha method for the details.
// The function returns the raw bytes ([]byte) that can be fed into Wallet.SignData
func prepareTypedDataForSign(typedData apitypes.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to generate domain separator: %w", err)
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to generate typed data hash: %w", err)
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))

	return rawData, nil
}

// SignTypedData's process is coming from SignerApi.SignTypedData.
// https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.18/signer/core#SignerAPI.SignTypedData
// See the source code of tha method for the details.
// The function returns the signed signature of the typedData.
// After the process of `PrepareTypedDataForSign`, the generated []byte can be signed by Wallet.SignData method.
// Wallet is defined here: https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.18/accounts#Wallet.
// However, we only need SignData method without even the mimeType or account.
// SignData method in the keystoreWallet simply gets the hash and then sign the hash with the private key.
func signTypeData(typedData apitypes.TypedData, signer ByteSigner) ([]byte, error) {
	rawData, err := prepareTypedDataForSign(typedData)
	if err != nil {
		return nil, err
	}
	sig, err := signer.EthSignRawData(rawData)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

// DeriveStarkKey gets the default deterministic stark key.
//
// - sign the typed data for key derivation.
// - append 0 to the signature.
// - convert the siganture to a big int.
// - in python implementation, it looks like the resulting big int is crypto.keccak
// hashed as an uint256, and uint256 is 32 bytes. however, in reality, it takes
// the whole 66 bytes and hashes it.
// - crypto.keccak hash the signature bytes
// - convert the resulted signature bytes into a big int.
// - right shift the big int by 5 -  this is our private key.
// - `starkex.PrivateKeyToEcPointOnStarkCurv` to get the x and y big ints.
// - convert private key, x, y into hex encoded strings (without the 0x).
//
// Function requires a signer, which can sign raw bytes.
func DeriveStarkKey(signer ByteSigner, isMainnet bool) (*StarkKey, error) {
	msg := getOnboardingTypedData(isMainnet, keyDerivationAction)

	signature, err := signTypeData(msg, signer)
	if err != nil {
		return nil, err
	}
	// here must append an extra byte of 0 (SIGNATURE_TYPE_NO_PREPEND)
	signature = append(signature, byte(0))

	// hash with all the bytes
	hashedSignature := hexutil.Encode(crypto.Keccak256(signature))

	privateKey, ok := new(big.Int).SetString(hashedSignature, 0)

	if !ok {
		return nil, fmt.Errorf("%s is not a valid int for private key", hashedSignature)
	}

	privateKey = privateKey.Rsh(privateKey, 5)

	x, y, err := starkex.PrivateKeyToEcPointOnStarkCurv(privateKey)
	if err != nil {
		return nil, err
	}

	return NewStarkKey(hex.EncodeToString(x.Bytes()), hex.EncodeToString(y.Bytes()), hex.EncodeToString(privateKey.Bytes())), nil
}
