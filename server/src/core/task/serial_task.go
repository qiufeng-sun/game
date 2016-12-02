// 串行任务接口及调度: 执行后会清空任务
package task

//
import (
	"core/safe/list"
)

//
type SerialTasks struct {
	slst *slist.SafeList // 任务队列. 执行后会被清空
}

// new
func NewSerialTasks() *SerialTasks {
	return &SerialTasks{slist.New()}
}

//
func (st *SerialTasks) AddTask(f func()) {
	st.slst.PushBack(f)
}

// 执行
func (st *SerialTasks) Exec() {
	// 遍历任务, 并执行
	f := st.slst.Begin()
	for ; f != nil; f = f.Next() {
		f.Value.(func())()
	}
	st.slst.End()

	// 清空任务
	st.removeAll()
}

//
func (st *SerialTasks) removeAll() {
	st.slst.Clear()
}

// 默认串行任务管理器
var g_sTasks *SerialTasks

// 初始化
func SInit() error {
	g_sTasks = NewSerialTasks()

	return nil
}

// 添加任务
func SAddTask(f func()) {
	g_sTasks.AddTask(f)
}

// 执行
func SExec() {
	g_sTasks.Exec()
}
