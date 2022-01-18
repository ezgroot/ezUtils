package dispatcher

import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	defaultTaskQueueMaxNum = 100 * 10000
)

// ErrorOfDispatcherHasStop dispatcher is stop
var ErrorOfDispatcherHasStop = errors.New("dispatcher has stop")

// State Dispatcher state
type State int

// state enum
const (
	StateOfNormal State = iota
	StateOfHalf
	StateOfBusy
	StateOfDeny
)

// Config of dispatcher.
type Config struct {
	Name         string
	MaxWorkerNum int
	MaxTaskNum   int
	IsTrance     bool
}

// Dispatcher job to worker to do.
type Dispatcher struct {
	maxWorkerNum     int
	workerCh         chan *Worker
	taskCh           chan Task
	name             string
	state            State
	isStop           bool
	isTrance         bool
	totalInTask      uint64
	totalDoneTask    uint64
	totalErrorTask   uint64
	totalRefusedTask uint64
}

func newDispatcher(c Config) *Dispatcher {
	if c.MaxWorkerNum <= 0 {
		c.MaxWorkerNum = runtime.NumCPU() * 2
	}

	if c.MaxTaskNum <= 0 {
		c.MaxTaskNum = defaultTaskQueueMaxNum
	}

	d := &Dispatcher{name: c.Name,
		maxWorkerNum: c.MaxWorkerNum,
		workerCh:     make(chan *Worker, c.MaxWorkerNum),
		taskCh:       make(chan Task, c.MaxTaskNum),
		isTrance:     c.IsTrance,
		state:        StateOfNormal,
		isStop:       false}

	return d
}

func (d *Dispatcher) assignTask() {
	for t := range d.taskCh {
		go func(task Task) {
			worker := <-d.workerCh
			worker.taskCh <- task
		}(t)
	}
}

func (d *Dispatcher) run(workDo workDoing, r Resource) {
	for i := 0; i < d.maxWorkerNum; i++ {
		worker := newWorker(d, i, workDo, r)
		worker.working()

		d.workerCh <- worker
	}

	go d.assignTask()
}

func (d *Dispatcher) setStop() {
	d.isStop = true
}

func (d *Dispatcher) setStart() {
	d.isStop = false
}

func (d *Dispatcher) getStopFlag() bool {
	return d.isStop
}

// AddTask add a job to do.
func (d *Dispatcher) AddTask(t Task) error {
	if d.getStopFlag() {
		return ErrorOfDispatcherHasStop
	}

	if d.GetCurrentTaskTodoNum() >= d.GetMaxTaskNum()/2 {
		d.state = StateOfHalf
	}

	if d.GetCurrentTaskTodoNum() >= d.GetMaxTaskNum()*7/10 {
		d.state = StateOfBusy
	}

	if d.GetCurrentTaskTodoNum() >= d.GetMaxTaskNum()*95/100 {
		atomic.AddUint64(&(d.totalRefusedTask), 1)
		d.state = StateOfDeny
		d.Status()
		return fmt.Errorf("dispatcher too busy, task = %s has been refused", t.InData.TaskUID)
	}

	t.StartTime = time.Now().UTC().UnixNano()

	atomic.AddUint64(&(d.totalInTask), 1)

	d.taskCh <- t

	return nil
}

// StopTemporary dispatcher, not accept task any more, but should finish all task that in queue,
// can start again by call StartAgain().
func (d *Dispatcher) StopTemporary() {
	d.setStop()
}

// StartAgain dispatcher.
func (d *Dispatcher) StartAgain() {
	d.setStart()
}

// StopForever dispatcher, not accept task any more, but should finish all task that in queue,
// then every worker will quit, can't be start again.
func (d *Dispatcher) StopForever(done chan struct{}) {
	d.setStop()

	go func() {
		for {
			if d.GetCurrentTaskTodoNum() == 0 {
				if d.GetCurrentFreeWorkerNum() == 0 {
					fmt.Printf("All workers stop accepting new task and continue to finish the work on their hand now.\n")
					done <- struct{}{}
					return
				}

				worker := <-d.workerCh
				worker.stop()
			}
		}
	}()
}

// Status show the status of the current dispatcher.
func (d *Dispatcher) Status() {
	fmt.Printf(
		"\n****************************************************************\n"+
			"*               %s dispatcher status\n"+
			"*         ----------------------------------------------\n"+
			"* state : %d\n"+
			"* isStop : %t\n"+
			"* total worker number : %d\n"+
			"* max task number : %d\n"+
			"* current free workers : %d\n"+
			"* current task wait todo : %d\n"+
			"* total in task number : %d\n"+
			"* total done task number : %d\n"+
			"* total error task number : %d\n"+
			"* total refused task number : %d\n"+
			"****************************************************************\n",
		d.GetName(), d.state, d.getStopFlag(), d.GetTotalWorkNum(), d.GetMaxTaskNum(),
		d.GetCurrentFreeWorkerNum(), d.GetCurrentTaskTodoNum(), d.GetTotalInTask(),
		d.GetTotalDoneTask(), d.GetTotalErrorTask(), d.GetTotalRefusedTask())
}

// GetName get dispatcher name.
func (d *Dispatcher) GetName() string {
	return d.name
}

// GetTotalWorkNum get max number of worker.
func (d *Dispatcher) GetTotalWorkNum() int {
	return cap(d.workerCh)
}

// GetMaxTaskNum get max number of job queue.
func (d *Dispatcher) GetMaxTaskNum() int {
	return cap(d.taskCh)
}

// GetCurrentFreeWorkerNum get current free worker number.
func (d *Dispatcher) GetCurrentFreeWorkerNum() int {
	return len(d.workerCh)
}

// GetCurrentTaskTodoNum get current job number wait to operate.
func (d *Dispatcher) GetCurrentTaskTodoNum() int {
	return len(d.taskCh)
}

// GetTotalInTask  get total number of job which put in job queue since dispatcher start.
func (d *Dispatcher) GetTotalInTask() uint64 {
	return atomic.LoadUint64(&(d.totalInTask))
}

// GetTotalDoneTask get total number of job which all workers finish since dispatcher start.
func (d *Dispatcher) GetTotalDoneTask() uint64 {
	return atomic.LoadUint64(&(d.totalDoneTask))
}

// GetTotalErrorTask get total number of job which all workers operate error since dispatcher start.
func (d *Dispatcher) GetTotalErrorTask() uint64 {
	return atomic.LoadUint64(&(d.totalErrorTask))
}

// GetTotalRefusedTask get total number of job which refuse because of job queue full since dispatcher start.
func (d *Dispatcher) GetTotalRefusedTask() uint64 {
	return atomic.LoadUint64(&(d.totalRefusedTask))
}

// GetState get dispatcher state.
func (d *Dispatcher) GetState() State {
	return d.state
}

// IsStop is dispatcher has stop.
func (d *Dispatcher) IsStop() bool {
	return d.getStopFlag()
}

// GetDispatch Get a worker task queue model.
func GetDispatch(c Config, workDo workDoing, r Resource) *Dispatcher {
	dispatcher := newDispatcher(c)
	dispatcher.run(workDo, r)

	return dispatcher
}
