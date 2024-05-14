package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

var config = viper.New()
var RootPath = GetRootPath()

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
	config.AddConfigPath(RootPath + "/config/") // 文件所在目录
	config.SetConfigName(file)                  // 文件名
	config.SetConfigType("ini")                 // 文件类型
	if err := config.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("config not found")
			fmt.Println(err.Error())
			panic(err.Error())
		}
	}
	key = section + "." + key
	return config.GetString(key)
}

// GetRootPath 获取项目根目录
func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../")
}
