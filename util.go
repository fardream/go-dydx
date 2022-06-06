package dydx

import (
	"encoding/json"
	"fmt"
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

func isStringValid(v string, validValues []string) (bool, int) {
	upper := strings.ToUpper(v)
	for i, t := range validValues {
		if upper == t {
			return true, i
		}
	}

	return false, 0
}

func unmarshalJsonForStringEnum[T ~string](input []byte, typename string, ot *T, validValues []string) error {
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	if ok, index := isStringValid(s, validValues); !ok {
		return fmt.Errorf("%s is not a valid %s", s, typename)
	} else {
		*ot = (T)(validValues[index])
		return nil
	}
}

func marshalJsonForStringEnum[T ~string](ot T, typename string, validValues []string) ([]byte, error) {
	if ok, index := isStringValid(string(ot), validValues); !ok {
		return nil, fmt.Errorf("%s is not a valid OrderType", ot)
	} else {
		s := string(validValues[index])
		return json.Marshal(s)
	}
}

func getProperStringEnum[T ~string](input string, validTypes []string, typename string) (T, error) {
	ok, index := isStringValid(input, validTypes)
	if !ok {
		return T(""), fmt.Errorf("invalid %s %s", typename, input)
	}
	return T(validTypes[index]), nil
}
