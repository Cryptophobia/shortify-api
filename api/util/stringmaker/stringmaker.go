package stringmaker

import (
	"math/rand"
	"time"
)

// Random charset that we want to use for now
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func NewStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewString(length int) string {
	return NewStringWithCharset(length, charset)
}

func GetCurrentCharset() string {
	return charset
}
