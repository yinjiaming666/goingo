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

func Number(start, end int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if end < start {
		t := end
		end = start
		start = t
	}
	end++
	return r.Intn(end-start) + start // (end-start)+start
}

// RandSlice 随机弹出切片中的元素
func RandSlice[T any](slice []T) (T, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(slice) == 0 {
		panic("empty slice")
	}
	randIndex := r.Intn(len(slice))
	return slice[randIndex], randIndex
}
