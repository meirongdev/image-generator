package pool

import (
	"log"
	"sync"

	"github.com/meirongdev/image-generator/internal/drawing"
)

type Task struct {
	ID        string
	ImageName string
}

type TaskPool struct {
	ch            chan *Task
	sm            sync.Map
	maxGoroutines int
}

func NewTaskPool(maxGoroutines, waitingTasks int) *TaskPool {
	return &TaskPool{
		ch:            make(chan *Task, waitingTasks),
		maxGoroutines: maxGoroutines,
	}
}

func (tp *TaskPool) SubmitTask(task *Task) {
	tp.ch <- task
}

func (tp *TaskPool) Start() {
	for i := range tp.maxGoroutines {
		go func(i int) {
			log.Printf("Worker %d is running", i)
			for t := range tp.ch {
				path := drawing.DrawOne(t.ImageName)
				tp.sm.Store(t.ID, path)
				log.Printf("Stored: %v\n", t)
			}
		}(i)
	}
}

func (tp *TaskPool) Stop() {
	log.Println("closing the task pool")
	// close the channel
	close(tp.ch)
}

func (tp *TaskPool) GetResult(id string) (string, bool) {
	path, ok := tp.sm.Load(id)
	if !ok {
		return "", false
	}
	return path.(string), true
}
