// json消息 -- 格式: 消息大小 + 消息Id + 消息内容(json)
package json

// import
import (
	"encoding/json"

	"core/net/msg"
)

// 序列化消息
func marshal(msgData interface{}) ([]byte, error) {
	return json.Marshal(msgData)
}

// 解析消息
func unmarshal(buf []byte, out interface{}) error {
	return json.Unmarshal(msg.GetMsgData(buf), out)
}

//
type JsonParser struct{}

//
func (_ JsonParser) Marshal(msgId uint32, msgData interface{}) ([]byte, []byte, error) {
	return msg.Marshal(msgId, msgData, marshal)
}

//
func (_ JsonParser) Unmarshal(buf []byte, out interface{}) error {
	return unmarshal(buf, out)
}
