// world server配置结构定义
package define

//
//import (
//	cb "core/db/couchbase"
//)

// world配置
type WorldCfg struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// 连接配置 -- log db
type LogDBCfg struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Pwd    string `json:"pwd"`
	DBName string `json:"db_name"`
}

// 连接配置 -- logon
type LogonConnCfg struct {
	Addr string `json:"addr"`
}

// 连接配置 -- 玩家
type PlayerConnCfg struct {
	LstAddr string `json:"lst_addr"`
	MaxNum  int    `json:"max_num"`
}

// 服务器配置
type ServerCfg struct {
	World WorldCfg `json:"world"`
	//	WorldDB    cb.ParamInit  `json:"world_db"`// to do
	LogDB      LogDBCfg      `json:"log_db"`
	LogonConn  LogonConnCfg  `json:"logon_conn"`
	PlayerConn PlayerConnCfg `json:"player_conn"`
}

// 配置对象
var g_srvCfg *ServerCfg = new(ServerCfg)

// 获取配置对象
func GetServerCfg() *ServerCfg {
	return g_srvCfg
}

//
func GetPlayerConn() *PlayerConnCfg {
	return &g_srvCfg.PlayerConn
}
