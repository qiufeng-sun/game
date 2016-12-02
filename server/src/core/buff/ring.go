// 环形缓冲
package buff

//
import (
	"errors"
	"io"
	"sync/atomic"
)

// 错误码
var (
	ErrWriteNone = errors.New("RingBuff: write none")
	ErrFull      = errors.New("RingBuff: full")
)

// io vector
type IoVector struct {
	Buff0 []byte // read_head
	Buff1 []byte // root
}

// 大小
func (iov *IoVector) Size0() int { return len(iov.Buff0) }
func (iov *IoVector) Size1() int { return len(iov.Buff1) }
func (iov *IoVector) Size() int  { return len(iov.Buff0) + len(iov.Buff1) }

// 获取连续的byte切片
func (iov *IoVector) Bytes(sz, start int) []byte {
	//
	var end = start + sz

	// 数据全部在iov.Buff0中
	if iov.Size0() >= end {
		return iov.Buff0[start:end]
	}

	// 大小不足
	if iov.Size() < end {
		return nil
	}

	// 数据全部在iov.Buff1中
	if iov.Size0() <= start {
		offset := start - iov.Size0()
		return iov.Buff1[offset : offset+sz]
	}

	// 两头 -- 申请新的空间
	var ret = make([]byte, sz)
	var n = copy(ret, iov.Buff0[start:])
	copy(ret[n:], iov.Buff1)

	return ret
}

// RingBuff 环形缓冲 -- 读写可在不同线程同时进行, 但若上层多线程同时读, 或
// 同时写需另作互斥
type RingBuff struct {
	buff  []byte // 消息缓冲区
	rPos  int    // 上次读取位置
	wPos  int    // 底层接收位置
	bytes int32  // 已缓冲数据大小 -- 这个需要原子操作
}

// 创建
func NewRingBuff(sz int) *RingBuff {
	return &RingBuff{buff: make([]byte, sz)}
}

// 获取已缓冲数据大小 -- 原子操作
func (r *RingBuff) GetBuffed() int {
	return int(atomic.LoadInt32(&r.bytes))
}

// 获取现有空闲空间
func (r *RingBuff) GetFreed() int {
	return len(r.buff) - r.GetBuffed()
}

// 清空
func (r *RingBuff) Clear() {
	r.rPos = 0
	r.wPos = 0
	r.bytes = 0
}

// 写入缓冲区 -- 注: 缓冲满后, 需要上层自己处理
func (r *RingBuff) ReadFrom(reader io.Reader) (int64, error) {
	// 已缓冲大小
	var buffed = r.GetBuffed()

	// 可写入大小
	var space = len(r.buff) - buffed

	// 缓冲满
	if 0 == space {
		return 0, ErrFull
	}

	// 计算可连续写入缓冲大小
	var endPos = len(r.buff)
	if buffed > r.wPos {
		endPos = r.wPos + space
	}

	// 写入缓冲
	readLen, err := reader.Read(r.buff[r.wPos:endPos])

	// 读取错误
	if err != nil {
		return 0, err
	}

	// 没有读取到东西
	if 0 == readLen {
		return 0, ErrWriteNone
	}

	// 更新缓冲大小
	atomic.AddInt32(&r.bytes, int32(readLen))

	// 更新写入位置
	r.wPos += readLen
	if len(r.buff) == r.wPos {
		r.wPos = 0
	}

	return int64(readLen), nil
}

// 写入缓冲区 -- 注: 缓冲满后, 需要上层自己处理
func (r *RingBuff) Write(p []byte) (n int, err error) {
	// 待写入大小
	var needSz = len(p)
	if 0 == needSz {
		return 0, nil
	}

	// 已缓冲大小
	var buffed = r.GetBuffed()

	// 可写入大小
	var space = len(r.buff) - buffed

	// 缓冲满, 或写不下
	if 0 == space || needSz > space {
		return 0, ErrFull
	}

	// 计算可连续写入缓冲大小
	var endPos = len(r.buff)
	if buffed > r.wPos {
		endPos = r.wPos + space
	}

	// 写入缓冲
	var sz = copy(r.buff[r.wPos:endPos], p)

	r.wPos += sz
	if len(r.buff) == r.wPos {
		r.wPos = 0
	}

	if needSz > sz {
		r.wPos = copy(r.buff[:], p[sz:])
	}

	// 更新缓冲大小
	atomic.AddInt32(&r.bytes, int32(needSz))

	return needSz, nil
}

// 读取缓冲区 -- 获取缓冲的数据切片
func (r *RingBuff) GetBuffedIoVector(vct *IoVector) (*IoVector, int) {
	// 已缓冲大小
	var buffed = r.GetBuffed()

	// 空
	if 0 == buffed {
		return nil, 0
	}

	// 空闲区域
	var space = len(r.buff) - buffed

	// 数据为连续的
	if space >= r.rPos {
		vct.Buff0 = r.buff[r.rPos : r.rPos+buffed]
		vct.Buff1 = nil
	} else { // 在两头
		vct.Buff0 = r.buff[r.rPos:]
		vct.Buff1 = r.buff[0 : r.rPos-space]
	}

	return vct, buffed
}

// 读取缓冲区 -- 仅移动rPos位置
func (r *RingBuff) Release(bytes int) {
	// 已缓冲大小
	var buffed = r.GetBuffed()

	// 移动
	var moved = bytes
	if bytes > buffed {
		moved = buffed
	}

	// 设置读取位置
	r.rPos += moved

	// 修正
	if r.rPos >= len(r.buff) {
		r.rPos -= len(r.buff)
	}

	// 更新缓冲大小
	atomic.AddInt32(&r.bytes, int32(-moved))
}
