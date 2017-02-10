//
package world

//
import (
	"runtime"

	"util/logs"

	"core/net/msg/protobuf"
	"core/net/socket"
	"core/task"

	"world/db"
	"world/define"
	"world/player/logon"
	"world/player/session"
	"world/res"

	"world/msg/hfunc"
)

// world server obj
type WorldServer struct {
}

func New() *WorldServer {
	return &WorldServer{}
}

type TaskTest struct {
}

func (t *TaskTest) Exec() {
	logs.Infoln("task test")
}

// 检查error
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

//
func (w *WorldServer) Init() bool {
	// 配置及资源读取
	if !res.LoadAll() {
		return false
	}

	var e error

	// 串行任务初始化
	e = task.SInit()
	checkErr(e)

	// 并行任务调度启动
	e = task.PInit(runtime.NumCPU() + 2)
	checkErr(e)

	// 临时 -- 测试
	//task.PAddTmpTask(&TaskTest{})

	// game db启动
	db.Init(res.Path_Cfg)

	// 注册消息处理函数
	hfunc.Register()

	// 网络层启动
	e = socket.Serve(define.GetPlayerConn().LstAddr,
		define.GetPlayerConn().MaxNum, &protobuf.PbParser{})
	checkErr(e)

	return true
}

//
func (w *WorldServer) Update() {
	// session更新
	session.Update()

	// 并发任务执行
	task.PExec()

	// 更新db反馈// to do
	//db.DB().ProcAsynRes()

	// 串行任务执行
	task.SExec()

	// 处理被踢玩家
	session.ProcDisconnect()

	// 未登入world玩家更新(已收到第一条消息)
	logon.Update()
}

//
func (w *WorldServer) Destroy() {
}

//
func (w *WorldServer) PreQuit() {
	// 关闭网络层
	socket.Stop()

	// 关闭db// to do
	//db.DB().Stop()
}

//
func (w *WorldServer) String() string {
	return "WorldServer"
}
