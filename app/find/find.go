package find

import (
	"strings"
	"wicho/whim/app/consts"
	"wicho/whim/app/input"
	"wicho/whim/app/row"
)

func editorFindCallback(appData *consts.EditorConfig, query string, aRune rune){
    if appData.StringFindData == nil {
        appData.StringFindData = &consts.FindData{
            LastMatch: -1,
            Direction: 1,
        }
    }

    if aRune == '\r' || aRune == consts.ESC {
        return
    } else if aRune == consts.RIGHT_ARROW || aRune == consts.DOWN_ARROW {
        appData.StringFindData.Direction = 1
    } else if aRune == consts.LEFT_ARROW || aRune == consts.UP_ARROW {
        appData.StringFindData.Direction = -1
    } else {
        appData.StringFindData.LastMatch = -1
        appData.StringFindData.Direction = 1
    }

    if appData.StringFindData.LastMatch == -1 {
        appData.StringFindData.Direction = 1
    }
    current := appData.StringFindData.LastMatch
    for i := 0; i < appData.NumRows; i++ {
        current += appData.StringFindData.Direction
        if current == -1 {
            current = appData.NumRows - 1
        } else if current == appData.NumRows {
            current = 0 
        }
        curRow := appData.Row[current]
        strIndex := strings.Index(*curRow.Runes, query)
        if strIndex > -1 {
            appData.StringFindData.LastMatch = current
            appData.CursorPosY = current
            appData.CursorPosX = row.EditorRowRxToCx(curRow, strIndex)
            appData.RowOffSet = appData.NumRows
            break
        }
    }
}

func EditorFind(appData *consts.EditorConfig){
    savedCursorPosX := appData.CursorPosX
    savedCursorPosY := appData.CursorPosY
    savedColOffset := appData.ColOffSet
    savedRowOffset := appData.RowOffSet
    appData.StringFindData = nil

    queryPrompt := input.EditorPrompt(appData, "(ESC to cancel) Search:", editorFindCallback)
    if queryPrompt == nil {
        return
    } else {
        appData.CursorPosX = savedCursorPosX
        appData.CursorPosY = savedCursorPosY
        appData.ColOffSet = savedColOffset
        appData.RowOffSet = savedRowOffset
    }
    appData.StringFindData = nil
}
