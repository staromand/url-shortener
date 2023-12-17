package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int8) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, size)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = chars[rnd.Intn(len(chars)-1)]
	}

	return string(result)
}
