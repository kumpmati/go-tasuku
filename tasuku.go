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

type TaskCtx struct {
	Context context.Context
	cancel  context.CancelFunc

	ongoing bool
	status  string // "success" | "warning" | "error" | "cancelled"
	title   string
	detail  string
	err     error
}

// Task runs the given function, and prints the status of the task
// to the terminal.
func Task[R any](title string, fn func(t *TaskCtx) (R, error)) (R, error) {
	ctx, cancel := context.WithCancel(context.Background())

	tc := TaskCtx{
		Context: ctx,
		cancel:  cancel,
		title:   title,
		status:  statusSuccess,
		ongoing: true,
	}

	wg := sync.WaitGroup{}

	go tc.print(&wg)

	result, err := fn(&tc)

	tc.ongoing = false

	if err != nil {
		tc.SetError(err)
	}

	cancel()  // cancel the context, this tells the `task.print` loop to stop.
	wg.Wait() // wait for print to finish

	return result, tc.err
}

// SetTitle updates the title of the task.
func (tc *TaskCtx) SetTitle(text string) {
	tc.title = text
}

// SetWarning sets the task state to "warning", updates the task detail text.
func (tc *TaskCtx) SetWarning(text string) {
	tc.detail = text
	tc.status = statusWarning
}

// SetError changes the task title and sets the status to "error".
func (tc *TaskCtx) SetError(err error) {
	tc.err = err
	tc.detail = tc.err.Error()
	tc.status = statusError
}

// ClearError clears the current error and detail, and changes the state to "success"
func (tc *TaskCtx) ClearError() {
	tc.err = nil
	tc.status = statusSuccess
	tc.detail = ""
}

// SetDetail changes the detail text that is shown after the task completes.
func (tc *TaskCtx) SetDetail(text string) {
	tc.detail = text
}

// Cancel cancels the task context and sets the task state to "cancelled".
// If `reason` is non-empty, it's set as the task detail.
func (tc *TaskCtx) Cancel(reason string) {
	tc.err = context.Canceled
	tc.status = statusCancelled
	if reason != "" {
		tc.detail = reason
	}

	tc.ongoing = false
	tc.cancel()
}
