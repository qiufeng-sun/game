// couchbase封装
// 同步操作, 直接调用响应接口XXXSync()
// 异步操作, 调用XXXAsyn(), 返回结果channel, 操作完成后底层会通知channel. channel由上层自己维护
// 操作接口在op.go中实现
package couchbase

// 
import (
	"log"
	
	gcb "github.com/couchbaselabs/go-couchbase"
)

// db操作
type IOp interface {
	Exec(engine *CbEngine)
}

// db初始化参数
type ParamInit struct {
	IP				string			`json:"ip"`
	Bucket			string			`json:"bucket"`
	Pwd				string			`json:"pwd"`
	MaxWait			int				`json:"max_wait"`		// 最大待执行操作个数
}

// 
type CbEngine struct {
	bucket			*gcb.Bucket		// bucket
	chOp			chan IOp		// 异步操作队列
	chStop			chan bool		// 结束通知(非阻塞)
	running			bool			// 是否正在执行
}

// 创建对象
func NewCbEngine() *CbEngine {
	return &CbEngine{}
}

// 获取名字
func (engine *CbEngine) Name() string {
	return "CbEngine"
}

//// run
//func (engine *CbEngine) Run(param *ParamInit) (err error) {
//	// 获得bucket
//	engine.bucket, err	= gcb.GetBucket(url(param), "default", param.Bucket)
//	if err != nil { return }
	
//	// 初始化其他属性
//	engine.chOp			= make(chan IOp, param.MaxWait)
//	engine.chStop		= make(chan bool, 1)
//	engine.running		= true
	
//	// 处理操作
//	for op := range engine.chOp {
//		op.Exec(engine)
//	}
	
//	// 通知操作完成
//	engine.chStop <- true
	
//	return nil
//}

func (engine *CbEngine) Serve(param *ParamInit) (err error) {
	// debug log
	//log.Println(url(param))
	
	// 获得bucket
	engine.bucket, err	= gcb.GetBucket(url(param), "default", param.Bucket)
	if err != nil { return }
	
	// 初始化其他属性
	engine.chOp			= make(chan IOp, param.MaxWait)
	engine.chStop		= make(chan bool, 1)
	engine.running		= true
	
	// 处理操作
	go func() {
		for op := range engine.chOp {
			op.Exec(engine)
		}
		
		// 通知操作完成
		engine.chStop <- true
	}()
	
	return nil
}

// stop
func (engine *CbEngine) Stop() {
	// 
	if !engine.running { return }
	
	// 关闭操作队列
	close(engine.chOp)
	
	// 等待操作完成
	<- engine.chStop
	
	// 清理
	engine.bucket.Close()
	engine.running = false
}

// 添加操作
func (engine *CbEngine) AddOp(op IOp) bool {
	select {
	case engine.chOp <- op:
		return true
	default:
		log.Println("too many db op to be done!")
		return false
	}
}

// url
func url(param *ParamInit) string {
	return "http://" + param.Bucket + ":" + param.Pwd + "@" + param.IP + "/"
}
