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
		return fmt.Errorf("cannot append to buffer, capacity exceeded: %d elements", len(b.Rows))
	}
	b.Lock()
	defer b.Unlock()
	b.Rows = append(b.Rows, r)
	return nil
}

func (b *BufferedWrite) Length() int {
	return len(b.Rows)
}

func (b *BufferedWrite) Capacity() int {
	return b.capacity
}

func (b *BufferedWrite) IsFull() bool {
	return len(b.Rows) == b.capacity
}

func (b *BufferedWrite) Reset() {
	b.Lock()
	defer b.Unlock()
	b.Rows = make([]Row, 0, b.capacity)
}

func (b *BufferedWrite) Flush() []Row {
	defer b.Reset()
	return b.Rows
}

func NewBufferedWrite(c int) *BufferedWrite {
	b := &BufferedWrite{capacity: c}
	b.Reset()
	return b
}
