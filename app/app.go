package app

import (
	"os"
	"strings"
	"time"
	"wicho/whim/app/data"
	"wicho/whim/app/fileio"
	"wicho/whim/input"

	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
    NOTHINGKEY= 999 //This key will do nothing when pressed
    CONTROLCASCII = 1000
    CONTROLFIRSTBYTE = 1001
    LEFT_ARROW = 1002
    RIGHT_ARROW = 1003
    DOWN_ARROW = 1004
    UP_ARROW = 1005
    PAGE_UP = 1006
    PAGE_DOWN = 1007
    HOME_KEY = 1008
    END_KEY = 1009
    DEL_KEY = 1010
    ESC = 1011
    BACKSPACE = 1012
    CONTROL_S = 1013
    WHIM_VERSION = "0.0.1"
)

func RunApp(){
    oldState := enableRawMode()
    defer disableRawMode(&oldState)
    AppData := initEditor(&oldState)
    argsWithProg := os.Args
    if len(argsWithProg) >= 2 {
        fileio.EditorOpen(&AppData, argsWithProg[1])
    }
    AppData.EditorSetStatusMessage("HELP: Ctrl-s = save | Ctrl-C, q = Quit")
    for {
        editorRefreshScreen(&AppData)
        editorProcessKeyPress(&AppData)
    }
}

func enableRawMode() term.State{
    oldState, err := term.MakeRaw(0)
    if err != nil {
        panic(err)
    }
    return *oldState
}

func disableRawMode(oldState *term.State) {
    term.Restore(0, oldState)
}

func editorReadKey() rune {
    inputBytes := input.ReaderBytes()
    inputRune, _ := utf8.DecodeRune(inputBytes)
    if unicode.IsControl(inputRune){
        //fmt.Println("--bytes-",inputBytes, "-Rune-", inputRune, "-String-", string(inputRune))
        //os.Exit(1)
        switch inputRune {
        case 3://CTRL-C
            return CONTROLCASCII
        case 27://First byte is A CTRL byte
            // Should add another switch statement that deals with this?
            returnRune, _ := utf8.DecodeRune(inputBytes[2:])
            switch returnRune{
            case 53://PAGE UP
                return PAGE_UP
            case 54://PAGE DOWN
                return PAGE_DOWN
            case 68://LEFT ARROW
                return LEFT_ARROW
            case 67://RIGHT ARROW
                return RIGHT_ARROW
            case 66://DOWN ARROW
                return DOWN_ARROW
            case 65://UP ARROW
                return UP_ARROW
            case 72://HOME KEY
                return HOME_KEY
            case 70://END KEY
                return END_KEY
            case 51://DEL KEY
                return DEL_KEY
            case 0://ESC KEY
                return ESC
            }
            return NOTHINGKEY
        case 127://BACKSPACE
            return BACKSPACE
        case 13://ENTER
            return '\r'
        case 19://Ctrl-s
            return CONTROL_S
        default:
            return NOTHINGKEY
        }
    }
    switch string(inputRune) {
    case "q":
        return CONTROLCASCII
    case "h":
        return LEFT_ARROW
    case "j":
        return DOWN_ARROW
    case "k":
        return UP_ARROW
    case "l":
        return RIGHT_ARROW
    }
    return inputRune
}

func editorProcessKeyPress(appData *data.EditorConfig){
    keyReadRune := editorReadKey()
    switch keyReadRune {
    case '\r':
        //TODO: Implement Enter
        break
    case BACKSPACE, DEL_KEY:
        //TODO: Implement Backspace
        data.EditorDelChar(appData)
        break
    case ESC:
        break
    case CONTROLCASCII:
        //TODO: Add a warning message if file is dirty and ask User to Confirm
        fmt.Println("<C>")
        disableRawMode(&appData.OldTerminalState)
        fmt.Print("\033[2J")
        fmt.Print("\033[H")
        defer os.Exit(0)
        break
    case CONTROL_S:
        fileio.EditorSave(appData)
        break
    case LEFT_ARROW, RIGHT_ARROW, UP_ARROW, DOWN_ARROW:
        editorMoveCursor(appData, keyReadRune)
        break
    case PAGE_UP, PAGE_DOWN:
        if keyReadRune == PAGE_UP {
            appData.CursorPosY = appData.RowOffSet
        } else if keyReadRune == PAGE_DOWN {
            appData.CursorPosY = appData.RowOffSet+appData.ScreenRows-1
            if appData.CursorPosY > appData.NumRows {
                appData.CursorPosY = appData.NumRows
            }
        }
        times := appData.ScreenRows
        for times>0 {
            directionRune := DOWN_ARROW
            if keyReadRune == PAGE_UP {
                directionRune = UP_ARROW
            }
            editorMoveCursor(appData, rune(directionRune))
            times--
        }
        break
    case HOME_KEY:
        appData.CursorPosX = 0
        break
    case END_KEY:
        // appData.CursorPosX = appData.ScreenColumns-1
        if appData.CursorPosY < appData.NumRows{
            appData.CursorPosX = appData.Row[appData.CursorPosY].Size
        }
        break
    case NOTHINGKEY:
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
        case LEFT_ARROW:
            if (appData.CursorPosX != 0){
                appData.CursorPosX--
            }
            break
        case RIGHT_ARROW:
            if pointerRow != nil && appData.CursorPosX < pointerRow.Size{
                appData.CursorPosX++
            }
            break
        case DOWN_ARROW:
            if appData.CursorPosY < appData.NumRows-1{
                appData.CursorPosY++
            }
            break
        case UP_ARROW:
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

func editorRefreshScreen(appData *data.EditorConfig){
    editorScroll(appData)
    appData.ABuf.WriteString("\033[?25l")
    appData.ABuf.WriteString("\033[H")
    editorDrawRows(appData)
    editorDrawStatusBar(appData)
    editorDrawMessageBar(appData)
    
    setCursorPosition(appData)

    appData.ABuf.WriteString("\033[?25h")

    fmt.Print(appData.ABuf.String())
    appData.ABuf.Reset()
}

func die(){
    fmt.Print("\033[2J")
    fmt.Print("\033[H")
    //Add os exit method!
    defer os.Exit(1)
}

func editorDrawRows(editorData *data.EditorConfig){
    for y:=0; y<editorData.ScreenRows; y++ {
        fileRow := y + editorData.RowOffSet
        // Maybe making it > will remove the ~ at the end?
        if fileRow >= editorData.NumRows{
            if editorData.NumRows == 0 && y == editorData.ScreenRows/3 {
                welcome := fmt.Sprintf("Whim Editor -- version %s", WHIM_VERSION)
                if len(welcome) > editorData.ScreenColumns{
                    welcome = "v:"+WHIM_VERSION
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

func editorScroll(editorData *data.EditorConfig) {
    editorData.RendorIndexX = 0
    if editorData.CursorPosY < editorData.NumRows {
        editorData.RendorIndexX = data.EditorRowCxToRx(editorData.Row[editorData.CursorPosY], editorData.CursorPosX)
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

func editorDrawStatusBar(appData *data.EditorConfig){
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
        rightSideStatus = fmt.Sprintf("%d/%d", appData.CursorPosY+1, appData.NumRows)
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

func editorDrawMessageBar(appData *data.EditorConfig){
    appData.ABuf.WriteString("\033[K")
    messageLength := len(appData.StatusMessage)
    if messageLength > appData.ScreenColumns{
        messageLength = appData.ScreenColumns
    }
    if messageLength>0 && (time.Since(appData.StatusMessageTime) < 5*time.Second){//1 sec = 1000000000 nano sec
        appData.ABuf.WriteString(appData.StatusMessage[:messageLength])
    }

}

func getWindowSize()(int, int){
    width, height, err := term.GetSize(0)
    if err != nil {
        die()
    }
    return width, height
}

func initEditor(oldState *term.State) data.EditorConfig{
    width, height := getWindowSize()
    height -= 2 //Making space for status bar
    initCursorX, initCursorY := 1,1
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

