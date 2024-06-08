package app

import (
	"os"
	"strings"
	"time"
	"wicho/whim/app/consts"
	"wicho/whim/app/data"
	"wicho/whim/app/fileio"
	"wicho/whim/app/find"
	"wicho/whim/app/output"
	"wicho/whim/app/terminal"

	"fmt"

	"golang.org/x/term"
)

//A Solution to errors with no file, is check for input once it get ones then we add a row
// This could be a mode
func RunApp(){
    oldState := terminal.EnableRawMode()
    defer terminal.DisableRawMode(&oldState)
    AppData := initEditor(&oldState)
    argsWithProg := os.Args
    if len(argsWithProg) >= 2 {
        fileio.EditorOpen(&AppData, argsWithProg[1])
    }
    output.EditorSetStatusMessage(&AppData, "HELP: Ctrl-s = save | Ctrl-C, q = Quit | Ctrl-F = Find")
    for {
        output.EditorRefreshScreen(&AppData)
        editorProcessKeyPress(&AppData)
    }
}

func editorProcessKeyPress(appData *consts.EditorConfig){
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
        if appData.CursorPosY < appData.NumRows{
            appData.CursorPosX = appData.Row[appData.CursorPosY].Size
        }
        break
    case consts.CONTROL_F:
        find.EditorFind(appData)
        break
    case consts.NOTHINGKEY:
        break
    default:
        data.EditorInsertChar(appData, keyReadRune)
        break
    }
}

func editorMoveCursor(appData *consts.EditorConfig, inputRune rune){
    // We are not (for now) implementing moving at new line by hitting left on the end of the line
    // Or right at the end of the line
    var pointerRow *consts.EditorRow
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


func initEditor(oldState *term.State) consts.EditorConfig{
    width, height := terminal.GetWindowSize()
    height -= 2 //Making space for status bar
    initCursorX, initCursorY := 0,0
    var newBuf strings.Builder
    newBuf.Reset()
    newRowSlice := make([]*consts.EditorRow, 0)
    return consts.EditorConfig{
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
        StringFindData: nil,
        Dirty: 0,
        EditorSyntax: nil,
    }
}

