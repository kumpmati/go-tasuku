package tasuku

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

const framerate = 10

var (
	successChar   = "✔"
	warningChar   = "⚠"
	errorChar     = "✖"
	cancelledChar = "-"
	arrowChar     = "→"

	loadingChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

// print starts an infinite loop that prints the current status
// of the task to the terminal. The loop is stopped when the task context
// is cancelled, printing the final status of the task before returning.
func (tc *TaskCtx) print(s <-chan taskState, wg *sync.WaitGroup) {
	wg.Add(1)

	var frame int
	var state taskState

	for {
		select {
		case newState := <-s:
			state = newState

		case <-time.After(time.Second / framerate):
			clearRow()
			fmt.Print(statusIcon(state.status, true, frame) + " " + state.title)
			frame++

		case <-tc.Context.Done():
			clearRow()
			fmt.Print(statusIcon(state.status, false, 0)+" "+color.HiWhiteString(state.title), "\n")

			if state.detail != "" {
				fmt.Print("  ", arrowChar+" "+state.detail, "\n")
			}

			wg.Done()
			return
		}
	}
}

// statusString returns the string to be printed based on
// the current status of the task.
func statusIcon(status string, ongoing bool, frame int) string {
	switch status {
	default:
		return ""
	case statusSuccess:
		if ongoing {
			return color.YellowString(loadingChars[frame%len(loadingChars)])
		}
		return color.GreenString(successChar)
	case statusError:
		return color.RedString(errorChar)
	case statusWarning:
		return color.YellowString(warningChar)
	case statusCancelled:
		return color.BlueString(cancelledChar)
	}
}

func clearRow() {
	w, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return
	}

	fmt.Print("\r\b", strings.Repeat(" ", w), "\r\b")
}
