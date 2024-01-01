package tasuku

import (
	"context"
	"sync"
)

const (
	statusSuccess   = "success"
	statusWarning   = "warning"
	statusError     = "error"
	statusCancelled = "cancelled"
)

type taskState struct {
	status  string // "success" | "warning" | "error" | "cancelled"
	title   string
	detail  string
	err     error
	ongoing bool
}

type TaskCtx struct {
	Context context.Context
	cancel  context.CancelFunc

	s            chan taskState
	currentState taskState // current state
}

// Task runs the given function, and prints the status of the task
// to the terminal.
func Task[R any](title string, fn func(t *TaskCtx) (R, error)) (R, error) {
	ctx, cancel := context.WithCancel(context.Background())

	tc := TaskCtx{
		Context:      ctx,
		cancel:       cancel,
		s:            make(chan taskState, 1),
		currentState: taskState{title: title, status: statusSuccess},
	}

	wg := sync.WaitGroup{}

	tc.set(func(s *taskState) {}) // provide initial state for print
	go tc.print(tc.s, &wg)

	result, err := fn(&tc)

	if err != nil {
		tc.SetError(err)
	}

	cancel()  // cancel the context, this tells the `task.print` loop to stop.
	wg.Wait() // wait for print to finish

	return result, tc.currentState.err
}

func (tc *TaskCtx) set(fn func(s *taskState)) {
	fn(&tc.currentState)
	tc.s <- tc.currentState
}

// SetTitle updates the title of the task.
func (tc *TaskCtx) SetTitle(text string) {
	tc.set(func(prev *taskState) {
		prev.title = text
	})
}

// SetWarning sets the task state to "warning", updates the task detail text.
func (tc *TaskCtx) SetWarning(text string) {
	tc.set(func(s *taskState) {
		s.detail = text
		s.status = statusWarning
	})
}

// SetError changes the task title and sets the status to "error".
func (tc *TaskCtx) SetError(err error) {
	tc.set(func(s *taskState) {
		s.err = err
		s.detail = s.err.Error()
		s.status = statusError
	})
}

// ClearError clears the current error and detail, and changes the state to "success"
func (tc *TaskCtx) ClearError() {
	tc.set(func(s *taskState) {
		s.err = nil
		s.status = statusSuccess
		s.detail = ""
	})
}

// SetDetail changes the detail text that is shown after the task completes.
func (tc *TaskCtx) SetDetail(text string) {
	tc.set(func(s *taskState) {
		s.detail = text
	})
}

// Cancel cancels the task context and sets the task state to "cancelled".
// If `reason` is non-empty, it's set as the task detail.
func (tc *TaskCtx) Cancel(reason string) {
	tc.set(func(s *taskState) {
		s.err = context.Canceled
		s.status = statusCancelled
		if reason != "" {
			s.detail = reason
		}
	})

	tc.cancel()
}
