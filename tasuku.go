package tasuku

import (
	"context"
	"sync"
)

const (
	statusSuccess   = "success"
	statusWarning   = "warning"
	statusErr       = "error"
	statusCancelled = "cancelled"
)

type TaskCtx struct {
	Context context.Context
	cancel  context.CancelFunc

	status string // "" | "success" | "warning" | "error" | "cancelled"
	title  string
	detail string
	err    error
}

// Task runs the given function, and prints the status of the task
// to the terminal.
func Task[R any](title string, fn func(t *TaskCtx) (R, error)) (R, error) {
	ctx, cancel := context.WithCancel(context.Background())

	task := TaskCtx{
		Context: ctx,
		cancel:  cancel,
		title:   title,
		status:  "",
	}

	wg := sync.WaitGroup{}

	go task.print(&wg)

	result, err := fn(&task)

	if err != nil {
		task.SetError(err)
	}

	cancel()  // cancel the context, this tells the `task.print` loop to stop.
	wg.Wait() // wait for print to finish

	return result, task.err
}

// SetTitle clears the Error and updates the Title of the task.
// It does not cancel the task context. Any previous errors will be cleared.
func (tc *TaskCtx) SetTitle(text string) {
	tc.err = nil
	tc.title = text
	tc.status = ""
}

// SetWarning changes the task title and sets the status to "warning",
// but does not cancel the task context. Any previous errors will be cleared.
func (tc *TaskCtx) SetWarning(text string) {
	tc.err = nil
	tc.detail = text
	tc.status = statusWarning
}

// SetError changes the task title and sets the status to "error",
// but does not cancel the task context.
func (tc *TaskCtx) SetError(err error) {
	tc.err = err
	tc.detail = tc.err.Error()
	tc.status = statusErr
}

// SetDetail changes the detail text that is shown after the task is done.
func (tc *TaskCtx) SetDetail(text string) {
	tc.detail = text
}

// Cancel changes the task title, sets the status to "cancelled"
// and cancels the task context.
func (tc *TaskCtx) Cancel(reason string) {
	tc.err = context.Canceled
	tc.status = statusCancelled
	tc.title = reason
	tc.cancel()
}
