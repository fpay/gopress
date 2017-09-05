# gopress

[![Build Status](https://travis-ci.org/fpay/gopress.svg)](https://travis-ci.org/fpay/gopress)
[![Go Report Card](https://goreportcard.com/badge/github.com/fpay/gopress)](https://goreportcard.com/report/github.com/fpay/gopress)
[![codecov](https://codecov.io/gh/fpay/gopress/branch/master/graph/badge.svg)](https://codecov.io/gh/fpay/gopress)

基于 [Echo](https://echo.labstack.com/) 的应用级Golang Web框架。

## 特性

1. 基于Echo
2. View层使用 [Handlebars](http://handlebarsjs.com/) ([raymond](https://github.com/aymerick/raymond)) 做为模板引擎，支持全局 Partials
3. Controller 层抽象
4. Service 层抽象
5. Service 注册中心
6. Logging 使用[logrus](https://github.com/sirupsen/logrus)，输出格式为JSON
7. 脚手架代码生成：service, controller, middleware, model, view, main

## 安装

使用 `go get` 命令安装最新版本。下面的代码同时安装 `gopress` 和它的命令行工具。

```sh
go get -u github.com/fpay/gopress/gopress
```

然后在项目中引入 `gopress`

```go
import "github.com/fpay/gopress"
```

## 脚手架

按照上面的步骤安装了命令行工具后，就可以使用项目脚手架了。

```sh
# 生成 main.go
gopress make entry

# 生成 users 和 posts 控制器
gopress make controller users posts

# 生成 login 页面模板
gopress make view login

# 生成 cache 服务
gopress make service cache
```

详情可执行 `gopress make --help` 查看具体使用方法。

## Sample Project

[gopress-kick-starter](https://github.com/jerray/gopress-kick-starter)

## License

MIT
