package key_utils

import (
	"strconv"
)

type KeyUtils struct {
	BaseName string
}

func (k *KeyUtils) GetUserKey(uid uint) string {
	return k.BaseName + ":user:" + strconv.Itoa(int(uid))
}

// GetTokenKey 存放 userJwt struct
func (k *KeyUtils) GetTokenKey(token string) string {
	return k.BaseName + ":token:" + token
}
