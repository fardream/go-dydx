package dydx_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"

	"github.com/fardream/go-dydx"
)

var spewConfig *spew.ConfigState = spew.NewDefaultConfig()

const privateKeyHex = "2b61fc77bbda78d7075d91ae6a58423df73f71ebe4aa36e3b587dc8e5fe1396f"

var privateKey *ecdsa.PrivateKey

func init() {
	spewConfig.DisableMethods = true
	var err error
	privateKey, err = crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}
}

func TestDeriveStarkKey(t *testing.T) {
	stark_key_python_testnet := dydx.NewStarkKey(
		"5ca235e41160c7c2a4695a98eecf7cb0ba6f7d26c4f899c192ca91912995115",
		"e07e9de2353c45e92b98e6c566fce8686f7fb620031aa7a79a5983daf22483",
		"3c6ce687a484ac1c50a48498092832957a3154c7c13237bc10df6965472e009",
	)

	stark_key_testnet, err := dydx.DeriveStarkKey(dydx.NewEcdsaPrivateKeySigner(privateKey), false)
	if err != nil {
		t.Fatalf("failed to derive stark key: %v", err)
	}

	if !cmp.Equal(stark_key_python_testnet, stark_key_testnet) {
		t.Fatalf("default implementation: %s is different from golang implementation: %s",
			spewConfig.Sdump(stark_key_python_testnet), spewConfig.Sdump(stark_key_testnet))
	}

	stark_key_python_mainnet := dydx.NewStarkKey(
		"470eb1f2bb80a4785e328efbcc41a4fad9060df9b210b848ad12fd4251a4fdc",
		"2f2fc80a218b2dee950c9ed69218e5b89bcb205c48b494c4fc6e856ecc31d7a",
		"1b480e13db79d66cc62dba6dc64536fa3242531a0c840add19cf1a04dd77866",
	)

	stark_key_mainnet, err := dydx.DeriveStarkKey(dydx.NewEcdsaPrivateKeySigner(privateKey), true)
	if err != nil {
		t.Fatalf("failed to derive stark key: %v", err)
	}

	if !cmp.Equal(stark_key_python_mainnet, stark_key_mainnet) {
		t.Fatalf("default implementation: %s is different from golang implementation: %s",
			spewConfig.Sdump(stark_key_python_mainnet), spewConfig.Sdump(stark_key_mainnet))
	}
}

func TestRecoverDefaultApiCredentials(t *testing.T) {
	api_key_mainnet, err := dydx.RecoverDefaultApiKeyCredentials(dydx.NewEcdsaPrivateKeySigner(privateKey), true)
	if err != nil {
		t.Fatalf("failed to recover api keys: %#v", err)
	}
	api_key_mainnet_python := dydx.NewApiKey("9fee0aab-870e-fea8-d724-29d1016e951e", "vydswDsZ8GjD9ARqfz02", "8k1btcHszt_6ShxTzSt1FRq5NwxaiIxhi9TTQclx")

	if !cmp.Equal(api_key_mainnet, api_key_mainnet_python) {
		t.Fatalf("default implementation: %s is different from golang implementation: %s", spewConfig.Sdump(api_key_mainnet_python), spewConfig.Sdump(api_key_mainnet))
	}

	api_key_testnet, err := dydx.RecoverDefaultApiKeyCredentials(dydx.NewEcdsaPrivateKeySigner(privateKey), false)
	if err != nil {
		t.Fatalf("failed to recover api keys: %#v", err)
	}
	api_key_testnet_python := dydx.NewApiKey("cdb76448-91af-1546-c841-2f7370164155", "QyeP9h46NTRBFQUX-eqN", "bP5omVCq4dVI1-t1vCr4t4fwC8NZ-_jJcns8espY")

	if !cmp.Equal(api_key_testnet, api_key_testnet_python) {
		t.Fatalf("default implementation: %s is different from golang implementation: %s", spewConfig.Sdump(api_key_testnet_python), spewConfig.Sdump(api_key_testnet))
	}
}
