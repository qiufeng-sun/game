// 平台http连接处理
package main

import (
	"fmt"
	"net/http"
)

// 测试
func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("platform http test")
	
	fmt.Fprintf(w, "Hello platform http!")
}

// 关服
func stopResponse(w http.ResponseWriter, r *http.Request) {
	StopServer()
}

// 添加充值记录
func rechargeResponse(w http.ResponseWriter, r *http.Request) {
	handler.HandleAddRecharge(w, r)
}

// 包初始化
func init() {
	http.HandleFunc("/test", testResponse)
	http.HandleFunc("/stop", stopResponse)
	http.HandleFunc("/recharge", rechargeResponse)
}
