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
    RowOffSet int
    CursorPosX int
    CursorPosY int
    NumRows int
    Row []*EditorRow
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

/*
    Row Operations
*/
func (editorData *EditorConfig) EditorAppendRow(aString string){
    // editorData.Row.Size = len(aString)
    // editorData.Row.Runes = &aString
    newRow := EditorRow{
        Size: len(aString),
        Runes: &aString,
    }
    newSlice := append(editorData.Row, &newRow)
    // fmt.Println("Slices", *newSlice[0])
    // fmt.Println("NewRow", newRow)
    //editorData.Die()
    editorData.Row = newSlice
    editorData.NumRows++
}
