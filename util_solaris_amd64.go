package pty

import (
	"os"

	"golang.org/x/sys/unix"
)

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func Getsize(t *os.File) (rows, cols int, err error) {
	ws, err := unix.IoctlGetWinsize(int(t.Fd()), unix.TIOCGWINSZ)
	return int(ws.Row), int(ws.Col), err
}
