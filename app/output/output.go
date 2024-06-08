package output

import (
	"fmt"
	"time"
	"unicode"
	"wicho/whim/app/consts"
	"wicho/whim/app/highlight"
	"wicho/whim/app/row"
	"wicho/whim/app/terminal"
)

func editorScroll(editorData *consts.EditorConfig) {
    editorData.RendorIndexX = 0
    if editorData.CursorPosY < editorData.NumRows {
        editorData.RendorIndexX = row.EditorRowCxToRx(editorData.Row[editorData.CursorPosY], editorData.CursorPosX)
    }

    if editorData.CursorPosY < editorData.RowOffSet {
        editorData.RowOffSet = editorData.CursorPosY
    }
    if editorData.CursorPosY >= editorData.RowOffSet+editorData.ScreenRows { 
        editorData.RowOffSet = editorData.CursorPosY - editorData.ScreenRows + 1
    }
    if editorData.RendorIndexX < editorData.ColOffSet{
        editorData.ColOffSet = editorData.RendorIndexX
    }
    if editorData.RendorIndexX >= editorData.ColOffSet + editorData.ScreenColumns {
        editorData.ColOffSet = editorData.RendorIndexX - editorData.ScreenColumns + 1
    }
}

func editorDrawRows(editorData *consts.EditorConfig){
    for y:=0; y<editorData.ScreenRows; y++ {
        fileRow := y + editorData.RowOffSet
        if fileRow >= editorData.NumRows{
            if editorData.NumRows == 0 && y == editorData.ScreenRows/3 {
                welcome := fmt.Sprintf("Whim Editor -- version %s", consts.WHIM_VERSION)
                if len(welcome) > editorData.ScreenColumns{
                    welcome = "v:"+ consts.WHIM_VERSION
                }
                padding := (editorData.ScreenColumns - len(welcome))/2
                if padding>0{
                    editorData.ABuf.WriteString("~")
                    padding--
                }
                for padding>0 {
                    padding--
                    editorData.ABuf.WriteString(" ")
                }
                editorData.ABuf.WriteString(welcome)
            } else {
                editorData.ABuf.WriteString("~")
            }
        } else {
            rowlength := editorData.Row[fileRow].RenderSize - editorData.ColOffSet
            if rowlength < 0 {
                rowlength = 0
            }
            renderString := ""
            renderHighlight := []int{}
            if rowlength > editorData.ScreenColumns {
                rowlength = editorData.ScreenColumns
                renderString = (*editorData.Row[fileRow].Render)[editorData.ColOffSet:rowlength]
                renderHighlight = editorData.Row[fileRow].Highlights[editorData.ColOffSet:rowlength]
            } else if rowlength == 0 {
                renderString = ""
                renderHighlight = []int{}
            } else {
                renderString = (*editorData.Row[fileRow].Render)[editorData.ColOffSet:]
                renderHighlight = editorData.Row[fileRow].Highlights[editorData.ColOffSet:]
            }
            renderRunes := []rune(renderString)
            currentColor := -1
            for j := 0; j < len(renderRunes); j++{
                //TODO: Test out opening a executable file if this if statement executes
                if (unicode.IsControl(renderRunes[j])){
                    symbol := "?"
                    if renderRunes[j] <= 26{
                        symbol = "@"
                    }
                    editorData.ABuf.WriteString("\033[7m")
                    editorData.ABuf.WriteString(symbol)
                    editorData.ABuf.WriteString("\033[m")
                    if currentColor != -1 {
                        unicodeColor := fmt.Sprintf("\033[%dm", currentColor)
                        editorData.ABuf.WriteString(unicodeColor)
                    }
                } else if renderHighlight[j] == consts.HL_NORMAL {
                    if currentColor != -1 {
                        editorData.ABuf.WriteString("\033[39m")
                        currentColor = -1
                    }
                    editorData.ABuf.WriteString(string(renderRunes[j]))
                } else {
                    colorInt := highlight.EditorSyntaxToColor(renderHighlight[j])
                    if colorInt != currentColor {
                        currentColor = colorInt
                        unicodeColor := fmt.Sprintf("\033[%dm", colorInt)
                        editorData.ABuf.WriteString(unicodeColor)
                    }
                    editorData.ABuf.WriteString(string(renderRunes[j]))
                }
            }
            editorData.ABuf.WriteString("\033[39m")
        }
        editorData.ABuf.WriteString("\033[K")
        editorData.ABuf.WriteString("\r\n")
    }
}

func editorDrawStatusBar(appData *consts.EditorConfig){
    appData.ABuf.WriteString("\033[7m")
    length := 0
    status := "[No Name]"
    rightSideStatus := ""
    rightSideLength := 0
    if appData.FileName != nil {
        status = *appData.FileName
        length = len(status)
        if length>20{
            length = 20
        }
        status = status[:length]+fmt.Sprintf(" - %d lines", appData.NumRows)
        if appData.Dirty > 0{
            status += "(Modified)"
        }
        length = len(status)
        editorSyntax := "no ft"
        if appData.EditorSyntax != nil {
            editorSyntax = appData.EditorSyntax.Filetype
        }
        rightSideStatus = fmt.Sprintf("%s | %d/%d", editorSyntax, appData.CursorPosY+1, appData.NumRows)
        rightSideLength = len(rightSideStatus)
    }
    appData.ABuf.WriteString(status)
    for length < appData.ScreenColumns {
        if (appData.ScreenColumns - length) == rightSideLength{
            appData.ABuf.WriteString(rightSideStatus)
            break
        } else {
            appData.ABuf.WriteString(" ")
            length++
        }
    }
    appData.ABuf.WriteString("\033[m")
    appData.ABuf.WriteString("\r\n")
}

func editorDrawMessageBar(appData *consts.EditorConfig){
    appData.ABuf.WriteString("\033[K")
    messageLength := len(appData.StatusMessage)
    if messageLength > appData.ScreenColumns{
        messageLength = appData.ScreenColumns
    }
    if messageLength>0 && (time.Since(appData.StatusMessageTime) < 5*time.Second){//1 sec = 1000000000 nano sec
        appData.ABuf.WriteString(appData.StatusMessage[:messageLength])
    }
}

func EditorRefreshScreen(appData *consts.EditorConfig){
    editorScroll(appData)
    appData.ABuf.WriteString("\033[?25l")
    appData.ABuf.WriteString("\033[H")
    editorDrawRows(appData)
    editorDrawStatusBar(appData)
    editorDrawMessageBar(appData)

    terminal.SetCursorPosition(appData)

    appData.ABuf.WriteString("\033[?25h")

    fmt.Print(appData.ABuf.String())
    appData.ABuf.Reset()
}

func EditorSetStatusMessage(data *consts.EditorConfig, messages ...string){
    newStatusMsg := ""
    for _, msg := range messages {
        newStatusMsg += msg
    }
    data.StatusMessage = newStatusMsg
    data.StatusMessageTime = time.Now()
    
}


