// protobuf消息 -- 格式: 消息大小 + 消息Id + 消息内容(protobuf)
package protobuf

// import
import (
	"github.com/golang/protobuf/proto"

	"core/net/msg"
)

//// 消息
//type PbMsg struct {
//	id					int32				// 消息id
//	data				proto.Message		// protobuf
//}

//// 创建
//func NewPbMsg(msgId int32, pb proto.Message) *PbMsg {
//	return &PbMsg{msgId, pb}
//}

//// 序列化为[]byte切片
//func (pb *PbMsg) Marshal() ([]byte, []byte, error) {
//	return Marshal(pb.id, pb.data)
//}

// 序列化消息
func marshal(msgData interface{}) ([]byte, error) {
	return proto.Marshal(msgData.(proto.Message))
}

// 解析消息
func unmarshal(buf []byte, out interface{}) error {
	return proto.Unmarshal(msg.GetMsgData(buf), out.(proto.Message))
}

//
type PbParser struct{}

//
func (_ PbParser) Marshal(msgId uint32, msgData interface{}) ([]byte, []byte, error) {
	return msg.Marshal(msgId, msgData, marshal)
}

//
func (_ PbParser) Unmarshal(buf []byte, out interface{}) error {
	return unmarshal(buf, out)
}
