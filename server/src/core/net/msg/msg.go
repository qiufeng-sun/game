// 消息 -- 格式: 消息头(消息大小＋消息id), 消息体
package msg

// import
import (
	"encoding/binary"

	"core/buff"

	. "msg/pbmsg"
)

// 常量
const (
	x_MsgSzBytes  = 4                        // 消息大小所占空间(header1)
	x_MsgIdSz     = 4                        // id大小(header2)
	x_MsgHeaderSz = x_MsgSzBytes + x_MsgIdSz // 消息头(消息大小＋消息id)大小
)

// 字节序
var g_byteOrder binary.ByteOrder = &binary.BigEndian

// 设置字节序
func SetByteOrder(order binary.ByteOrder) {
	g_byteOrder = order
}

// 错误码转换
func Err2Int32(err EMsgErr) *int32 {
	n := int32(err)
	return &n
}

// 从iov中读取uint32
func Uint32(iov *buff.IoVector) (uint32, bool) {
	var data = iov.Bytes(4, 0)
	if nil == data {
		return 0, false
	}

	return g_byteOrder.Uint32(data), true
}

func Uint32ByBytes(data []byte) (uint32, bool) {
	if len(data) < 4 {
		return 0, false
	}

	return g_byteOrder.Uint32(data), true
}

// 将uint32写入[]byte中, 并返回
func PutUint32(data []byte, val uint32) []byte {
	g_byteOrder.PutUint32(data, val)
	return data
}

// 检查消息是否完整 -- 返回值: (消息大小, 是否完成)
func Check(iov *buff.IoVector) (int, bool) {
	var sz, ok = Uint32(iov)
	if !ok || int(sz) <= 0 {
		return 0, false
	}

	return int(sz), iov.Size() >= int(sz)+x_MsgSzBytes
}

// 获取一条完整消息的切片 -- 返回值: (消息id + 消息数据, 是否有消息)
func Get(iov *buff.IoVector) ([]byte, bool) {
	// 消息完整性检查
	sz, ok := Check(iov)
	if !ok {
		return nil, false
	}

	return iov.Bytes(sz, x_MsgSzBytes), true
}

// 计算消息在缓存中占用的空间
func CalBuffSpace(bufSz int) int { return bufSz + x_MsgSzBytes }

// 获取消息id
func ParseMsgId(buf []byte) (int32, bool) {
	msgId, ok := Uint32ByBytes(buf)

	return int32(msgId), ok
}

// 获取消息bytes
func GetMsgData(buf []byte) []byte {
	return buf[x_MsgIdSz:]
}

// @return 消息头，消息体，error
type Marshaler interface {
	Marshal(msgId uint32, msgData interface{}) ([]byte, []byte, error)
}

// parser
type Unmarshaler interface {
	Unmarshal(data []byte, out interface{}) error
}

// msg
type Parser interface {
	Marshaler
	Unmarshaler
}

// marshal
// @return: 消息头(消息大小＋消息id), 消息体, error
func Marshal(msgId uint32, msgData interface{},
	marshal func(data interface{}) ([]byte, error)) ([]byte, []byte, error) {
	b2, e := marshal(msgData)
	if e != nil {
		return nil, nil, e
	}

	b1 := make([]byte, x_MsgHeaderSz)
	PutUint32(b1[:x_MsgSzBytes], uint32(len(b2)+x_MsgIdSz))
	PutUint32(b1[x_MsgSzBytes:], msgId)

	return b1, b2, nil
}
