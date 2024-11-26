package utils

import (
	"math/rand"
	"strings"
)

// di go 1.20 udah ga perlu pake ini kalo buat random doang
// func init() {
// 	rand.Seed( )
// }

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1) // rand.Intn(10) [0, 9] jadi misal, min 5 max 10 maka 5, 6, 7, 8, 9 karena rand.Intn(5+1) [0, 5] + 5
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
