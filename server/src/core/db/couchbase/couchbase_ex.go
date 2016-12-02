// couchbase封装扩展
// 提供异步操作结果channel统一管理调用
package couchbase

// 
import (
	"log"
	
//	gcb "github.com/couchbaselabs/go-couchbase"
)

// db操作
type IOpEx interface {
	Exec(engine *CbEngineEx)
}

// 结果
type IResult interface {
	Handle()
}

// 
type CbEngineEx struct {
	*CbEngine						// bucket
	chOp			chan IOpEx		// 异步操作队列
	chResult		chan IResult	// 异步操作结果队列
	chStop			chan bool		// 结束通知(非阻塞)
}

// 创建对象
func NewCbEngineEx() *CbEngineEx {
	return &CbEngineEx{ NewCbEngine(), nil, nil, nil }
}

// 获取名字
func (engine *CbEngineEx) Name() string {
	return "CbEngineEx"
}

// run
func (engine *CbEngineEx) Serve(param *ParamInit) (err error) {
	// 启动CbEngine
	if err := engine.CbEngine.Serve(param); err != nil {
		return err
	}
	
	// 初始化其他属性
	engine.chOp			= make(chan IOpEx, param.MaxWait)
	engine.chResult		= make(chan IResult, param.MaxWait * 2)
	engine.chStop		= make(chan bool, 1)
	
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
func (engine *CbEngineEx) Stop() {
	// 
	if !engine.running { return }
	
	// 关闭操作队列
	close(engine.chOp)
	
	// 等待操作完成
	<- engine.chStop
	
	// 关闭CbEngine
	engine.CbEngine.Stop()
}

// 处理异步操作结果
func (engine *CbEngineEx) ProcAsynRes() {
	for {
		select {
		case res := <-engine.chResult:
			res.Handle()
		default:
			return
		}
	}
}

// 添加操作
func (engine *CbEngineEx) AddOpEx(op IOpEx) bool {
	select {
	case engine.chOp <- op:
		return true
	default:
		log.Println("too many db op to be done!")
		return false
	}
}

// 反馈通知
func (engine *CbEngineEx) AddChRes(chRes IResult) {
	engine.chResult <- chRes
}
