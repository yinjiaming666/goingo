package confg

import (
	"app/tools/conv"
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	viper    *viper.Viper
	Path     string // 配置目录
	FileName string // 文件名
	FileType string // 文件类型
}

//var v = viper.New()

func (c *Config) Init() *Config {
	if c.Path == "" {

	}
	if c.FileType == "" {
		c.FileType = "ini"
	}
	c.viper = viper.New()
	c.viper.AddConfigPath(c.Path)     // 文件所在目录
	c.viper.SetConfigName(c.FileName) // 文件名
	c.viper.SetConfigType(c.FileType) // 文件类型
	if err := c.viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("config not found")
			fmt.Println(err.Error())
			panic(err.Error())
		}
	}
	return c
}

func Get[T conv.BuiltinT](c *Config, section, key string) T {
	key = section + "." + key
	var t T
	switch any(t).(type) {
	case string:
		return any(c.viper.GetString(key)).(T)
	case bool:
		return any(c.viper.GetBool(key)).(T)
	case int:
		return any(c.viper.GetInt(key)).(T)
	case int32:
		return any(c.viper.GetInt32(key)).(T)
	case int64:
		return any(c.viper.GetInt64(key)).(T)
	case uint:
		return any(c.viper.GetUint(key)).(T)
	case uint16:
		return any(c.viper.GetUint16(key)).(T)
	case uint32:
		return any(c.viper.GetUint32(key)).(T)
	case uint64:
		return any(c.viper.GetUint64(key)).(T)
	case float64:
		return any(c.viper.GetFloat64(key)).(T)
	}
	return any(c.viper.GetString(key)).(T)
}
