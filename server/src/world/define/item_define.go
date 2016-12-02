// 物品定义
package define

//
import (
//	"encoding/xml"
)

// 物品静态属性
type ItemEntry struct {
//	XMLName     xml.Name 		`xml:"item"`
	Id			int				`xml:"id,attr"`
	Type		int				`xml:"type,attr"`
	Qlty		int				`xml:"qlty,attr"`
}

// 
type ItemEntryMap map[int]*ItemEntry

// 物品资源管理对象定义
type ItemEntrys struct {
//	XMLName     xml.Name 		`xml:"root"`
	Entrys		[]ItemEntry		`xml:"item"`
	entrys		ItemEntryMap
}

// 资源接口实现
// IEntrys
func (entrys *ItemEntrys) InitMap() {
	for _, entry := range entrys.Entrys {
		entrys.entrys[entry.Id] = &entry
	}
}

// 资源管理对象
var g_itemEntrys = ItemEntrys{ entrys: make(ItemEntryMap) }

// 获取资源管理对象
func GetItemEntrys() *ItemEntrys {
	return &g_itemEntrys
}

// 获取指定
func GetItemEntry(id int) *ItemEntry {
	return g_itemEntrys.entrys[id]
}

