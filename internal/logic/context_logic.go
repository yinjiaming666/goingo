package logic

import (
	"app/internal/model"
	"app/tools/jwt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ContextLogic struct {
	Context *gin.Context
}

var ContextLogicInstance *ContextLogic

func init() {
	ContextLogicInstance = &ContextLogic{}
}

func (c *ContextLogic) SetContext(gc *gin.Context) {
	c.Context = gc
}

// GetIndexUserInfo 获取保存在上下文中的用户信息
func (c *ContextLogic) GetIndexUserInfo() (*model.User, error) {
	user, ok := c.Context.Get(string(jwt.IndexJwtType))
	if !ok {
		return nil, errors.New("not found Context user")
	}
	m := user.(*model.User)
	return m, nil
}
