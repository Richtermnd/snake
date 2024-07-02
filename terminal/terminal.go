package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func GetTerminalSize() (int, int) {
	w, h, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}
	return w, h
}

func ResetAndWrite(b []byte) (int, error) {
	ResetCursor()
	return os.Stdout.Write(b)
}

func ResetCursor() {
	// move cursor ton (0, 0) position
	fmt.Print("\033[H")
}
