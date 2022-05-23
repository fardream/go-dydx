package dydx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ApiKey is the format from browser's local storage.
type ApiKey struct {
	WalletAddress string `json:"walletAddress"`
	Secret        string `json:"secret"`
	Key           string `json:"key"`
	Passphrase    string `json:"passphrase"`
	LegacySigning bool   `json:"legacySigning"`
	WalletType    string `json:"walletType"`
}

func ParseApiKeyMap(input []byte) (map[string]*ApiKey, error) {
	result := make(map[string]*ApiKey)
	if err := json.Unmarshal(input, &result); err != nil {
		return nil, fmt.Errorf("cannot parse json: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("there is no data in the input")
	}
	return result, nil
}

func NewApiKey(ethAddress, key, passphrase, secret string) *ApiKey {
	return &ApiKey{
		WalletAddress: ethAddress,
		Key:           key,
		Passphrase:    passphrase,
		Secret:        secret,
		LegacySigning: false,
	}
}

func (a *ApiKey) Sign(requestPath, method, isoTimestamp string, body []byte) string {
	message := fmt.Sprintf("%s%s%s%s", isoTimestamp, method, requestPath, body)
	secret, _ := base64.URLEncoding.DecodeString(a.Secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// StarkKey
type StarkKey struct {
	WalletAddress        string `json:"walletAddress"`
	PublicKey            string `json:"publicKey"`
	PublicKeyYCoordinate string `json:"publicKeyYCoordinate"`
	PrivateKey           string `json:"privateKey"`
	LegacySigning        bool   `json:"legacySigning"`
	WalletType           string `json:"walletType"`
}

func ParseStarkKeyMap(input []byte) (map[string]*StarkKey, error) {
	result := make(map[string]*StarkKey)
	if err := json.Unmarshal(input, &result); err != nil {
		return nil, fmt.Errorf("cannot parse json: %w", err)
	}
	return result, nil
}

func NewStarkKey(ethAddress, publicKey, publicKeyYCoordinate, privateKey string) *StarkKey {
	return &StarkKey{
		WalletAddress:        ethAddress,
		PublicKey:            publicKey,
		PrivateKey:           privateKey,
		PublicKeyYCoordinate: publicKeyYCoordinate,
		LegacySigning:        false,
	}
}
