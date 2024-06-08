package data

import (
	"fmt"
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
    row.EditorRowInsertChar(appData, appData.Row[appData.CursorPosY], appData.CursorPosX, r)
    appData.Dirty++
    appData.CursorPosX += 1
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
        row.EditorRowDelChar(appData, curRow, appData.CursorPosX - 1)
        appData.CursorPosX--
        appData.Dirty++
    } else {
        appData.CursorPosX = appData.Row[appData.CursorPosY-1].Size
        row.EditorRowAppendString(appData, appData.Row[appData.CursorPosY-1], *curRow.Runes)
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
        row.EditorUpdateRow(appData, curRow)
    }
    appData.CursorPosY++
    appData.CursorPosX = 0
}

