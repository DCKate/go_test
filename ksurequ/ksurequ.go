package ksurequ

import (
	"sync"
)

type Sque struct {
	Slen int
	Qu   []interface{}
	lock sync.Mutex
}

// SetLen only use with function PushSureQueue
func (q *Sque) SetLen(l int) {
	q.Slen = l
}

// PushSureQueue make sure pop first element only when Queue is full
func (q *Sque) PushSureQueue(m interface{}) (int, interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Slen == len(q.Qu) {
		x := q.Qu[0]
		q.Qu = append(q.Qu[1:], m)
		return 1, x
	}
	q.Qu = append(q.Qu, m)
	return -1, nil

}

func (q *Sque) popMamber() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.Qu) > 0 {
		x := q.Qu[0]
		q.Qu = q.Qu[1:]
		return x
	}
	return nil
}

func (q *Sque) pushMamber(m interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.Qu = append(q.Qu, m)
}
