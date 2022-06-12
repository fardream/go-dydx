package dydx_test

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fardream/go-dydx"
)

var _ = "keep"

func ExampleDeriveStarkKey() {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate an ethereum key: %#v", err)
	}

	stark_key, err := dydx.DeriveStarkKey(dydx.NewEcdsaPrivateKeySigner(key), false)
	if err != nil {
		log.Fatalf("failed to derivate stark key: %#v", err)
	}

	fmt.Println(stark_key)
}
