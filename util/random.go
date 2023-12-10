package util

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// generates random integer [min, max]
func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// generates random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generates random Owner
func RandomOwner() string {
	return RandomString(6)
}

// generates random Balance
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

// generates random Currency
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)

	return currencies[r.Intn(n)]
}
