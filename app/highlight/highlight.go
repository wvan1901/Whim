package highlight

import (
	"slices"
	"strings"
	"unicode"
	"wicho/whim/app/consts"
)

func EditorUpdateSyntax(appSyntax *consts.EditorSyntax, curRow *consts.EditorRow){
    rowHl := make([]int, curRow.RenderSize)
    renderRunes := []rune(*curRow.Render)

    if appSyntax == nil{
        curRow.Highlights = rowHl
        return
    }
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

        if slices.Contains(appSyntax.Flags, "HL_HIGHLIGHT_NUMBERS"){
            if unicode.IsDigit(curRune) && (prevSeperator || prevHl == consts.HL_NUMBER) || (curRune == '.' && prevHl == consts.HL_NUMBER) {
                rowHl[i] = consts.HL_NUMBER
                i++
                prevSeperator = false
                continue
            }
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

func EditorSelectSyntaxHighlight(appData *consts.EditorConfig){
    appData.EditorSyntax = nil
    if appData.FileName == nil {
        return
    }
    fileNameSplit := strings.Split(*appData.FileName, ".")
    if len(fileNameSplit) < 2 {
        return
    }
    extension := fileNameSplit[len(fileNameSplit)-1]

    HLDB := consts.HLDB()
    for j := 0; j < len(HLDB); j++{
        curEditorSyntax := HLDB[j]
        for _, item := range curEditorSyntax.Filematch {
            exten := item[1:]            
            if exten == extension {
                appData.EditorSyntax = &curEditorSyntax

                for filerow := 0; filerow < appData.NumRows; filerow++{
                    EditorUpdateSyntax(appData.EditorSyntax, appData.Row[filerow])
                }

                return
            }
        }
    }
}
