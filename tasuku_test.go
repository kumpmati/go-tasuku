package tasuku

import (
	"context"
	"errors"
	"testing"
)

func TestTaskSuccess(test *testing.T) {
	num, err := Task("task success", func(t *TaskCtx) (int, error) { return 5, nil })

	if num != 5 || err != nil {
		test.Fail()
	}
}

func TestTaskSetTitle(test *testing.T) {
	num, err := Task("task set title", func(t *TaskCtx) (int, error) {
		t.SetTitle("success")

		if t.currentState.title != "success" {
			test.Fail()
		}

		return 1, nil
	})

	if num != 1 || err != nil {
		test.Fail()
	}
}

func TestTaskSetDetail(test *testing.T) {
	num, err := Task("task set detail", func(t *TaskCtx) (int, error) {
		t.SetDetail("detail")

		if t.currentState.detail != "detail" {
			test.Fail()
		}

		return 1, nil
	})

	if num != 1 || err != nil {
		test.Fail()
	}
}

func TestTaskSetWarning(test *testing.T) {
	num, err := Task("task set warning", func(t *TaskCtx) (int, error) {
		t.SetWarning("warning")

		if t.currentState.detail != "warning" {
			test.Fail()
		}

		return 1, nil
	})

	if num != 1 || err != nil {
		test.Fail()
	}
}

func TestTaskError(test *testing.T) {
	_, err := Task("task set error", func(t *TaskCtx) (int, error) { return 1, errors.New("some error") })

	if err == nil || err.Error() != "some error" {
		test.Fail()
	}
}

func TestTaskCancel(test *testing.T) {
	_, err := Task("task cancel", func(t *TaskCtx) (int, error) {
		t.Cancel("some reason")

		return 1, nil
	})

	if !errors.Is(err, context.Canceled) {
		test.Error("expected context canceled error")
	}

	customErr := errors.New("my error")

	_, err = Task("task cancel with err return", func(t *TaskCtx) (int, error) {
		t.Cancel("some reason")

		return 1, customErr
	})

	if !errors.Is(err, customErr) {
		test.Error("expected custom error")
	}
}
