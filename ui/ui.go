package ui

import (
	"log"
	"strconv"

	"github.com/johnmastroberti/gochat/msg"
	gc "github.com/rthornton128/goncurses"
)

// Initialize ncurses and set global settings
func cursesInit() *gc.Window {
	// Initialize ncurses
	stdscr, err := gc.Init()
	if err != nil {
		log.Panic(err)
	}

	// Turn off character echo, hide the cursor and disable input buffering
	gc.Echo(false)
	gc.Raw(true)
	gc.Cursor(0)

	return stdscr
}

// Construct the display and input windows
func makeWindows(stdscr *gc.Window) (*gc.Window, *gc.Window) {
	// Compute window coordinates and geometry
	rows, cols := stdscr.MaxYX()

	inputHeight, inputWidth := 3, cols
	dispHeight, dispWidth := rows-inputHeight, cols

	iy, ix := dispHeight, 0
	dy, dx := 0, 0

	// Create the windows
	dwin, err := gc.NewWindow(dispHeight, dispWidth, dy, dx)
	if err != nil {
		log.Panic(err)
	}
	iwin, err := gc.NewWindow(inputHeight, inputWidth, iy, ix)
	if err != nil {
		log.Panic(err)
	}
	dwin.Box(0, 0)
	iwin.Box(0, 0)
	dwin.Refresh()
	iwin.Refresh()
	return dwin, iwin
}

func cleanupCurses(dwin *gc.Window, iwin *gc.Window) {
	iwin.Delete()
	dwin.Delete()
	gc.End()
}

// RunUILoop runs the main loop of the interface, where input
// is accepted from the user and information is displayed
// to the screen
func RunUILoop(uiEvents chan msg.Message, userInput chan msg.Message) {
	stdscr := cursesInit()
	dispWin, inputWin := makeWindows(stdscr)

	defer cleanupCurses(dispWin, inputWin)

	userKeys := make(chan gc.Key)
	ready := make(chan bool)

	go getUserInput(userKeys, ready, inputWin)
	defer close(ready) // will cause getUserInput to return

	for {
		select {
		case e := <-uiEvents:
			if handleEvent(e, dispWin) {
				return
			}

		case c := <-userKeys:
			handleInput(c, inputWin, userInput)

		// indicate that we are ready for user input
		// if there is no event or input to process
		case ready <- true:
		}
	}

}

// Returns true if ui should exit
func handleEvent(e msg.Message, dispWin *gc.Window) bool {
	// For now, just print to the display window
	dispWin.Erase()
	dispWin.Box(0, 0)
	dispWin.MovePrint(1, 1, "Type: ", e.Type)
	dispWin.MovePrint(2, 1, "From: ", e.From)
	dispWin.MovePrint(3, 1, "To: ", e.To)
	dispWin.MovePrint(4, 1, "Content: ", e.Content)
	dispWin.MovePrint(5, 1, "Username: ", e.Username)
	dispWin.MovePrint(6, 1, "Email: ", e.Email)
	dispWin.MovePrint(7, 1, "Password: ", e.Password)

	dispWin.Refresh()

	return e.Content == "/exit"
}

var inputContents []byte

func handleInput(c gc.Key, inputWin *gc.Window, userInput chan msg.Message) {
	updated := false
	switch {
	case c == gc.KEY_BACKSPACE || c == 127:
		if len(inputContents) > 0 {
			inputContents = inputContents[:len(inputContents)-1]
			updated = true
		}

	case c == gc.KEY_ENTER || c == gc.KEY_RETURN:
		userInput <- msg.Message{Type: "USERINPUT", Content: string(inputContents)}
		inputContents = inputContents[:0]
		updated = true

	case isCharacter(c):
		inputContents = append(inputContents, byte(c))
		updated = true

	default:
		userInput <- msg.Message{Type: "USERINPUT", Content: "Unhandled character: " + gc.KeyString(c) + " (" + strconv.Itoa(int(c)) + ")"}
	}

	if updated {
		// Redraw input window
		inputWin.Erase()
		inputWin.Box(0, 0)
		// Only diplay as much of the input as the window can hold
		_, maxLen := inputWin.MaxYX()
		maxLen -= 2
		if maxLen > len(inputContents) {
			maxLen = len(inputContents)
		}
		inputWin.MovePrintf(1, 1, string(inputContents[len(inputContents)-maxLen:]))
		if maxLen < len(inputContents) {
			inputWin.MovePrintf(1, 1, "...")
		}
		inputWin.Refresh()
	}
}

func isCharacter(c gc.Key) bool {
	return c >= 32 && c <= 126
}

func getUserInput(input chan gc.Key, ready chan bool, inputWin *gc.Window) {
	for range ready {
		input <- inputWin.GetChar()
	}
}
