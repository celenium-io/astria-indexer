package astria

import (
	"github.com/pactus-project/pactus/util/bech32m"
)

const (
	prefix = "astria"
)

func IsAddress(s string) bool {
	if len(s) != 45 {
		return false
	}
	p, _, err := bech32m.Decode(s)
	if err != nil {
		return false
	}
	return p == prefix
}

func EncodeAddress(b []byte) (string, error) {
	return bech32m.EncodeFromBase256(prefix, b)
}

func DecodeAddress(s string) ([]byte, error) {
	_, b, err := bech32m.DecodeToBase256(s)
	return b, err
}
