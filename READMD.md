# Goingo

基于 Gin + Gorm 整合的开发框架，用于快速构建 API 服务

#### 使用技术
- 路由，中间件 [Gin](https://gin-gonic.com/zh-cn/)
- model [Gorm](https://gorm.io/zh_CN/)
- 配置文件解析 [viper](https://github.com/spf13/viper/)
- [jwt](https://github.com/golang-jwt/jwt/)
- [redis](https://redis.uptrace.dev/zh/guide/)


#### 目录结构
```
├── config              // 项目配置文件
│   ├── dev.ini         // 开发环境
│   ├── prod.ini        // 测试环境
│   └── server.ini      // 服务器信息
├── deploy              // 打包上传到正式环境
│   ├── deploy.go
│   ├── deploy.sh
│   └── run.sh
├── internal            // 业务代码
│   ├── logic           // 业务逻辑
│   ├── middleware      // 中间件
│   ├── model           // 模型
│   ├── router          // 路由
│   └── server          // 接口
├── log                 // 运行日志
├── main.go             // 入口文件
└── tools               // 通用工具
```

#### 运行
```shell
go run main.go -mode=dev
# 运行参数
# -mode=dev    运行测试环境 dev.ini
# -mode=prod   运行正式环境 prod.ini
# -initDb=true 根据结构体初始化数据库
```
#### 打包上传到服务器
```shell
cd deploy && go run deploy.go
```
