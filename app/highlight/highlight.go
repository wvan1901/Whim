package highlight

import (
	"strings"
	"unicode"
	"wicho/whim/app/consts"
)

func EditorUpdateSyntax(curRow *consts.EditorRow){
    rowHl := make([]int, curRow.RenderSize)
    renderRunes := []rune(*curRow.Render)

    prevSeperator := true

    i := 0
    for i < curRow.RenderSize {
        curRune := renderRunes[i]
        prevHl := consts.HL_NORMAL
        if i > 0 {
            prevHl = rowHl[i-1]
        } else {
            prevHl = consts.HL_NORMAL
        }

        if unicode.IsDigit(curRune) && (prevSeperator || prevHl == consts.HL_NUMBER) || (curRune == '.' && prevHl == consts.HL_NUMBER) {
            rowHl[i] = consts.HL_NUMBER
            i++
            prevSeperator = false
            continue
        }

        prevSeperator = isSeperator(curRune)
        i++
    }
    curRow.Highlights = rowHl
}

func isSeperator(aRune rune) bool {
    isRuneSpace := unicode.IsSpace(aRune)
    isEmpty := aRune == '\u0000'
    isInvalid := strings.ContainsRune(",.()+-/*=~%<>[];", aRune)
    return isRuneSpace || isEmpty || isInvalid
}

func EditorSyntaxToColor(hl int) int {
    switch hl {
    case consts.HL_NUMBER:
        return 31
    case consts.HL_MATCH:
        return 34
    default:
        return 37
    }
}
