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
			clearRow()
			fmt.Print(statusIcon(tc, frame) + " " + tc.title)
			frame++

		case <-tc.Context.Done():
			clearRow()
			fmt.Print(statusIcon(tc, 0)+" "+color.HiWhiteString(tc.title), "\n")

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
func statusIcon(tc *TaskCtx, frame int) string {
	switch tc.status {
	default:
		return ""
	case statusSuccess:
		if tc.ongoing {
			return color.YellowString(loadingCharacters[frame%len(loadingCharacters)])
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
