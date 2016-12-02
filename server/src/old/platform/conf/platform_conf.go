package pfconf

import (
	"conf"
)

// 配置文件
const (
	File_PlatformConf			= "conf/platform_conf.json" 
)

// 平台id定义
const (
	X_PlatformId_Mi				= 1
)

// 配置
type PlatformConf struct {
	LsnHttp		conf.HttpConf	`json:"lsnHttp"`
	DailHttp	conf.HttpConf	`json:"dailHttp"`
	PlatformId	int				`json:"platformId"`
	AppId		string			`json:"appId"`
	AppKey		string			`json:"appKey"`
}

// 配置对象
var Cfg PlatformConf
