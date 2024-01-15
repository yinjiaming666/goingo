package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/spf13/viper"
	"goingo/logger"
)

var config = viper.New()

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

func GetConfig(file string, section string, key string) string {
	config.AddConfigPath("./config/") // 文件所在目录
	config.SetConfigName(file)        // 文件名
	config.SetConfigType("ini")       // 文件类型
	if err := config.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logger.Error("config not found", err.Error())
		} else {
			logger.Error("config read error", err.Error())
		}
		panic(err.Error())
	}
	key = section + "." + key
	return config.GetString(key)
}
