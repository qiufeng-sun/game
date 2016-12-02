// 消息处理函数注册
package hfunc

//
import (
	. "msg/pbmsg"
)

// param: msgId, msgHandler, gmLv, inWorld, serial
func Register() {
	// 进入world前的消息注册(param: msgId, msgHandler, gmLv)
	RegOutWorldHandler(EMsgId_ID_NcCreatePlayer, HandleCreatePlayer, 0)

	// 进入world后的串行消息注册(param: msgId, msgHandler, gmLv)
	//RegInWorldSerialHandler()

	// 进入world后的并行消息注册(param: msgId, msgHandler, gmLv)
	//RegInWorldParallelHandler
}
