package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	characters := []rune("QWERTYUIOPASDFGHJKLZXCVBNM" +
		"qwertyuiopasdfghjklzxcvbnm" +
		"1234567890")

	buf := make([]rune, length)
	for i := range buf {
		buf[i] = characters[rnd.Intn(len(characters))]
	}

	return string(buf)
}
