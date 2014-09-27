package gobq

import (
	"testing"
)

func TestAppend(t *testing.T) {
	b := NewBufferedWrite(10)
	r := Row{
		"Category": "Test",
	}
	b.Append(r)
	if b.Length() != 1 {
		t.Error("length should be 1 after appending 1 row")
	}
	f := b.Flush()
	if f[0]["Category"] != r["Category"] {
		t.Error("flushed row should be equal to appended row")
	}
}

func TestLength(t *testing.T) {
	b := NewBufferedWrite(10)
	if b.Length() != 0 {
		t.Error("length should be 0 after the buffer is initiated")
	}
	r := Row{
		"Category": "Test",
	}
	b.Append(r)
	if b.Length() != 1 {
		t.Error("length should be 1 after appending 1 row")
	}
	b.Flush()
	if b.Length() != 0 {
		t.Error("length should be 0 after the buffer is flushed")
	}
}

func TestCapacity(t *testing.T) {
	b := NewBufferedWrite(1)
	r := Row{
		"Category": "Test",
	}
	b.Append(r)
	if err := b.Append(r); err == nil {
		t.Error("Append should return an error when we try to append beyond the buffer's capacity")
	}
}

func TestIsFull(t *testing.T) {
	b := NewBufferedWrite(2)
	r := Row{
		"Category": "Test",
	}
	if b.IsFull() {
		t.Error("Buffer should not be full when the capacity had not been reached")
	}
	b.Append(r)
	if b.IsFull() {
		t.Error("Buffer should not be full when the capacity had not been reached")
	}
	b.Append(r)
	if !b.IsFull() {
		t.Error("Buffer should be full when the capacity had been reached")
	}
}

func TestReset(t *testing.T) {
	b := NewBufferedWrite(2)
	r := Row{
		"Category": "Test",
	}
	b.Append(r)
	b.Reset()
	if b.Length() != 0 {
		t.Error("length should be 0 after the buffer was reset")
	}
}

func TestFlush(t *testing.T) {
	b := NewBufferedWrite(10)
	r := Row{
		"Category": "Test",
	}
	b.Append(r)
	f := b.Flush()
	if f[0]["Category"] != r["Category"] {
		t.Error("flushed row should be equal to appended row")
	}
	if b.Length() != 0 {
		t.Error("length should be 0 after the buffer was flushed")
	}
}
