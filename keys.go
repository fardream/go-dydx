package dydx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
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

func (c *ApiKey) String() string {
	return fmt.Sprintf("key: %s - passphrase: %s - secret: (redacted)", c.Key, c.Passphrase)
}

func (c *ApiKey) Set(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m, err := ParseApiKeyMap(data)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("only one keys is allowed: %s", data)
	}
	for _, v := range m {
		*c = *v
	}
	return nil
}

func (c *ApiKey) Type() string {
	return "api-key-map-file"
}

// StarkKey is the private key on the Stark L2.
// Below fields are from the browser cache but unused.
// - WalletAddress        string `json:"walletAddress"`
// - LegacySigning        bool   `json:"legacySigning"`
// - WalletType           string `json:"walletType"`
type StarkKey struct {
	PublicKey            string `json:"publicKey"`
	PublicKeyYCoordinate string `json:"publicKeyYCoordinate"`
	PrivateKey           string `json:"privateKey"`
}

func ParseStarkKeyMap(input []byte) (map[string]*StarkKey, error) {
	result := make(map[string]*StarkKey)
	if err := json.Unmarshal(input, &result); err != nil {
		return nil, fmt.Errorf("cannot parse json: %w", err)
	}
	return result, nil
}

func NewStarkKey(publicKey, publicKeyYCoordinate, privateKey string) *StarkKey {
	return &StarkKey{
		PublicKey:            publicKey,
		PrivateKey:           privateKey,
		PublicKeyYCoordinate: publicKeyYCoordinate,
	}
}

func (c *StarkKey) String() string {
	return fmt.Sprintf("public key: %s - public key y: %s - private key: (redacted)", c.PublicKey, c.PublicKeyYCoordinate)
}

func (c *StarkKey) Set(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m, err := ParseStarkKeyMap(data)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("only one keys is allowed: %s", data)
	}
	for _, v := range m {
		*c = *v
	}
	return nil
}

func (c *StarkKey) Type() string {
	return "stark-key-map-file"
}
