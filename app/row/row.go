package row

import (
	"fmt"
	"wicho/whim/app/consts"
	"wicho/whim/app/highlight"
)

const (
    SPACES_IN_TAB = 4
)

// IDEA: Make the cursor go to the end of the tab
// IDEA: Make the visual part not go over 1 for x and y
// TODO: Fix Tabs when going to a diffrent line
func EditorRowCxToRx(row *consts.EditorRow, cursorPosX int) int {
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

// This is the oppisite of EditorRowCxToRx
func EditorRowRxToCx(curRow *consts.EditorRow, renderPosX int) int {
    curRx := 0
    for cursorPosX := 0; cursorPosX < curRow.Size; cursorPosX++{
        if (*curRow.Runes)[cursorPosX] == '\t'{
            curRx += (SPACES_IN_TAB - 1) - (curRx % SPACES_IN_TAB)
        }
        curRx++
        if curRx > renderPosX {
            return cursorPosX
        }
    }
    return curRow.Size-1
}

//editorUpdateRow
func EditorUpdateRow(row *consts.EditorRow){
    row.Render = row.Runes
    row.RenderSize = len(*row.Render)
    if row.RenderSize < 1 {
        return
    }

    sliceOfRunes := []rune(*row.Render)[:0]
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

    highlight.EditorUpdateSyntax(row)
}

//editorInsertRow

//editorFreeRow

//EditorDelRow
func EditorDelRow(appData *consts.EditorConfig, at int){
    if at < 0 || at >= appData.NumRows{
        return
    }
    newRows := append(appData.Row[:at], appData.Row[at+1:]...)
    appData.Row = newRows

    appData.NumRows--
    appData.Dirty++
}

//EditorRowInsertChar
func EditorRowInsertChar(row *consts.EditorRow, at int, r rune){
    if (at < 0 || at > row.Size){
        at = row.Size
    }
    newRunes := (*row.Runes)[:at] + string(r) + (*row.Runes)[at:]
    row.Runes = &newRunes

    row.Size += 1
    EditorUpdateRow(row)
}

//EditorRowAppendString
func EditorRowAppendString(row *consts.EditorRow, newString string){
    newRowString := *row.Runes + newString
    row.Runes = &newRowString
    //TODO: Why is there a + '\n' for the Size?
    row.Size = row.Size + len(newString) + '\n'
    EditorUpdateRow(row)
}

//EditorRowDelChar
func EditorRowDelChar(row *consts.EditorRow, at int){
    if at < 0 || at > row.Size {
        return
    }
    sliceRunes := []rune(*row.Runes)
    result := append(sliceRunes[0:at], sliceRunes[at+1:]...)
    newRunes := string(result)
    row.Runes = &newRunes
    row.Size = len(*row.Runes)
    EditorUpdateRow(row)
}

