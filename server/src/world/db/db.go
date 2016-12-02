// 游戏数据库
package db

//
import (
	"log"

	. "core/db/couchbase"
)

// CbEngineEx
var g_cb *CbEngineEx

// 获取g_cb
func DB() *CbEngineEx { return g_cb }

// serve ex
func Serve(param *ParamInit) error {
	g_cb = NewCbEngineEx()

	log.Println(g_cb.Name(), "serve start")
	defer log.Println(g_cb.Name(), "serve end")

	return g_cb.Serve(param)
}
