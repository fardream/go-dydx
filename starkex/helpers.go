package starkex

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/fardream/decimal"
	"golang.org/x/crypto/sha3"
)

func ToJsonString(input interface{}) string {
	js, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Println("ToJsonString error:", err.Error())
	}
	return string(js)
}

func getHash(str1, str2 string) string {
	return PedersenHash(str1, str2)
}

// NonceByClientId generate nonce by clientId
func NonceByClientId(clientId string) *big.Int {
	h := sha256.New()
	h.Write([]byte(clientId))

	a := new(big.Int)
	a.SetBytes(h.Sum(nil))
	res := a.Mod(a, big.NewInt(NONCE_UPPER_BOUND_EXCLUSIVE))
	return res
}

// SerializeSignature Convert a Sign from an r, s pair to a 32-byte hex string.
func SerializeSignature(r, s *big.Int) string {
	return IntToHex32(r) + IntToHex32(s)
}

// IntToHex32 Normalize to a 32-byte hex string without 0x prefix.
func IntToHex32(x *big.Int) string {
	str := x.Text(16)
	if len(str) < 64 {
		return strings.Repeat("0", 64-len(str)) + str
	} else {
		return str
	}
}

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// FactToCondition Generate the condition, signed as part of a conditional transfer.
func FactToCondition(factRegistryAddress string, fact string) *big.Int {
	data := strings.TrimPrefix(factRegistryAddress, "0x") + fact
	hexBytes, _ := hex.DecodeString(data)
	// int(Web3.keccak(data).hex(), 16) & BIT_MASK_250
	d := sha3.NewLegacyKeccak256()
	d.Write(hexBytes)
	var hash Hash
	r, ok := d.(io.Reader)
	if ok {
		r.Read(hash[:])
	}
	fst := hash.Big()
	fst.And(fst, BIT_MASK_250)
	return fst
}

// GetTransferErc20Fact get erc20 fact
// tokenDecimals is COLLATERAL_TOKEN_DECIMALS
func GetTransferErc20Fact(recipient string, tokenDecimals int, humanAmount, tokenAddress, salt string) (string, error) {
	// token_amount = float(human_amount) * (10 ** token_decimals)
	amount, err := decimal.NewFromString(humanAmount)
	if err != nil {
		return "", err
	}
	saltInt, ok := big.NewInt(0).SetString(salt, 0) // with prefix: 0x
	if !ok {
		return "", fmt.Errorf("invalid salt: %v,can not parse to big.Int", salt)
	}
	tokenAmountStr := amount.Mul(decimal.New(10, int32(tokenDecimals-1))).String()
	tokenAmount, ok := big.NewInt(0).SetString(tokenAmountStr, 10)
	if !ok {
		return "", fmt.Errorf("cannot get token amount: %s", tokenAmountStr)
	}

	recp_addr, err := hex.DecodeString(strings.TrimPrefix(recipient, "0x"))
	if err != nil {
		return "", err
	}
	toekn_addr, err := hex.DecodeString(strings.TrimPrefix(tokenAddress, "0x"))
	if err != nil {
		return "", err
	}

	var b []byte
	b = append(b, recp_addr...)
	b = append(b, common.LeftPadBytes(math.U256Bytes(tokenAmount), 32)...)
	b = append(b, toekn_addr...)
	b = append(b, common.LeftPadBytes(math.U256Bytes(saltInt), 32)...)
	d := sha3.NewLegacyKeccak256()

	d.Write(b)

	return hex.EncodeToString(d.Sum(nil)), nil
}

func GenerateKRfc6979(msgHash, priKey *big.Int, seed int) *big.Int {
	msgHash = big.NewInt(0).Set(msgHash) // copy
	bitMod := msgHash.BitLen() % 8
	if bitMod <= 4 && bitMod >= 1 && msgHash.BitLen() > 248 {
		msgHash.Mul(msgHash, big.NewInt(16))
	}
	var extra []byte
	if seed > 0 {
		buf := new(bytes.Buffer)
		var data interface{}
		if seed < 256 {
			data = uint8(seed)
		} else if seed < 65536 {
			data = uint16(seed)
		} else if seed < 4294967296 {
			data = uint32(seed)
		} else {
			data = uint64(seed)
		}
		_ = binary.Write(buf, binary.BigEndian, data)
		extra = buf.Bytes()
	}
	return generateSecret(EC_ORDER, priKey, sha256.New, msgHash.Bytes(), extra)
}
