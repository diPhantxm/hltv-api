package taskScheduler

import "testing"

var result int = 0

func add(params ...interface{}) {
	num := params[0].([]interface{})[0].(int)
	result += num
}

func TestStraightScheduler(t *testing.T) {
	const N int = 3
	tests := []Task{
		NewTask(add, 3),
		NewTask(add, 5),
		NewTask(add, 28),
	}

	expectedResult := result
	result = 0

	for i := 0; i < N; i++ {
		go PushTests(tests)
	}

	for !GetStraightScheduler().IsEmpty() {

	}

	if expectedResult != result {
		t.Errorf("Got %d, Expected: %d\n", result, expectedResult)
	}
}

func PushTests(tasks []Task) {
	for _, task := range tasks {
		GetStraightScheduler().Add(task.Func, task.Params...)
	}
}
