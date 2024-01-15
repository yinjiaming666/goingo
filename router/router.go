package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goingo/logger"
	"goingo/middleware"
	"io"
	"os"
	"time"
)

type RouteGateway interface {
	initRoute()
}

func InitRouter(port string) {
	r := gin.New()
	r.Use(middleware.CORSMiddleware()) // 解决跨域

	f, _ := os.OpenFile(logger.AccessLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	c := gin.LoggerConfig{
		Output:    f,
		SkipPaths: []string{"/test"},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\" POSTFORM - [%s] \n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
				params.Request.PostForm,
			)
		},
	}
	r.Use(gin.LoggerWithConfig(c))

	r.Use(middleware.RespMiddleware()) // 响应中间件

	AdminRoute := AdminRoute{group: r.Group("api/admin")}
	IndexRoute := IndexRoute{group: r.Group("api/index")}
	AdminRoute.initRoute()
	IndexRoute.initRoute()

	r.Run(":" + port)
}
