package main

import (
	"errors"
	"time"

	"github.com/kumpmati/go-tasuku"
)

func main() {
	tasuku.Task("basic", func(t *tasuku.TaskCtx) (int, error) {
		<-time.After(time.Second)
		return 1, nil
	})

	tasuku.Task("change title", func(t *tasuku.TaskCtx) (int, error) {
		<-time.After(time.Second)
		t.SetTitle("new title")

		<-time.After(time.Second)
		return 1, nil
	})

	tasuku.Task("set warning", func(t *tasuku.TaskCtx) (int, error) {
		<-time.After(time.Second)
		t.SetWarning("warning reason")

		<-time.After(time.Second)
		return 1, nil
	})

	tasuku.Task("set error", func(t *tasuku.TaskCtx) (int, error) {
		<-time.After(time.Second)
		t.SetError(errors.New("error reason"))

		<-time.After(time.Second)
		return 1, nil
	})

	tasuku.Task("clear error", func(t *tasuku.TaskCtx) (int, error) {
		<-time.After(time.Second)
		t.SetError(errors.New("error reason"))

		<-time.After(time.Second)
		t.ClearError()

		return 1, nil
	})

	tasuku.Task("cancel", func(t *tasuku.TaskCtx) (int, error) {
		<-time.After(time.Second)
		t.Cancel("some reason")

		<-time.After(time.Second)

		return 1, nil
	})
}
