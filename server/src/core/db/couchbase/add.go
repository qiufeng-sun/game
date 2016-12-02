// 新增db记录(like sql insert)
package couchbase

// 
import (
//	"log"
//	"errors"
	
)

//// 序列化
//type IMarshaler interface {
//	Marshal() ([]byte, error)
//}


// Add: 增加新记录. 存在或者失败则返回false

// 不关心操作结果
type OpAdd0 struct {
	key			string
	val			[]byte
}

func NewOpAdd0(key string, val []byte) *OpAdd0 {
	return &OpAdd0{key, val}
}

func (op *OpAdd0) Exec(engine *CbEngine) {
	engine.AddRaw(op.key, 0, op.val)
}

func Add0(key string, val []byte) bool {
	return AddOp(NewOpAdd0(key, val))
}

// 新增记录结果
type ResAdd struct {
	added		bool
	err			error
}

func NewResAdd(added bool, err error) *ResAdd {
	return &ResAdd{added, err}
}

func (r *ResAdd) Ok() bool		{ return r.added }
func (r *ResAdd) Exist() bool	{ return !r.added && nil == r.err }
func (r *ResAdd) Error() error	{ return r.err }

// 关心操作结果
type OpAdd1 struct {
	*OpAdd0
	chRes		chan<- *ResAdd
}

func NewOpAdd1(key string, val []byte, chRes chan<- *ResAdd) *OpAdd1 {
	return &OpAdd1{NewOpAdd0(key, val), chRes}
}

func (op *OpAdd1) Exec(engine *CbEngine) {
	op.chRes <- NewResAdd(engine.AddRaw(op.key, 0, op.val))
}

func Add1(key string, val []byte, chRes chan<- *ResAdd) bool {
	return AddOp(NewOpAdd1(key, val, chRes))
}
