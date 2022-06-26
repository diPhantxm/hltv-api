package taskScheduler

type Task struct {
	Func   func(...interface{})
	Params []interface{}
}

func NewTask(function func(...interface{}), params ...interface{}) Task {
	return Task{
		Func:   function,
		Params: params,
	}
}

type TaskScheduler interface {
	Add(task func(...interface{}), params ...interface{})
	IsEmpty() bool
	Length() int
}
