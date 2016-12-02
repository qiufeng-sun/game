// 平台驱动
package driver

import (
	"fmt"
	"strconv"
	"net/http"
)

// 平台handler接口
type IHandler interface {
	HandleAddRecharge(w http.ResponseWriter, r *http.Request)
}

// handler
var handlers = make(map[int] IHandler)

// 注册平台
func Register(platformId int, handler IHandler) {
	if nil == handler {
		panic("platform handler: Register hander is nil")
	}
	if _, dup := handlers[platformId]; dup {
		panic("platform handler: Register called twice for handler " + strconv.Itoa(platformId))
	}
	handlers[platformId] = handler
}

// 根据平台id获取handler
func GetHandler(platformId int) (IHandler, error) {
	handler, ok := handlers[platformId]
	if !ok {
		fmt.Println("cannt found platform handler. platform id:", platformId)
		return nil, fmt.Errorf("cannt found platform handler. platform id: %d", platformId)
	}
	
	return handler, nil
}