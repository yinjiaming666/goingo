package random

import (
	"math/rand"
	"time"
)

const defaultLens = 30

var runes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func Str(lens int) string {
	if lens <= 0 {
		lens = defaultLens
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, lens)
	for i := range b {
		b[i] = runes[r.Intn(len(runes))]
	}
	return string(b)
}
