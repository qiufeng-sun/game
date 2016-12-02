// world server入口文件
package main

//
import (
	"runtime/debug"

	"util/logs"

	"core/server"
	"world/world"
)

// 程序入口
func main() {
	defer func() {
		if e := recover(); e != nil {
			logs.GetLogger().Critical("panic:%v", e)
			logs.Warn(string(debug.Stack()))
		}
		logs.Close()
	}()

	server.Run(world.New())
}
