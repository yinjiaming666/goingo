package router

import (
	"app/internal/global"
	middleware2 "app/internal/middleware"
	"app/tools/logger"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type RouteGateway interface {
	initRoute()
}

func InitRouter(port string) {
	var err error

	if global.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware2.CORSMiddleware()) // 解决跨域

	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	f, _ := os.OpenFile(logger.AccessLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	c := gin.LoggerConfig{
		Output:    f,
		SkipPaths: []string{"/favicon.ico"},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\" POSTFORM - [%s] \n",
				params.ClientIP,
				params.TimeStamp.Format(time.DateTime),
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

	r.Use(middleware2.RespMiddleware()) // 响应中间件
	rr := r.Group(global.Version)
	AdminRoute := AdminRoute{group: rr.Group("api/admin")}
	IndexRoute := IndexRoute{group: rr.Group("api/index")}
	AdminRoute.initRoute()
	IndexRoute.initRoute()

	err = r.Run(":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
}
