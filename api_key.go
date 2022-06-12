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
// Below are from browser but not necessary
//
// - WalletAddress string `json:"walletAddress"`
// - LegacySigning bool   `json:"legacySigning"`
// - WalletType    string `json:"walletType"`
type ApiKey struct {
	Secret     string `json:"secret"`
	Key        string `json:"key"`
	Passphrase string `json:"passphrase"`
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

func NewApiKey(key, passphrase, secret string) *ApiKey {
	return &ApiKey{
		Key:        key,
		Passphrase: passphrase,
		Secret:     secret,
	}
}

// Sign a request for dydx
func (a *ApiKey) Sign(requestPath, method, isoTimestamp string, body []byte) string {
	message := fmt.Sprintf("%s%s%s%s", isoTimestamp, method, requestPath, body)
	secret, _ := base64.URLEncoding.DecodeString(a.Secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// String prints out the key, for cobra cli.
func (c *ApiKey) String() string {
	return fmt.Sprintf("key: %s - passphrase: %s - secret: (redacted)", c.Key, c.Passphrase)
}

// Set reads the file containing the ApiKey, for cobra cli.
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

// Type is for cobra cli.
func (c *ApiKey) Type() string {
	return "api-key-map-file"
}
