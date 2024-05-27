package find

import (
	"strings"
	"wicho/whim/app/data"
	"wicho/whim/app/input"
)

func EditorFind(appData *data.EditorConfig){
    queryPrompt := input.EditorPrompt(appData, "(ESC to cancel) Search:")
    if queryPrompt == nil {
        return
    }

    for i := 0; i < appData.NumRows; i++ {
        curRow := appData.Row[i]
        strIndex := strings.Index(*curRow.Runes, *queryPrompt)
        if strIndex > -1 {
            appData.CursorPosY = i
            appData.CursorPosX = data.EditorRowRxToCx(curRow, strIndex)
            appData.RowOffSet = appData.NumRows
            break
        }
    }
}
