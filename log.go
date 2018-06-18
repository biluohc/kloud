package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/cihub/seelog"
	"github.com/labstack/echo"
)

var LOGXMLDEVEL string = "log.devel.xml"
var LOGXMLRELEASE string = "log.dist.xml"
var LOGXMLDEBUG string = "log.debug.xml"

// 日后要自己管理http的log, Gin的自定义中间件,不用标准的logger
func init() {
	var logConfig = LOGXMLDEVEL
	var envConfig = os.Getenv("GO_LOG")

	if envConfig == "1" {
		logConfig = LOGXMLDEBUG
	} else if envConfig == "2" {
		logConfig = LOGXMLRELEASE
	}

	logger, err := log.LoggerFromConfigAsFile(logConfig)

	if err != nil {
		fmt.Println(err)
	}

	log.ReplaceLogger(logger)
}

func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		err = next(c)
		if err != nil {
			c.Error(err)
		}

		stop := time.Now()

		if res.Status < 400 {
			log.Infof("[%s %s] %3d %s %s #%s &%s", stop.Sub(start).String(), c.RealIP(), res.Status, req.Method, req.RequestURI, req.UserAgent(), req.Referer())
		} else if res.Status < 500 {
			log.Warnf("[%s %s] %3d %s %s #%s &%s", stop.Sub(start).String(), c.RealIP(), res.Status, req.Method, req.RequestURI, req.UserAgent(), req.Referer())
		} else {
			log.Errorf("[%s %s] %3d %s %s #%s &%s", stop.Sub(start).String(), c.RealIP(), res.Status, req.Method, req.RequestURI, req.UserAgent(), req.Referer())
		}
		// 不能返回err，否则echo将重复调用 错误处理函数，然后重复堆积响应 --> 非法JSON
		return nil
	}
}
