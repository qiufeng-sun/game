// 平台对接http server
package main

import (
	"fmt"
	"net/http"
	
	"server"
	"platform/driver"
	"conf"
	."platform/conf"
	
	_ "platform/mi"
)

// 平台对接http server
type Platform struct {
	server.ServerBase
}

// 配置
var handler driver.IHandler

// 
func (s *Platform) String() string {
	return "Platform http server"
}

// 初始化
func (s *Platform) Init() bool {
	var err error
	
	// 读取配置
	if err = conf.LoadJsonConf(File_PlatformConf, &Cfg); err != nil {
		return false
	}
	
	// handler
	if handler, err = driver.GetHandler(Cfg.PlatformId); err != nil {
		return false
	}

	// http server -- 监听平台连接
	go func() {
		if err = http.ListenAndServe(conf.GetHttpAddr(&Cfg.LsnHttp), nil); err != nil {
			fmt.Println("platform http err:", err)
		}
	}()

	return true
}
