package mi

import (
	"fmt"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
	"crypto/hmac"
	"crypto/sha1"
	
	"platform/driver"
	."platform/conf"
	"conf"
	"db"
)

// 订单处理返回码
const (
	e_Invalid				= "invalid"			// dail或者db错误导致 -- 不反馈给平台
	e_Success				= "200"				// 成功
	e_Failed_CpOrderId		= "1506"			// 
	e_Failed_AppId			= "1515"			// 
	e_Failed_UId			= "1516"			// 
	e_Failed_Signature		= "1525"			// 
	e_Failed_OrderInfo		= "3515"			// 订单信息不一致，用于和CP 的订单校验
)

// 订单状态
const (
	e_Order_AddOk			= "TRADE_SUCCESS"	// 平台支付成功
)

// 订单请求中userInfo中数据记录顺序
const (
	e_CpUserInfo_PlayerId	= iota
	e_CpUserInfo_ProductId
	e_CpUserInfo_PayFee
)

// 充值信息枚举 -- 按字母升序定义
const (
	e_AddOrder_appId			= iota
	e_AddOrder_cpOrderId
	e_AddOrder_cpUserInfo
	e_AddOrder_orderId
	e_AddOrder_orderStatus
	e_AddOrder_payFee
	e_AddOrder_payTime
	e_AddOrder_productCode
	e_AddOrder_productCount
	e_AddOrder_productName
	e_AddOrder_uid
	e_AddOrder_Num
)

// 充值信息key名
var rechargeKey	[e_AddOrder_Num]string

// 充值信息
type rechargeInfo struct {
	args		[e_AddOrder_Num]string
	signature	string
}

// 平台接口对象
type myDriver struct {} 

// 充值
func (d *myDriver) HandleAddRecharge(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	var info rechargeInfo
	
	r.ParseForm()
	if "GET" == r.Method {
		// 充值参数
		for i:=0; i<e_AddOrder_Num; i++ {
			info.args[i] = r.Form.Get(rechargeKey[i])
		}
		
		// 签名
		info.signature = r.Form.Get("signature")
		
		// 充值
		retErr := recharge(&info)
		if retErr != e_Success {
			fmt.Println("recharge err:", retErr)
			
			if e_Invalid == retErr {
				return
			}
		}

		// 反馈
		fmt.Fprintf(w, "{\"errcode\":\"" + retErr + "\"}")
	}
}


// 充值
func recharge(info *rechargeInfo) string {
	// 支付失败
	if info.args[e_AddOrder_orderStatus] != e_Order_AddOk {
		fmt.Println(1)
		return e_Failed_OrderInfo
	}
	
	// appId不符
	if info.args[e_AddOrder_appId] != Cfg.AppId {
		return e_Failed_AppId
	}
	
	// rpc参数
	var err error
	var rpcArgs db.ArgsAddRecharge
	
	// cp order id
	rpcArgs.CpOrderId = info.args[e_AddOrder_cpOrderId]

	// cp user info
	cpUserInfos := strings.Split(info.args[e_AddOrder_cpUserInfo], "|")
	
	if rpcArgs.PlayerId, err = strconv.ParseInt(cpUserInfos[e_CpUserInfo_PlayerId], 10, 64); err != nil {
		fmt.Println(2)
		return e_Failed_OrderInfo
	}
	
	if rpcArgs.ProductId, err = strconv.Atoi(cpUserInfos[e_CpUserInfo_ProductId]); err != nil {
		fmt.Println(3)
		return e_Failed_OrderInfo
	}
		
	var payFee int
	if payFee, err = strconv.Atoi(info.args[e_AddOrder_payFee]); err != nil {
		fmt.Println(4)
		return e_Failed_OrderInfo
	}
	
	if rpcArgs.PayFee, err = strconv.Atoi(cpUserInfos[e_CpUserInfo_PayFee]); err != nil {
		fmt.Println(5)
		return e_Failed_OrderInfo
	}
	
	if rpcArgs.PayFee * 100 != payFee {
		fmt.Println(6)
		return e_Failed_OrderInfo
	}
	
	// 商品个数
	if rpcArgs.ProductNum, err = strconv.Atoi(info.args[e_AddOrder_productCount]); err != nil {
		fmt.Println(7)
		return e_Failed_OrderInfo
	}
	
	// 签名
	if !checkSignature(info) {
		return e_Failed_Signature
	}
	
	// 平台id
	rpcArgs.Platform = Cfg.PlatformId
	
	// 平台order id
	rpcArgs.OrderId = info.args[e_AddOrder_orderId]
	
	// 支付时间
	rpcArgs.PayTime = info.args[e_AddOrder_payTime]
	
	// 连接billing
	client, err := rpc.DialHTTP("tcp", conf.GetHttpAddr(&Cfg.DailHttp))
	if err != nil {
		fmt.Println("billing connect err:", err)
		return e_Invalid
	}
	defer client.Close()

	// Synchronous call
	var errNum uint16
	client.Call("Recharge.AddRecharge", rpcArgs, &errNum)
	
	// 不是记录已存在
	if errNum != db.E_MySql_Success && errNum != db.E_MySql_Duplicate {
		return e_Invalid
	}
	
	return e_Success
}

// 签名检查
func checkSignature(info *rechargeInfo) bool {
	// 签名字符串 -- 升序后拼接成par1=val1&par2=val2&par3=val3
	var srcStr string
	for i:=0; i<e_AddOrder_Num; i++ {
		if info.args[i] != "" {
			srcStr += rechargeKey[i] + "=" + info.args[i] + "&"
		}
	}
	
	srcLen := len(srcStr)
	if srcLen <= 0 {
		return false
	}
	
	// 
	var checkStr string
	
	// hmac
	hmacInst := hmac.New(sha1.New, []byte(Cfg.AppKey))
	hmacInst.Write([]byte(srcStr)[:srcLen - 1])
	res := hmacInst.Sum([]byte(""))
	for _,ch := range res {
		checkStr += fmt.Sprintf("%02x", ch)
	}
	
	return checkStr == info.signature
}

// 包初始化
func init() {
	driver.Register(X_PlatformId_Mi, &myDriver{})

	rechargeKey[e_AddOrder_appId]		= "appId"
	rechargeKey[e_AddOrder_cpOrderId]	= "cpOrderId"
	rechargeKey[e_AddOrder_cpUserInfo]	= "cpUserInfo"
	rechargeKey[e_AddOrder_orderId]		= "orderId"
	rechargeKey[e_AddOrder_orderStatus]	= "orderStatus"
	rechargeKey[e_AddOrder_payFee]		= "payFee"
	rechargeKey[e_AddOrder_payTime]		= "payTime"
	rechargeKey[e_AddOrder_productCode]	= "productCode"
	rechargeKey[e_AddOrder_productCount]= "productCount"
	rechargeKey[e_AddOrder_productName]	= "productName"
	rechargeKey[e_AddOrder_uid]			= "uid"	
}

