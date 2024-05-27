package highlight

import (
    "unicode"
	"wicho/whim/app/consts"
)

func EditorUpdateSyntax(curRow *consts.EditorRow){
    rowHl := make([]int, curRow.RenderSize)
    renderRunes := []rune(*curRow.Render)
    for i := 0; i < curRow.RenderSize; i++ {
        if unicode.IsDigit(renderRunes[i]){
            rowHl[i] = consts.HL_NUMBER
        }
    }
    curRow.Highlights = rowHl
}

func EditorSyntaxToColor(hl int) int {
    switch hl {
    case consts.HL_NUMBER:
        return 31
    default:
        return 37
    }
}
