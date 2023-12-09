package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = chars[rnd.Intn(len(chars))]
	}

	return string(bytes)
}
