package taskScheduler

import (
	"sync"
	"time"
)

type straightScheduler struct {
	mutex    sync.Mutex
	Interval time.Duration
	tasks    []Task
}

func newStraightScheduler(interval time.Duration) *straightScheduler {
	ss := &straightScheduler{
		Interval: interval,
	}

	go ss.run()

	return ss
}

var straightInstance *straightScheduler = nil

func GetStraightScheduler() *straightScheduler {
	if straightInstance == nil {
		straightInstance = newStraightScheduler(time.Duration(400) * time.Millisecond)
	}

	return straightInstance
}

func (s *straightScheduler) Add(task func(...interface{}), params ...interface{}) {
	s.mutex.Lock()
	s.tasks = append(s.tasks, NewTask(task, params))
	s.mutex.Unlock()
}

func (s *straightScheduler) IsEmpty() bool {
	return s.Length() == 0
}

func (s *straightScheduler) Length() int {
	return len(s.tasks)
}

func (s *straightScheduler) run() {
	for {
		time.Sleep(s.Interval) // prevents ban for too frequent request to website

		if len(s.tasks) != 0 {
			s.tasks[0].Func(s.tasks[0].Params...)
			s.tasks = s.tasks[1:]
		}
	}
}
