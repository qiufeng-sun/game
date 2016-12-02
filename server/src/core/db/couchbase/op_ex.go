// db操作
package couchbase

// 
import (
//	"log"
//	"errors"
)

//// 通过实现IHResInsert接口, 可以在couchbase.ProcAsynRes()中处理反馈
//// 例: 获取数据
//	func (o *Item) HandleInsert() {
//		// 处理结果
//	}
//	var item Item
//	couchbase.InsertAnsyEx(key, val, item)

// 处理创建记录结果接口
type IHResInsert interface {
	HandleInsert(data *OpInsertData, res *ResInsert)
}

// 函数接口
type FuncHResInsert	func(data *OpInsertData, res *ResInsert)
func (f FuncHResInsert) HandleInsert(data *OpInsertData, res *ResInsert) {
	f(data, res)
}

// 增加新记录
type OpInsertEx struct {
	OpInsertData
	*ResInsert
	handle		IHResInsert
}

func NewOpInsertEx(key string, val []byte, info interface{}, handle IHResInsert) *OpInsertEx {
	return &OpInsertEx{OpInsertData{key, val, info}, nil, handle}
}

func (op *OpInsertEx) Exec(engine *CbEngineEx) {
	added, err := engine.InsertSync(op.Key, op.Val)
	
	if op.handle != nil {
		op.ResInsert = NewResInsert(added, err)
		engine.AddChRes(op)
	}
}

func (op *OpInsertEx) Handle() {
	if op.handle != nil {
		op.handle.HandleInsert(&op.OpInsertData, op.ResInsert)
	} else {
		panic("logical error! OpInsertEx.Handle().")
	}
}

// 异步添加新记录. 不需要反馈时, 传入handle为nil
func (engine *CbEngineEx) InsertAsynEx(key string, val []byte, info interface{}, handle IHResInsert) bool {
	return engine.AddOpEx(NewOpInsertEx(key, val, info, handle))
}

func (engine *CbEngineEx) InsertAsynExFunc(key string, val []byte, info interface{}, f func(data *OpInsertData, res *ResInsert)) bool {
	return engine.InsertAsynEx(key, val, info, FuncHResInsert(f))
}

//// 例:
//// 通过实现IHResGet接口, 可以在couchbase.ProcAsynRes()中处理反馈
//// 例: 获取数据
//	func (o *Item) HandleGet(res *ResGet) {
//		// 处理结果
//	}
//	var item Item
//	couchbase.GetAsynEx(key, item)

// 获取记录结果后处理接口
type IHResGet interface {
	HandleGet(data *OpGetData, res *ResGet)
}

// 函数定义
type FuncHResGet	func(data *OpGetData, res *ResGet)
func (f FuncHResGet) HandleGet(data *OpGetData, res *ResGet) {
	f(data, res)
}

// 获取记录
type OpGetEx struct {
	OpGetData
	*ResGet
	handle		IHResGet
}

func NewOpGetEx(key string, info interface{}, handle IHResGet) *OpGetEx {
	return &OpGetEx{OpGetData{key, info}, nil, handle}
}

func (op *OpGetEx) Exec(engine *CbEngineEx) {
	res, err := engine.GetSync(op.Key)
	
	if op.handle != nil {
		op.ResGet = NewResGet(res, err)
		engine.AddChRes(op)
	} 
}

func (op *OpGetEx) Handle() {
	if op.handle != nil {
		op.handle.HandleGet(&op.OpGetData, op.ResGet)
	} else {
		panic("logical error! OpGetEx.Handle()")
	}
}

func (engine *CbEngineEx) GetAsynEx(key string, info interface{}, handle IHResGet) bool {
	return engine.AddOpEx(NewOpGetEx(key, info, handle))
}

func (engine *CbEngineEx) GetAsynExFunc(key string, info interface{}, f func(data *OpGetData, res *ResGet)) bool {
	return engine.GetAsynEx(key, info, FuncHResGet(f))
}

//// 通过实现IHResSet接口, 可以在couchbase.ProcAsynRes()中处理反馈
//// 例: 获取数据
//	func (o *Item) HandleSet() {
//		// 处理结果
//	}
//	var item Item
//	couchbase.SetAysnEx(key, val, item)

// 处理创建记录结果接口
type IHResSet interface {
	HandleSet(data *OpSetData, res *ResSet)
}

// 函数接口
type FuncHResSet	func(data *OpSetData, res *ResSet)
func (f FuncHResSet) HandleSet(data *OpSetData, res *ResSet) {
	f(data, res)
}

// 更新记录
type OpSetEx struct {
	OpSetData
	*ResSet
	handle		IHResSet
}

func NewOpSetEx(key string, val []byte, info interface{}, handle IHResSet) *OpSetEx {
	return &OpSetEx{OpSetData{key, val, info}, nil, handle}
}

func (op *OpSetEx) Exec(engine *CbEngineEx) {
	err := engine.SetSync(op.Key, op.Val)
	
	if op.handle != nil {
		op.ResSet = NewResSet(err)
		engine.AddChRes(op)
	}
}

func (op *OpSetEx) Handle() {
	if op.handle != nil {
		op.handle.HandleSet(&op.OpSetData, op.ResSet)
	} else {
		panic("logical error! OpSetEx.Handle().")
	}
}

// 异步更新记录. 不需要反馈时, 传入handle为nil
func (engine *CbEngineEx) SetAsynEx(key string, val []byte, info interface{}, handle IHResSet) bool {
	return engine.AddOpEx(NewOpSetEx(key, val, info, handle))
}

func (engine *CbEngineEx) SetAsynExFunc(key string, val []byte, info interface{}, f func(data *OpSetData, res *ResSet)) bool {
	return engine.SetAsynExFunc(key, val, info, FuncHResSet(f))
}
