package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RandomInt64(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	lengthAplhabet := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rnd.Intn(lengthAplhabet)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomBalance() int64 {
	return RandomInt64(1, 10_000)
}

func RandomEntriesAmount() int64 {
	return RandomInt64(-1_000_000, 1_000_000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "IDR"}
	n := len(currencies)

	return currencies[rnd.Intn(n)]
}
