package context

import (
	"app/internal/model"
	"app/tools/jwt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// GetIndexUserInfo 获取保存在上下文中的用户信息
func GetIndexUserInfo(c *gin.Context) (*model.User, error) {
	user, ok := c.Get(string(jwt.IndexJwtType))
	if !ok {
		return nil, errors.New("not found Context user")
	}
	m := user.(*model.User)
	return m, nil
}
