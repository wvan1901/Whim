package data

import (
	"golang.org/x/term"
    "strings"
)

type EditorConfig struct {
    OldTerminalState term.State
    ScreenRows int
    ScreenColumns int
    CursorPosX int
    CursorPosY int
    ABuf strings.Builder
}
