package tasuku

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

const framerate = 10

var (
	rowClearChar  = "\033[H\033[2J" // TODO: make sure this works on all platforms
	successChar   = "✔"
	warningChar   = "⚠"
	errorChar     = "✖"
	cancelledChar = "-"
	arrowChar     = "→"

	loadingCharacters = []string{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}
)

// print starts an infinite loop that prints the current status
// of the task to the terminal. The loop is stopped when the task context
// is cancelled, printing the final status of the task before returning.
func (tc *TaskCtx) print(wg *sync.WaitGroup) {
	wg.Add(1)

	frame := 0

	for {
		select {

		case <-time.After(time.Second / framerate):
			fmt.Print(rowClearChar, statusIcon(tc.status, frame)+" "+tc.title)
			frame++

		case <-tc.Context.Done():
			fmt.Print(rowClearChar, statusIcon(tc.status, 0)+" "+color.HiWhiteString(tc.title), "\n")

			if tc.detail != "" {
				fmt.Print("  ", arrowChar+" "+tc.detail, "\n")
			}

			wg.Done()
			return
		}
	}
}

// statusString returns the string to be printed based on
// the current status of the task.
func statusIcon(status string, frame int) string {
	switch status {
	default:
		return color.YellowString(loadingCharacters[frame%len(loadingCharacters)])
	case statusSuccess:
		return color.GreenString(successChar)
	case statusErr:
		return color.RedString(errorChar)
	case statusWarning:
		return color.YellowString(warningChar)
	case statusCancelled:
		return color.BlueString(cancelledChar)
	}
}
