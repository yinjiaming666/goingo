package model

import (
	"app/tools"
	"app/tools/logger"
	"app/tools/resp"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type User struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint    `redis:"Id" json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Password        string  `redis:"Password" json:"password" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Nickname        string  `redis:"Nickname" json:"nickname" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	HeadImg         string  `redis:"HeadImg" json:"head_img" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Status          uint8   `redis:"Status" json:"status" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	CreateTime      int64   `redis:"CreateTime" json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;"` // 自动写入时间戳
	CreateTimeStr   string  `redis:"-" json:"create_time_str" gorm:"-:all"`                                              // -:all 无读写迁移权限，该字段不在数据库中
	Age             uint8   `redis:"Age" json:"age" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	Sex             uint8   `redis:"Sex" json:"sex" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	Money           float64 `redis:"Money" json:"money" gorm:"type:DECIMAL(10,2) UNSIGNED NOT NULL;default:0.00"`
}

func (user *User) GetUserInfo() {
	Db().Where(user).First(user)
}

func (user *User) DoRegister() uint {
	user.Password = tools.Md5(user.Password, UserPwdSalt)
	user.CreateTime = int64(uint(time.Now().Unix()))
	user.SaveCache()
	tx := Db().Create(user)
	if tx.Error != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "错误", Body: map[string]any{
			"error": tx.Error.Error(),
		}}).Response()
	}
	return user.Id
}

func (user *User) SetUser() {
	if user.Password != "" {
		user.Password = tools.Md5(user.Password, UserPwdSalt)
	}
	user.SaveCache()
	go func() {
		Db().Select("nickname", "password", "age", "sex", "token", "status").Model(user).Updates(user)
	}()
}

func (user *User) SaveCache() {
	key := KeyUtils.GetUserKey(user.Id)
	ctx := context.Background()
	if _, err := RedisClient.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key, "Id", user.Id)
		rdb.HSet(ctx, key, "Nickname", user.Nickname)
		rdb.HSet(ctx, key, "HeadImg", user.HeadImg)
		rdb.HSet(ctx, key, "Status", user.Status)
		rdb.HSet(ctx, key, "CreateTime", user.CreateTime)
		rdb.HSet(ctx, key, "Age", user.Age)
		rdb.HSet(ctx, key, "Sex", user.Sex)
		//rdb.HSet(ctx, key, "Money", user.Money)
		//rdb.HSet(ctx, key, "Password", user.Password)
		return nil
	}); err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
}

func (user *User) GetCache() {
	key := KeyUtils.GetUserKey(user.Id)
	if err := RedisClient.HGetAll(context.Background(), key).Scan(user); err != nil {
		if errors.Is(err, redis.Nil) {
			user.Id = 0
			return
		}
		panic(err)
	}
}

func (user *User) CacheXiaFen(field string, num float64) float64 {
	var script = `
local key = KEYS[1]
local field = ARGV[1]
local num = ARGV[2]
local a = redis.call('HGET', key, field) or 0
a = tonumber(a)
num = tonumber(num)
if a >= num then
	local res = redis.call('HINCRBYFLOAT', key, field, num)
	return res
else
	return -1
end
`
	keys := []string{KeyUtils.GetUserKey(user.Id)}
	res, err := RedisClient.Eval(context.Background(), script, keys, field, num).Float64()
	if err != nil {
		logger.Error("xia fen fail", "err", err)
		return -1
	}
	return res
}

func (user *User) CacheShangFen(field string, num float64) float64 {
	k := KeyUtils.GetUserKey(user.Id)
	r, err := RedisClient.HIncrByFloat(context.Background(), k, field, num).Result()
	if err != nil {
		return -1
	}
	return r
}

func (user *User) CacheResetFen(field string, num float64) {
	k := KeyUtils.GetUserKey(user.Id)
	_, err := RedisClient.HMSet(context.Background(), k, field, num).Result()
	if err != nil {
		logger.Error("CacheSetFen Fail", "err", err.Error())
	}
}
