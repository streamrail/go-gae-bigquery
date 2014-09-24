package tracking

import (
	"fmt"
	"sync"
)

type BufferedWrite struct {
	sync.RWMutex
	Rows     []map[string]interface{}
	Capacity int
}

func (b *BufferedWrite) Append(row map[string]interface{}) error {
	b.Lock()
	defer b.Unlock()
	if len(b.Rows) == b.Capacity {
		return fmt.Errorf("cannot append to buffer, capacity exceeded: %d elements", len(b.Rows))
	}
	b.Rows = append(b.Rows, row)
	return nil
}

func (b *BufferedWrite) Length() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.Rows)
}

func (b *BufferedWrite) IsFull() bool {
	len := b.Length()
	return len == b.Capacity
}

func (b *BufferedWrite) Reset() {
	b.Lock()
	defer b.Unlock()
	b.Rows = b.Rows[:0]
}

func (b *BufferedWrite) Flush() []map[string]interface{} {
	b.RLock()
	defer b.Reset()
	defer b.RUnlock()
	return b.Rows
}

func NewBufferedWrite(c int) *BufferedWrite {
	return &BufferedWrite{
		Rows:     *new([]map[string]interface{}),
		Capacity: c,
	}
}
