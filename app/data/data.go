package data

import (
	"golang.org/x/term"
    "strings"
    "fmt"
    "os"
)

type EditorConfig struct {
    OldTerminalState term.State
    ScreenRows int
    ScreenColumns int
    CursorPosX int
    CursorPosY int
    NumRows int
    Row EditorRow
    ABuf strings.Builder
}

type EditorRow struct {
    Size int
    Runes *string
}


func (appData *EditorConfig) Die(){
    fmt.Print("\033[2J")
    fmt.Print("\033[H")
    //Add os exit method! 
    term.Restore(0, &appData.OldTerminalState)
    defer os.Exit(1)
}
