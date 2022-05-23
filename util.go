package dydx

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var namespace uuid.UUID

func init() {
	if err := namespace.UnmarshalText([]byte("0f9da948-a6fb-4c45-9edc-4685c3f3317d")); err != nil {
		panic(err)
	}
}

func GetUserId(ethAddress string) string {
	return uuid.NewSHA1(namespace, []byte(ethAddress)).String()
}

func GetAccountIdFromEth(ethAddress string) string {
	return uuid.NewSHA1(namespace, []byte(GetUserId(strings.ToLower(ethAddress))+strconv.Itoa(0))).String()
}

func GetIsoDateStr(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
