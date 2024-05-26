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
func (editorData *EditorConfig) EditorInsertRow(at int, aString string){
    if at < 0 || at > editorData.NumRows {
        return
    }

    newRow := EditorRow{
        Size: len(aString),
        Runes: &aString,
        Render: nil,
        RenderSize: 0,
    }

    //NOTE: This Was Erroring Due To Trying At Add Row At Index
    var firstHalfSlice []*EditorRow
    firstHalfSlice = append(firstHalfSlice, editorData.Row[:at]...)
    var secondHalfSlice []*EditorRow
    secondHalfSlice = append(secondHalfSlice, editorData.Row[at:]...)
    newSlice := append(firstHalfSlice, &newRow)
    newSlice = append(newSlice, secondHalfSlice...)
    editorData.Row = newSlice
    editorUpdateRow(&newRow)
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
    newRunes := (*row.Runes)[:at] + string(r) + (*row.Runes)[at:]
    row.Runes = &newRunes

    row.Size += 1
    editorUpdateRow(row)
}

func EditorInsertChar(appData *EditorConfig, r rune){
    if appData.CursorPosY == appData.NumRows {
        appData.EditorInsertRow(appData.NumRows, "")
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

func editorRowDelChar(row *EditorRow, at int){
    if at < 0 || at > row.Size {
        return
    }
    sliceRunes := []rune(*row.Runes)
    result := append(sliceRunes[0:at], sliceRunes[at+1:]...)
    newRunes := string(result)
    row.Runes = &newRunes
    row.Size = len(*row.Runes)
    editorUpdateRow(row)
}

func EditorDelChar(appData *EditorConfig){
    if appData.CursorPosY == appData.NumRows {
        return
    }
    if appData.CursorPosX == 0 && appData.CursorPosY == 0{
        return
    }
    curRow := appData.Row[appData.CursorPosY]
    if appData.CursorPosX > 0 {
        editorRowDelChar(curRow, appData.CursorPosX - 1)
        appData.CursorPosX--
        appData.Dirty++
    } else {
        appData.CursorPosX = appData.Row[appData.CursorPosY-1].Size
        editorRowAppendString(appData.Row[appData.CursorPosY-1], *curRow.Runes)
        editorDelRow(appData, appData.CursorPosY)
        appData.CursorPosY--
    }
}

func editorDelRow(appData *EditorConfig, at int){
    if at < 0 || at >= appData.NumRows{
        return
    }
    newRows := append(appData.Row[:at], appData.Row[at+1:]...)
    appData.Row = newRows

    appData.NumRows--
    appData.Dirty++
}

func editorRowAppendString(row *EditorRow, newString string){
    newRowString := *row.Runes + newString
    row.Runes = &newRowString
    //TODO: Why is there a + '\n' for the Size?
    row.Size = row.Size + len(newString) + '\n'
    editorUpdateRow(row)
}

func EditorInsertNewLine(appData *EditorConfig){
    if appData.CursorPosX == 0 {
        appData.EditorInsertRow(appData.CursorPosY, "")
    } else {
        curRow := appData.Row[appData.CursorPosY]
        leftSideString := (*curRow.Runes)[:appData.CursorPosX]
        rightSideString := (*curRow.Runes)[appData.CursorPosX:]
        curRow.Runes = &leftSideString
        fmt.Println("Add:", &leftSideString)
        appData.EditorInsertRow(appData.CursorPosY+1, rightSideString)
        curRow.Size = len(*curRow.Runes)
        editorUpdateRow(curRow)
    }
    appData.CursorPosY++
    appData.CursorPosX = 0
}
