package db

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// 添加充值记录参数
type ArgsAddRecharge struct {
	Platform int
	OrderId string
	CpOrderId string
	PlayerId int64
	ProductId int
	ProductNum int
	PayFee int
	PayTime string
}

// 充值rpc对象
type Recharge int

// 充值
func (r *Recharge) AddRecharge(args *ArgsAddRecharge, errNum *uint16) error {
	sql := "INSERT recharge SET platform=?,orderId=?,cpOrderId=?,player=?,product=?,num=?,payFee=?,payTime=?"
	
	stmt, err := sqlDB.Prepare(sql);
	if err != nil {
		(*errNum) = err.(*mysql.MySQLError).Number
		fmt.Println("AddRecharge err:", err)
		return nil
	}
	
	if _, err := stmt.Exec(args.Platform, args.OrderId, args.CpOrderId, args.PlayerId, 
					args.ProductId, args.ProductNum, args.PayFee, args.PayTime); err != nil {
		(*errNum) = err.(*mysql.MySQLError).Number
		fmt.Println("AddRecharge err:", err)
		return nil
	}
	
	*errNum = E_MySql_Success
	return nil
}
