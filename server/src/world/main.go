// world server入口文件
package main

//
import (
	"core"
	"core/server"
	"world/world"
)

// 程序入口
func main() {
	defer core.PrintPanic()

	server.Run(world.New())
}
