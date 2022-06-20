package dydx

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/google/uuid"

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
// Derived from the python implementation of official dydx client:
// https://github.com/dydxprotocol/dydx-v3-python/blob/914fc66e542d82080702e03f6ad078ca2901bb46/dydx3/modules/onboarding.py#L116-L145
//
// - sign the typed data for key derivation.
//
// - append 0 to the signature.
//
// - crypto.keccak hash the signature bytes.
// in python implementation, it convert the siganture to a big int.
// and the resulting big int is crypto.keccak
// hashed as an uint256 - note uint256 is 32 bytes. however, in reality, it takes
// the whole 66 bytes and hashes it.
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
		return nil, fmt.Errorf("failed to sign typed data %#v: %w", msg, err)
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

// RecoverDefaultApiKeyCredentials recovers the default credentials.
//
// Implementation is a carbon-copy of the code in python version of official dydx client:
// https://github.com/dydxprotocol/dydx-v3-python/blob/914fc66e542d82080702e03f6ad078ca2901bb46/dydx3/modules/onboarding.py#L147-L184
func RecoverDefaultApiKeyCredentials(signer SignTypedData, isMainnet bool) (*ApiKey, error) {
	msg := getOnboardingTypedData(isMainnet, onboardingAction)
	signature_raw, err := signer.EthSignTypedData(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign typed data %#v: %w", msg, err)
	}
	// here must append an extra byte of 0 (SIGNATURE_TYPE_NO_PREPEND)
	signature_raw = append(signature_raw, 0)

	// Python will convert signature into hex, the convert the values into big int, and then hash that bytes of the big int.
	// we simply use the raw bytes in signature_raw.
	//
	// signature := hexutil.Encode(signature_raw)

	// In python implementation, signature_raw is hex encoded with 0x prefix, the 64 characters without
	// the 0x prefix is taken to make a 32 byte big int. Then the whole big int is hashed as uint256, which
	// is fortunately 32 byte long.
	// ```
	// r_hex := signature[2:66]
	// r_int, ok := new(big.Int).SetString(r_hex, 16)
	// if !ok {
	// 	return nil, fmt.Errorf("%s is not a proper big int hex", r_hex)
	// }
	// hashed_r_bytes := crypto.Keccak256(r_int.Bytes())
	// ```
	hashed_r_bytes := crypto.Keccak256(signature_raw[:32])

	secret_bytes := hashed_r_bytes[:30]

	// Similarly to what happend to hashed_r_bytes, here python converts everything into big int and hash that,
	// here we simply take the 32-64 byte of the signature.
	//
	// ```
	// s_hex := signature[66:130]
	// s_int, ok := new(big.Int).SetString(s_hex, 16)
	// if !ok {
	// 	return nil, fmt.Errorf("%s is not a proper big int hex", s_hex)
	// }
	// hashed_s_bytes := crypto.Keccak256(s_int.Bytes())
	// ```
	hashed_s_bytes := crypto.Keccak256(signature_raw[32:64])
	key_bytes := hashed_s_bytes[:16]
	passphrase_bytes := hashed_s_bytes[16:31]

	key_uuid, err := uuid.FromBytes(key_bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain UUID from bytes: %w", err)
	}

	return NewApiKey(key_uuid.String(), base64.URLEncoding.EncodeToString(passphrase_bytes), base64.URLEncoding.EncodeToString(secret_bytes)), nil
}
