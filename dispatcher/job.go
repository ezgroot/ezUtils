package dispatcher

import (
	"context"
	"time"
)

// InData of Task.
type InData struct {
	TaskUID string
	Data    []byte
	If      interface{}
}

// OutData of Task.
type OutData struct {
	Code int
	Err  error
	Data []byte
	If   interface{}
}

// Task that worker to operate.
type Task struct {
	OriginalCtx context.Context
	InData      InData
	OutCh       chan OutData
	TimeOut     time.Duration // per - Millisecond
	StartTime   int64
	EndTime     int64
}
