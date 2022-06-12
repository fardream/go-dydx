package dydx

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const eip712DomainName = "EIP712Domain"

// SignTypedData is an interface to sign apitypes.TypedData on ethereum network.
// For a reference implementation, see:
// https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.18/signer/core#SignerAPI.SignTypedData
type SignTypedData interface {
	EthSignTypedData(typedData apitypes.TypedData) ([]byte, error)
}

// ecdsaPrivateKeySigner implements the SignTypedData interface for ecdsa.PrivateKey
type ecdsaPrivateKeySigner ecdsa.PrivateKey

// NewEcdsaPrivateKeySigner converts an *ecdsa.PrivateKey into a signer to EthSignTypedData
func NewEcdsaPrivateKeySigner(key *ecdsa.PrivateKey) *ecdsaPrivateKeySigner {
	return (*ecdsaPrivateKeySigner)(key)
}

// prepareTypedDataForSign does part of the process in geth's SignerApi.SignTypedData.
// https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.18/signer/core#SignerAPI.SignTypedData
// See the source code of tha method for the details.
// The function returns the raw bytes ([]byte) that can be fed into Wallet.SignData
func prepareTypedDataForSign(typedData apitypes.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct(eip712DomainName, typedData.Domain.Map())
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
func (key *ecdsaPrivateKeySigner) EthSignTypedData(typedData apitypes.TypedData) ([]byte, error) {
	rawData, err := prepareTypedDataForSign(typedData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert typed data into raw bytes: %w", err)
	}
	sig, err := crypto.Sign(crypto.Keccak256(rawData), (*ecdsa.PrivateKey)(key))
	if err != nil {
		return nil, fmt.Errorf("failed to sign the raw bytes for typed data: %w", err)
	}
	// Legacy Signing
	sig[64] += 27

	return sig, nil
}
