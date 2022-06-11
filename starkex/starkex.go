package starkex

import (
	"fmt"
	"math/big"
)

func NewSigner(starkPrivateKey string) *Signer {
	s := new(Signer)
	s.starkPrivateKey = starkPrivateKey

	return s
}

func WithdrawSign(starkPrivateKey string, param WithdrawSignParam) (string, error) {
	return NewSigner(starkPrivateKey).SignWithdraw(param)
}

func TransferSign(starkPrivateKey string, param TransferSignParam) (string, error) {
	return NewSigner(starkPrivateKey).SignTransfer(param)
}

func OrderSign(starkPrivateKey string, param OrderSignParam) (string, error) {
	return NewSigner(starkPrivateKey).SignOrder(param)
}

func PrivateKeyToEcPointOnStarkCurv(priv_key *big.Int) (*big.Int, *big.Int, error) {
	if priv_key.Sign() < 0 || priv_key.Cmp(EC_ORDER) >= 0 {
		return nil, nil, fmt.Errorf("private key is invalid: %s", priv_key.String())
	}

	x := ecMult(priv_key, pedersenCfg.ConstantPoints[1], pedersenCfg.ALPHA, FIELD_PRIME)

	return x[0], x[1], nil
}
