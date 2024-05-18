package data

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
    SPACES_IN_TAB = 4
)

type EditorConfig struct {
    OldTerminalState term.State
    ScreenRows int
    ScreenColumns int
    RowOffSet int
    ColOffSet int
    CursorPosX int
    CursorPosY int
    RendorIndexX int
    NumRows int
    Row []*EditorRow
    FileName *string
    ABuf strings.Builder
    StatusMessage string
    StatusMessageTime time.Time
    Dirty int
}

type EditorRow struct {
    Size int
    Runes *string
    Render *string
    RenderSize int
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
    newRow := EditorRow{
        Size: len(aString),
        Runes: &aString,
        Render: nil,
        RenderSize: 0,
    }
    //Render string
    editorUpdateRow(&newRow)
    newSlice := append(editorData.Row, &newRow)
    editorData.Row = newSlice
    editorData.NumRows++
    editorData.Dirty++
}

func editorUpdateRow(row *EditorRow){
    row.Render = row.Runes
    row.RenderSize = len(*row.Render)
    if row.RenderSize < 1 {
        return
    }
    // Maybe making a new slice would be better?
    sliceOfRunes := []rune(*row.Render)[:0]
    //sliceOfRunes := make([]rune, 0)
    for _, aRune := range []rune(*row.Render) {
        if aRune == '\t' {
            sliceOfSpaces := make([]rune, 0)
            for j := 0; j<SPACES_IN_TAB; j++ {
                sliceOfSpaces = append(sliceOfSpaces, '-')
            }
            sliceOfRunes = append(sliceOfRunes, sliceOfSpaces...)
        } else {
            sliceOfRunes = append(sliceOfRunes, aRune)
        }
    }
    renderString := string(sliceOfRunes)
    row.Render = &renderString
    row.RenderSize = len(renderString)
}

// IDEA: Make the cursor go to the end of the tab
// IDEA: Make the visual part not go over 1 for x and y
// TODO: Fix Tabs when going to a diffrent line
func EditorRowCxToRx(row *EditorRow, cursorPosX int) int {
    renderX := 0
    for i, aRune := range []rune(*row.Runes) {
        if i >= cursorPosX{
            fmt.Println("BREAK| I:", i+1)
            break 
        }
        if aRune == '\t' {
            // renderX += (SPACES_IN_TAB-1) - (renderX % SPACES_IN_TAB)
            renderX += (SPACES_IN_TAB-1)
        }
        renderX++
    }
    return renderX
}

func editorRowInsertChar(row *EditorRow, at int, r rune){
    if (at < 0 || at > row.Size){
        at = row.Size
    }
    // TODO: Insert Rune in row
    newRunes := (*row.Runes)[:at] + string(r) + (*row.Runes)[at:]
    row.Runes = &newRunes

    row.Size += 1
    editorUpdateRow(row)
}

func EditorInsertChar(appData *EditorConfig, r rune){
    if appData.CursorPosY == appData.NumRows {
        appData.EditorAppendRow("")
    }
    editorRowInsertChar(appData.Row[appData.CursorPosY], appData.CursorPosX, r)
    //? Should this be moved to editorRowInsertChar?
    appData.Dirty++
    appData.CursorPosX += 1
}

//TODO: this should probably be moved elsewhere
func (data *EditorConfig) EditorSetStatusMessage(messages ...string){
    newStatusMsg := ""
    for _, msg := range messages {
        newStatusMsg += msg
    }
    data.StatusMessage = newStatusMsg
    data.StatusMessageTime = time.Now()
    
}
