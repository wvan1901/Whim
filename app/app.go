package app

import (
	"os"
	"strings"
	"time"
	"wicho/whim/app/consts"
	"wicho/whim/app/data"
	"wicho/whim/app/fileio"
	"wicho/whim/app/output"
	"wicho/whim/app/terminal"

	"fmt"

	"golang.org/x/term"
)

func RunApp(){
    oldState := terminal.EnableRawMode()
    defer terminal.DisableRawMode(&oldState)
    AppData := initEditor(&oldState)
    argsWithProg := os.Args
    if len(argsWithProg) >= 2 {
        fileio.EditorOpen(&AppData, argsWithProg[1])
    }
    AppData.EditorSetStatusMessage("HELP: Ctrl-s = save | Ctrl-C, q = Quit")
    for {
        output.EditorRefreshScreen(&AppData)
        editorProcessKeyPress(&AppData)
    }
}

func editorProcessKeyPress(appData *data.EditorConfig){
    keyReadRune := terminal.EditorReadKey()
    switch keyReadRune {
    case '\r':
        //TODO: Implement Enter
        data.EditorInsertNewLine(appData)
        break
    case consts.BACKSPACE, consts.DEL_KEY:
        data.EditorDelChar(appData)
        break
    case consts.ESC:
        break
    case consts.CONTROLCASCII:
        //TODO: Add a warning message if file is dirty and ask User to Confirm
        fmt.Println("<C>")
        terminal.DisableRawMode(&appData.OldTerminalState)
        fmt.Print("\033[2J")
        fmt.Print("\033[H")
        defer os.Exit(0)
        break
    case consts.CONTROL_S:
        fileio.EditorSave(appData)
        break
    case consts.LEFT_ARROW, consts.RIGHT_ARROW, consts.UP_ARROW, consts.DOWN_ARROW:
        editorMoveCursor(appData, keyReadRune)
        break
    case consts.PAGE_UP, consts.PAGE_DOWN:
        if keyReadRune == consts.PAGE_UP {
            appData.CursorPosY = appData.RowOffSet
        } else if keyReadRune == consts.PAGE_DOWN {
            appData.CursorPosY = appData.RowOffSet+appData.ScreenRows-1
            if appData.CursorPosY > appData.NumRows {
                appData.CursorPosY = appData.NumRows
            }
        }
        times := appData.ScreenRows
        for times>0 {
            directionRune := consts.DOWN_ARROW
            if keyReadRune == consts.PAGE_UP {
                directionRune = consts.UP_ARROW
            }
            editorMoveCursor(appData, rune(directionRune))
            times--
        }
        break
    case consts.HOME_KEY:
        appData.CursorPosX = 0
        break
    case consts.END_KEY:
        // appData.CursorPosX = appData.ScreenColumns-1
        if appData.CursorPosY < appData.NumRows{
            appData.CursorPosX = appData.Row[appData.CursorPosY].Size
        }
        break
    case consts.NOTHINGKEY:
        break
    default:
        // appData.ABuf.WriteRune(keyReadRune)
        // fmt.Println(keyReadRune)
        // disableRawMode(&appData.OldTerminalState)
        // defer os.Exit(0)
        data.EditorInsertChar(appData, keyReadRune)
        break
    }
}

func editorMoveCursor(appData *data.EditorConfig, inputRune rune){
    // We and implementing moving at new line by hitting left on the end of the line
    // Or right at the end of the line
    var pointerRow *data.EditorRow
    // Should it be appData.NumRows-1?
    if appData.CursorPosY >= appData.NumRows {
        pointerRow = nil
    } else {
        pointerRow = appData.Row[appData.CursorPosY]
    }

    switch inputRune {
        case consts.LEFT_ARROW:
            if (appData.CursorPosX != 0){
                appData.CursorPosX--
            }
            break
        case consts.RIGHT_ARROW:
            if pointerRow != nil && appData.CursorPosX < pointerRow.Size{
                appData.CursorPosX++
            }
            break
        case consts.DOWN_ARROW:
            if appData.CursorPosY < appData.NumRows-1{
                appData.CursorPosY++
            }
            break
        case consts.UP_ARROW:
            if appData.CursorPosY != 0 {
                appData.CursorPosY--
            }
            break
    }
    if appData.CursorPosY >= appData.NumRows {
        pointerRow = nil
    } else {
        pointerRow = appData.Row[appData.CursorPosY]
    }
    rowLength := 0
    if pointerRow != nil {
        rowLength = pointerRow.Size
    }
    if appData.CursorPosX > rowLength {
        appData.CursorPosX = rowLength
    }
}

func editorDrawRows(editorData *data.EditorConfig){
    for y:=0; y<editorData.ScreenRows; y++ {
        fileRow := y + editorData.RowOffSet
        // Maybe making it > will remove the ~ at the end?
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
            // rowlength := editorData.Row[fileRow].Size - editorData.ColOffSet
            rowlength := editorData.Row[fileRow].RenderSize - editorData.ColOffSet
            if rowlength < 0 {
                rowlength = 0
            }
            if rowlength > editorData.ScreenColumns {
                rowlength = editorData.ScreenColumns
                // shortenedString := (*editorData.Row[fileRow].Runes)[editorData.ColOffSet:rowlength]
                shortenedString := (*editorData.Row[fileRow].Render)[editorData.ColOffSet:rowlength]
                editorData.ABuf.WriteString(shortenedString) 
            } else if rowlength == 0 {
                editorData.ABuf.WriteString("")
            } else {
                editorData.ABuf.WriteString((*editorData.Row[fileRow].Render)[editorData.ColOffSet:])
            }

        }
        editorData.ABuf.WriteString("\033[K")
        // Maybe this will fix the ~ at the end of file
        // if y < editorData.ScreenRows -1 {
        //     editorData.ABuf.WriteString("\r\n")
        // }
        editorData.ABuf.WriteString("\r\n")
    }
}

func initEditor(oldState *term.State) data.EditorConfig{
    width, height := terminal.GetWindowSize()
    height -= 2 //Making space for status bar
    initCursorX, initCursorY := 0,0
    var newBuf strings.Builder
    newBuf.Reset()
    newRowSlice := make([]*data.EditorRow, 0)
    return data.EditorConfig{
        OldTerminalState: *oldState, 
        ScreenRows: height,
        RowOffSet: 0,
        ColOffSet: 0,
        RendorIndexX: 0,
        ScreenColumns: width,
        CursorPosX: initCursorX,
        CursorPosY: initCursorY,
        ABuf: newBuf,
        Row: newRowSlice,
        NumRows: 0,
        FileName: nil,
        StatusMessage: "",
        StatusMessageTime: time.Now(),
        Dirty: 0,
    }
}

func setCursorPosition(data *data.EditorConfig) {
    cursorPos := fmt.Sprintf("\033[%d;%dH", (data.CursorPosY - data.RowOffSet)+1, (data.RendorIndexX - data.ColOffSet)+1)
    data.ABuf.WriteString(cursorPos)
}

