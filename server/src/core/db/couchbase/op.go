// db操作
package couchbase

// 
import (
//	"log"
//	"errors"
)

//// 异步添加记录(sql insert)
//	if op, added := couchbase.InsertAsyn(k, val, true); !added {
//		// 异步队列满了...
//	} else {
//		do something other......
		
//		// 等待db操作结果
//		res := <-op.ChRes

//		// 处理res
//	}

// 同步添加新记录(sql insert)
func (engine *CbEngine) InsertSync(key string, val []byte) (bool, error) {
	return engine.bucket.AddRaw(key, 0, val)
}

// 新增记录结果
type ResInsert struct {
	added		bool
	err			error
}

func NewResInsert(added bool, err error) *ResInsert {
	return &ResInsert{added, err}
}

func (r *ResInsert) Ok() bool		{ return r.added }
func (r *ResInsert) Exist() bool	{ return !r.added && nil == r.err }
func (r *ResInsert) Error() error	{ return r.err }

// OpInsertData
type OpInsertData struct {
	Key			string
	Val			[]byte
	Info		interface{}			// 缓存上层信息
}

// OpInsert 增加新记录. 存在或者失败则返回false
type OpInsert struct {
	OpInsertData
	ChRes		chan *ResInsert
}

func NewOpInsert(key string, val []byte, info interface{}, chRes chan *ResInsert) *OpInsert {
	return &OpInsert{OpInsertData{key, val, info}, chRes}
}

func (op *OpInsert) Exec(engine *CbEngine) {
	added, err := engine.InsertSync(op.Key, op.Val)
	
	if op.ChRes != nil {
		op.ChRes <- NewResInsert(added, err)
	}
}

// 异步添加新记录
func (engine *CbEngine) InsertAsyn(key string, val []byte, info interface{}, needFeedback bool) (*OpInsert, bool) {
	var chRes chan *ResInsert
	if needFeedback {
		chRes = make(chan *ResInsert, 1)		// 非阻塞
	}
	
	var op = NewOpInsert(key, val, info, chRes)
	
	return op, engine.AddOp(op)
}

//// 异步获取记录(sql select)
//	if op, added := couchbase.GetAsyn(key); !added {
//		// 异步队列满了
//	} else {
//		do something other......
		
//		// 等待db操作结果
//		res := <-op.ChRes

//		// 处理res
//	}

// 同步获取数据
func (engine *CbEngine) GetSync(key string) ([]byte, error) {
	return engine.bucket.GetRaw(key)
}

// 获取数据结果
type ResGet struct {
	Res			[]byte
	Err			error
}

func NewResGet(res []byte, err error) *ResGet {
	return &ResGet{res, err}
}

// OpGetData
type OpGetData struct {
	Key			string
	Info		interface{}			// 缓存上层信息
}

// OpGet 获取数据
type OpGet struct {
	OpGetData
	ChRes		chan *ResGet
}

func NewOpGet(key string, info interface{}, chRes chan *ResGet) *OpGet {
//	if nil == chRes { panic("OpGet: chRes cannt be nil!") }
	
	return &OpGet{OpGetData{key, info}, chRes}
}

func (op *OpGet) Exec(engine *CbEngine) {
	op.ChRes <- NewResGet(engine.GetSync(op.Key))
}

// 异步获取数据
func (engine *CbEngine) GetAsyn(key string, info interface{}) (*OpGet, bool) {
	// 生成op
	op := NewOpGet(key, info, make(chan *ResGet, 1))
	
	return op, engine.AddOp(op)
}

//// 异步更新数据(sql replace)
//	if op, added := couchbase.SetAsyn(k, val, true); !added {
//		// 异步队列满了...
//	} else {
//		do something other......
		
//		// 等待db操作结果
//		res := <-op.ChRes

//		// 处理res
//	}

// 同步更新数据
func (engine *CbEngine) SetSync(key string, val []byte) error {
	return engine.bucket.SetRaw(key, 0, val)
}

// 更新记录
type ResSet struct {
	Err			error
}

func NewResSet(err error) *ResSet {
	return &ResSet{err}
}

// OpSetData
type OpSetData struct {
	Key			string
	Val			[]byte
	Info		interface{}			// 缓存上层信息
}

// OpSet 更新记录. 不存在, 则创建新的记录
type OpSet struct {
	OpSetData
	ChRes		chan *ResSet
}

func NewOpSet(key string, val []byte, info interface{}, chRes chan *ResSet) *OpSet {
	return &OpSet{OpSetData{key, val, info}, chRes}
}

func (op *OpSet) Exec(engine *CbEngine) {
	err := engine.SetSync(op.Key, op.Val)
	
	if op.ChRes != nil {
		op.ChRes <- NewResSet(err)
	} 
}

// 异步添加新记录
func (engine *CbEngine) SetAsyn(key string, val []byte, info interface{}, needFeedback bool) (*OpSet, bool) {
	var chRes chan *ResSet
	if needFeedback {
		chRes = make(chan *ResSet, 1)
	}
	
	var op = NewOpSet(key, val, info, chRes)
	
	return op, engine.AddOp(op)
}




