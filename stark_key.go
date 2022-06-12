package dydx

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

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

func stripLeadingZeros(s string) string {
	for {
		s1 := strings.TrimPrefix(s, "0")
		if s1 == s {
			return s1
		} else {
			s = s1
		}
	}
}

func NewStarkKey(publicKey, publicKeyYCoordinate, privateKey string) *StarkKey {
	return &StarkKey{
		PublicKey:            stripLeadingZeros(publicKey),
		PrivateKey:           stripLeadingZeros(privateKey),
		PublicKeyYCoordinate: stripLeadingZeros(publicKeyYCoordinate),
	}
}

// String prints out the key, for cobra cli
func (c *StarkKey) String() string {
	return fmt.Sprintf("public key: %s - public key y: %s - private key: (redacted)", c.PublicKey, c.PublicKeyYCoordinate)
}

// Set reads in the file, for cobra cli
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

// Type is for cobra cli
func (c *StarkKey) Type() string {
	return "stark-key-map-file"
}
