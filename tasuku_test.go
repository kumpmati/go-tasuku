package tasuku

import (
	"errors"
	"testing"
	"time"
)

func TestTaskSuccess(test *testing.T) {
	num, err := Task("returns 5", func(t *TaskCtx) (int, error) { return 5, nil })

	if num != 5 || err != nil {
		test.Fail()
	}
}

func TestTaskSetTitle(test *testing.T) {
	num, err := Task("task 1", func(t *TaskCtx) (int, error) {
		t.SetTitle("success")

		if t.title != "success" || t.status != statusSuccess {
			test.Fail()
		}

		return 1, nil
	})

	if num != 1 || err != nil {
		test.Fail()
	}
}

func TestTaskSetWarning(test *testing.T) {
	num, err := Task("task 2", func(t *TaskCtx) (int, error) {
		t.SetWarning("warning")

		if t.title != "warning" || t.status != statusWarning {
			test.Fail()
		}

		return 1, nil
	})

	if num != 1 || err != nil {
		test.Fail()
	}
}

func TestTaskError(test *testing.T) {
	_, err := Task("task 3", func(t *TaskCtx) (int, error) { return 1, errors.New("some error") })

	if err == nil || err.Error() != "some error" {
		test.Fail()
	}
}

func TestTaskCancel(test *testing.T) {
	_, err := Task("task 4", func(t *TaskCtx) (int, error) {
		t.Cancel("some reason")

		return 1, nil
	})

	if err == nil {
		test.Error("received nil error")
	}
}

func TestTaskStuff(test *testing.T) {
	Task("asd", func(t *TaskCtx) (int, error) {
		<-time.After(time.Second)
		t.SetTitle("new title")
		<-time.After(time.Second)
		t.SetWarning("some warning")
		<-time.After(time.Second)
		t.SetError(errors.New("some error"))

		return 1, nil
	})

	<-time.After(time.Second)
}
