package main

import (
	"fmt"
	"os"

	"github.com/fardream/go-dydx"
)

type (
	starkKey dydx.StarkKey
	apiKey   dydx.ApiKey
)

func (c *starkKey) String() string {
	return "empty"
}

func (c *starkKey) Set(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m, err := dydx.ParseStarkKeyMap(data)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("only one keys is allowed: %s", data)
	}
	for _, v := range m {
		*c = (starkKey)(*v)
	}
	return nil
}

func (c *starkKey) Type() string {
	return "stark-key-map-file"
}

func (c *apiKey) String() string {
	return "empty"
}

func (c *apiKey) Set(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m, err := dydx.ParseApiKeyMap(data)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("only one keys is allowed: %s", data)
	}
	for _, v := range m {
		*c = (apiKey)(*v)
	}
	return nil
}

func (c *apiKey) Type() string {
	return "api-key-map-file"
}
