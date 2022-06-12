package dydx

import (
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

func getOnboardingTypedData(isMainnet bool, action string) apitypes.TypedData {
	result := apitypes.TypedData{
		Types: apitypes.Types{
			eip712DomainName: []apitypes.Type{
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

// DeriveStarkKey gets the default deterministic stark key.
//
// - sign the typed data for key derivation.
//
// - append 0 to the signature.
//
// - convert the siganture to a big int.
//
// - in python implementation, it looks like the resulting big int is crypto.keccak
// hashed as an uint256, and uint256 is 32 bytes. however, in reality, it takes
// the whole 66 bytes and hashes it.
//
// - crypto.keccak hash the signature bytes
//
// - convert the resulted signature bytes into a big int.
//
// - right shift the big int by 5 -  this is our private key.
//
// - `starkex.PrivateKeyToEcPointOnStarkCurv` to get the x and y big ints.
//
// - convert private key, x, y into hex encoded strings (without the 0x).
//
// Function requires a signer to sign typed data
func DeriveStarkKey(signer SignTypedData, isMainnet bool) (*StarkKey, error) {
	msg := getOnboardingTypedData(isMainnet, keyDerivationAction)

	signature, err := signer.EthSignTypedData(msg)
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
