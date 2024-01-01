# go-tasuku

[![Tests](https://github.com/kumpmati/go-tasuku/actions/workflows/go.yml/badge.svg)](https://github.com/kumpmati/go-tasuku/actions/workflows/go.yml)

The minimal task visualizer for go.

This is a go implementation of the great [privatenumber/tasuku](https://github.com/privatenumber/tasuku), a task visualizer library for Node.js.

## Install

```
go get -u github.com/kumpmati/go-tasuku
```

## Features

🚧 This library is still very much in progress, so some features from the JavaScript counterpart haven't been implemented yet:

- [x] Single task visualization
- [ ] TODO: parallel, nested & grouped tasks
- [ ] TODO: clearing completed tasks

## Usage

Task states

- `⣷` Loading - task has not completed yet
- `⚠` Warning - `SetWarning` was called
- `✖` Error - task returned an error, or `SetError` was called
- `-` Cancelled - task was cancelled manually
- `✔` Success - task completed without any errors

### Basic

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) { return 1, nil })
```

### SetTitle

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  <-time.After(time.Second * 5)
  t.SetTitle("this is taking longer than expected...")

  <-time.After(time.Second)
  t.SetTitle("done!")

  return 2, nil
})

// Terminal output
// ⣷ my task

// After 5 seconds
// ⣷ this is taking longer than expected...

// Completed
// ✔ done!
```

### SetDetail

Adds an extra message below the title after the task has completed.

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  t.SetDetail("some details")

  return 3, nil
})

// ✔ my task
//   → some details
```

### SetWarning

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  t.SetWarning("some warning")

  return 4, nil
})

// ⚠ my task
//   → some warning
```

### SetError / Returning errors

To show an error, either return the error at the end, or call `t.SetError` inside the task.

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  return 5, errors.new("some error")
})

// ✖ my task
//   → some error
```

or

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  // caller will receive this error, since the task function returns a nil error
  t.SetError(errors.New("some error"))

  return 6, nil
})

// ✖ my task
//  → some error
```

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  // this will set the task status to error during execution, but the
  // caller will only receive the "second error" in the return statement
  t.SetError(errors.New("first error"))

  return 7, errors.New("second error")
})

// ✖ my task
//  → second error
```

To clear any error set by `t.SetError`, you can call `t.ClearError` before returning from the task.

### Cancel

```go
result, err := tasuku.Task("my task", func(t *tasuku.TaskCtx) (int, error) {
  if condition {
    t.Cancel("cancellation reason")
    return 1, nil // return nil error so that the cancellation error is returned to caller
  }

  return 6, errors.New("custom error")
})

// condition == true:
// - my task
//  → cancellation reason

// condition == false:
// ✖ my task
//  → custom error
```
