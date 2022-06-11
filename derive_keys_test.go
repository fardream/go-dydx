package dydx_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"

	"github.com/fardream/go-dydx"
)

func TestDeriveStarkKey(t *testing.T) {
	config := spew.NewDefaultConfig()
	config.DisableMethods = true
	private_key_hex := "2b61fc77bbda78d7075d91ae6a58423df73f71ebe4aa36e3b587dc8e5fe1396f"
	private_key, err := crypto.HexToECDSA(private_key_hex)
	if err != nil {
		t.Fatalf("%s is not marshalled into private key: %#v", private_key_hex, err)
	}
	stark_key_python_testnet := dydx.NewStarkKey(
		"05ca235e41160c7c2a4695a98eecf7cb0ba6f7d26c4f899c192ca91912995115",
		"e07e9de2353c45e92b98e6c566fce8686f7fb620031aa7a79a5983daf22483",
		"03c6ce687a484ac1c50a48498092832957a3154c7c13237bc10df6965472e009",
	)

	stark_key_testnet, err := dydx.DeriveStarkKey(dydx.NewEcdsaPrivateKeySigner(private_key), false)
	if err != nil {
		t.Fatalf("failed to derive stark key: %v", err)
	}

	if !cmp.Equal(stark_key_python_testnet, stark_key_testnet) {
		t.Fatalf("default implementation: %s is different from golang implementation: %s",
			config.Sdump(stark_key_python_testnet), config.Sdump(stark_key_testnet))
	}

	stark_key_python_mainnet := dydx.NewStarkKey(
		"0470eb1f2bb80a4785e328efbcc41a4fad9060df9b210b848ad12fd4251a4fdc",
		"02f2fc80a218b2dee950c9ed69218e5b89bcb205c48b494c4fc6e856ecc31d7a",
		"01b480e13db79d66cc62dba6dc64536fa3242531a0c840add19cf1a04dd77866",
	)

	stark_key_mainnet, err := dydx.DeriveStarkKey(dydx.NewEcdsaPrivateKeySigner(private_key), true)
	if err != nil {
		t.Fatalf("failed to derive stark key: %v", err)
	}

	if !cmp.Equal(stark_key_python_mainnet, stark_key_mainnet) {
		t.Fatalf("default implementation: %s is different from golang implementation: %s",
			config.Sdump(stark_key_python_mainnet), config.Sdump(stark_key_mainnet))
	}
}
