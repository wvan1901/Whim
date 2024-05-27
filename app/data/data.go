package data

import (
	"fmt"
	"time"
	"wicho/whim/app/consts"
	"wicho/whim/app/row"
)

const (
    SPACES_IN_TAB = 4
)

func EditorInsertChar(appData *consts.EditorConfig, r rune){
    if appData.CursorPosY == appData.NumRows {
        row.EditorInsertRow(appData, appData.NumRows, "")
    }
    row.EditorRowInsertChar(appData.Row[appData.CursorPosY], appData.CursorPosX, r)
    //? Should this be moved to editorRowInsertChar?
    appData.Dirty++
    appData.CursorPosX += 1
}

//TODO: this should probably be moved elsewhere
func EditorSetStatusMessage(data *consts.EditorConfig, messages ...string){
    newStatusMsg := ""
    for _, msg := range messages {
        newStatusMsg += msg
    }
    data.StatusMessage = newStatusMsg
    data.StatusMessageTime = time.Now()
    
}

func EditorDelChar(appData *consts.EditorConfig){
    if appData.CursorPosY == appData.NumRows {
        return
    }
    if appData.CursorPosX == 0 && appData.CursorPosY == 0{
        return
    }
    curRow := appData.Row[appData.CursorPosY]
    if appData.CursorPosX > 0 {
        row.EditorRowDelChar(curRow, appData.CursorPosX - 1)
        appData.CursorPosX--
        appData.Dirty++
    } else {
        appData.CursorPosX = appData.Row[appData.CursorPosY-1].Size
        row.EditorRowAppendString(appData.Row[appData.CursorPosY-1], *curRow.Runes)
        row.EditorDelRow(appData, appData.CursorPosY)
        appData.CursorPosY--
    }
}

func EditorInsertNewLine(appData *consts.EditorConfig){
    if appData.CursorPosX == 0 {
        row.EditorInsertRow(appData, appData.CursorPosY, "")
    } else {
        curRow := appData.Row[appData.CursorPosY]
        leftSideString := (*curRow.Runes)[:appData.CursorPosX]
        rightSideString := (*curRow.Runes)[appData.CursorPosX:]
        curRow.Runes = &leftSideString
        fmt.Println("Add:", &leftSideString)
        row.EditorInsertRow(appData, appData.CursorPosY+1, rightSideString)
        curRow.Size = len(*curRow.Runes)
        row.EditorUpdateRow(curRow)
    }
    appData.CursorPosY++
    appData.CursorPosX = 0
}

