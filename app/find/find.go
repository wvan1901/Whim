package find

import (
	"strings"
	"wicho/whim/app/consts"
	"wicho/whim/app/data"
	"wicho/whim/app/input"
)

func editorFindCallback(appData *data.EditorConfig, query string, aRune rune){
    if aRune == '\r' || aRune == consts.ESC {
        return
    }

    for i := 0; i < appData.NumRows; i++ {
        curRow := appData.Row[i]
        strIndex := strings.Index(*curRow.Runes, query)
        if strIndex > -1 {
            appData.CursorPosY = i
            appData.CursorPosX = data.EditorRowRxToCx(curRow, strIndex)
            appData.RowOffSet = appData.NumRows
            break
        }
    }
}

func EditorFind(appData *data.EditorConfig){
    queryPrompt := input.EditorPrompt(appData, "(ESC to cancel) Search:", editorFindCallback)
    if queryPrompt == nil {
        return
    }
}
