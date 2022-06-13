package dydx_test

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fardream/go-dydx"
)

const createUserIsMainnet = false

// An example that generate a new private key for ethereum network,
// use the default stark key and api credentials, then create the user on the
// testnet, and request token airdrops.
func ExampleClient_CreateUser() {
	// create private key on the ethereum network
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("testing private key: %s\n", hexutil.Encode(crypto.FromECDSA(privateKey)))

	// public key address
	ethAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	fmt.Printf("testing public key: %s\n", ethAddress.String())

	// signer
	signer := dydx.NewEcdsaPrivateKeySigner(privateKey)

	// stark key
	starkKey, err := dydx.DeriveStarkKey(signer, createUserIsMainnet)
	if err != nil {
		panic(err)
	}

	// api key
	apiKey, err := dydx.RecoverDefaultApiKeyCredentials(signer, false)
	if err != nil {
		panic(err)
	}

	// create a new client
	client, err := dydx.NewClient(starkKey, apiKey, ethAddress.String(), createUserIsMainnet)
	if err != nil {
		panic(err)
	}

	// create user on the dydx
	resp, err := client.CreateUser(context.Background(), signer, nil)
	if err != nil {
		panic(err)
	} else {
		spew.Dump(resp)
	}

	// request airdrop
	_, err = client.RequestTestnetTokens(context.Background())
	if err != nil {
		panic(err)
	}
}
