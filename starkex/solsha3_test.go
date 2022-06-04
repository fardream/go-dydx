package starkex

import "testing"

func TestSolsha3(t *testing.T) {
	hashnew, err := GetTransferErc20Fact("0x5dDDF0B14cc2132BDD2a0c4C5265b5Dc925465e0", 5, "30.2", "0xBE9a129909EbCb954bC065536D2bfAfBd170d27A", "3")
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}
	hashold := "f183a49d4f0cc97c8a4443f2603f0067b15a45a8a28259848fb06e03b3eeb329"
	if hashnew != hashold {
		t.Fatalf("new hash: %s and old hash: %s dont match", hashnew, hashold)
	}
}
