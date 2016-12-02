// 资源读取
package res

// import
import (
	"util/loader"
	"world/define"
)

// 资源读取结果
var g_allOk bool

// 读取资源
func LoadAll() bool {
	g_allOk = true

	var e error

	e = loader.ParseJsonFile(File_WorldCfg, define.GetServerCfg())
	checkParseErr(e)

	e = loader.ParseXmlFile(File_ItemEntry, define.GetItemEntrys())
	checkParseErr(e)

	return g_allOk
}

// 结果检查
func checkParseErr(e error) {
	if g_allOk && e != nil {
		g_allOk = false
	}
}
