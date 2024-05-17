package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string, salt string) string {
	h := md5.New()
	b := []byte(str)
	h.Write(b)

	if salt != "" {
		s := []byte(salt)
		h.Write(s)
	}

	return hex.EncodeToString(h.Sum(nil))
}
