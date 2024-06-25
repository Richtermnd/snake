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

func ClearAndWrite(b []byte) (int, error) {
	ClearScreen()
	return os.Stdout.Write(b)
}

func ClearScreen() {
	// exec.Command("clear")
	fmt.Print("\033[H\033[2J")
}
