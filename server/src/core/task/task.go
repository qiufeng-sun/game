// 并行任务接口及调度：分长期任务和临时任务。临时任务执行后，会被清除。
package task

//
import (
	"sync"
)

// 任务接口
type ITask interface {
	Exec() // 注意: 任务结束后需要调用TaskOver()
}

//
type TaskFunc func()

func (this TaskFunc) Exec() {
	this()
}

// 任务调度器
type ParallelTasks struct {
	task    []ITask    // 长期任务(所有任务结束后不清除)
	tmpTask []ITask    // 临时任务(所有任务结束后清除)
	chTask  chan ITask // 任务通讯channel

	waitTask sync.WaitGroup // waiter -- 等待所有任务完成
	waitStop sync.WaitGroup // waiter -- 等待所有处理协程结束
}

// 创建任务调度器
func NewParallelTasks() *ParallelTasks {
	return &ParallelTasks{}
}

// 启动任务调度器
func (m *ParallelTasks) Serve(goNum int) {
	// 初始化协程相关
	m.waitStop.Add(goNum)
	m.chTask = make(chan ITask, goNum*2)

	// 启动处理协程
	for i := 0; i < goNum; i++ {
		go func() {
			for task := range m.chTask {
				task.Exec()
				m.waitTask.Done()
			}

			m.waitStop.Done()
		}()
	}
}

// 关闭任务调度器
func (m *ParallelTasks) Stop() {
	close(m.chTask)
	m.waitStop.Wait()
}

// 添加长期任务
func (m *ParallelTasks) AddTask(t ITask) {
	m.task = append(m.task, t)
}

// 删除长期任务
func (m *ParallelTasks) RemoveTask(t ITask) {
	for i, task := range m.task {
		if task == t {
			// 不是最后一个
			end := len(m.task) - 1
			if i != end {
				m.task[i] = m.task[end]
			}
			m.task = m.task[:end]

			return
		}
	}
}

// 添加临时任务
func (m *ParallelTasks) AddTmpTask(t ITask) {
	m.tmpTask = append(m.tmpTask, t)
}

// 删除所有临时任务
func (m *ParallelTasks) removeAllTmpTask() {
	m.tmpTask = m.tmpTask[:0]
}

// 执行任务
func (m *ParallelTasks) Exec() {
	var num = len(m.task) + len(m.tmpTask)
	if 0 == num {
		return
	}

	m.waitTask.Add(num)
	defer m.waitTask.Wait()

	// 永久任务
	for _, t := range m.task {
		m.chTask <- t
	}

	// 临时任务
	if len(m.tmpTask) > 0 {
		for _, t := range m.tmpTask {
			m.chTask <- t
		}

		// 清除临时任务
		m.removeAllTmpTask()
	}
}

// 默认并行任务管理器
var g_pTasks *ParallelTasks

// 启动默认任务调度器
func PInit(goNum int) error {
	g_pTasks = NewParallelTasks()

	g_pTasks.Serve(goNum)

	return nil
}

// 添加长期任务
func PAddTask(t ITask) {
	g_pTasks.AddTask(t)
}

func PAddTaskFunc(f func()) {
	PAddTask(TaskFunc(f))
}

// 删除长期任务
func PRemoveTask(t ITask) {
	g_pTasks.RemoveTask(t)
}

// 添加临时任务
func PAddTmpTask(t ITask) {
	g_pTasks.AddTmpTask(t)
}

// 执行任务
func PExec() {
	g_pTasks.Exec()
}
