package gobq

import (
	"fmt"
	"sync"
)

type Row map[string]interface{}

type BufferedWrite struct {
	sync.RWMutex
	Rows     []Row
	capacity int
}

func (b *BufferedWrite) Append(r Row) error {
	if b.IsFull() {
		return fmt.Errorf("cannot append to buffer, capacity exceeded: %d elements", len(b.rows))
	}
	b.Lock()
	defer b.Unlock()
	b.rows = append(b.Rows, r)
	return nil
}

func (b *BufferedWrite) Length() int {
	return len(b.rows)
}

func (b *BufferedWrite) Capacity() int {
	return b.capacity
}

func (b *BufferedWrite) IsFull() bool {
	return len(b.rows) == b.capacity
}

func (b *BufferedWrite) Reset() {
	b.Lock()
	defer b.Unlock()
	b.rows = make([]Row, 0, b.capacity)
}

func (b *BufferedWrite) Flush() []Row {
	defer b.Reset()
	return b.rows
}

func NewBufferedWrite(c int) *BufferedWrite {
	b := &BufferedWrite{capacity: c}
	b.Reset()
	return b
}
