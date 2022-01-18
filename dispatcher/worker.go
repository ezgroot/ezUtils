package dispatcher

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	defaultTaskTimeout = time.Duration(3000) * time.Millisecond
)

type workDoing func(ctx context.Context, res Resource, task Task) error

// Worker is who operate job task.
type Worker struct {
	dispatcher *Dispatcher
	workerID   int
	taskCh     chan Task
	workDo     workDoing
	quitCh     chan struct{}
	resource   Resource
}

func newWorker(dp *Dispatcher, id int, wd workDoing, res Resource) *Worker {
	w := &Worker{
		dispatcher: dp,
		workerID:   id,
		taskCh:     make(chan Task, 1),
		workDo:     wd,
		quitCh:     make(chan struct{}, 1),
		resource:   res,
	}

	return w
}

func (w *Worker) working() {
	go func() {
		for {
			select {
			case task := <-w.taskCh:
				{
					var err error
					var out OutData

					if task.TimeOut <= 0 {
						task.TimeOut = defaultTaskTimeout
					}

					var newCtx context.Context
					var cancel context.CancelFunc
					if task.OriginalCtx == nil {
						task.OriginalCtx = context.Background()
					}

					newCtx, cancel = context.WithTimeout(task.OriginalCtx, task.TimeOut*time.Millisecond)

					select {
					default:
					case <-newCtx.Done():
						{
							atomic.AddUint64(&(w.dispatcher.totalErrorTask), 1)

							out.Err = newCtx.Err()
							if task.OutCh != nil {
								task.OutCh <- out
							}
						}

						goto TaskEnd
					}

					if w.resource != nil {
						w.resource, err = w.resource.Check(newCtx)
						if err != nil {
							atomic.AddUint64(&(w.dispatcher.totalErrorTask), 1)

							out.Err = fmt.Errorf("worker resource check error : %s", err)
							if task.OutCh != nil {
								task.OutCh <- out
							}

							goto TaskEnd
						}
					}

					// workDo will make sure send result data to job out chan if need
					err = w.workDo(newCtx, w.resource, task)
					if err != nil {
						atomic.AddUint64(&(w.dispatcher.totalErrorTask), 1)
					} else {
						atomic.AddUint64(&(w.dispatcher.totalDoneTask), 1)
					}

				TaskEnd:
					task.EndTime = time.Now().UTC().UnixNano()
					if w.dispatcher.isTrance {
						fmt.Printf("from dispatcher, module: %s, workerID : %d, jobID : %s, cost %s\n",
							w.dispatcher.name, w.workerID, task.InData.TaskUID,
							diffTimeStrapToShow(task.StartTime, task.EndTime))
					}

					cancel()
				}
			case <-w.quitCh:
				{
					if w.resource != nil {
						err := w.resource.Release()
						if err != nil {
							fmt.Printf("release worker resource error = %s\n", err)
						}
					}

					return
				}
			}

			w.dispatcher.workerCh <- w
		}
	}()
}

func (w *Worker) stop() {
	w.quitCh <- struct{}{}
}
