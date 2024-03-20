package app

import (
	"os"
    "strings"
	"wicho/whim/app/data"
	"wicho/whim/input"

	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
    CONTROLCASCII = 3
    ARROWFIRSTBYTE = 27
    LEFT_ARROW = 68
    RIGHT_ARROW = 67
    DOWN_ARROW = 66
    UP_ARROW = 65
    WHIM_VERSION = "0.0.1"
)

func RunApp(){
    oldState := enableRawMode()
    defer disableRawMode(&oldState)
    AppData := initEditor(&oldState)
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

func editorReadKey() []byte {
    return input.ReaderBytes()// gets the array of bytes
}

func editorProcessKeyPress(appData *data.EditorConfig){
    keyRead := editorReadKey()
    inputRune, _ := utf8.DecodeRune([]byte(string(keyRead)))
    if unicode.IsControl(inputRune){
        switch keyRead[0] {
        case CONTROLCASCII:
            fmt.Println("<C>")
            disableRawMode(&appData.OldTerminalState)
            fmt.Print("\033[2J")
            fmt.Print("\033[H")
            defer os.Exit(0)
        break
        case ARROWFIRSTBYTE:
            editorMoveCursor(appData, keyRead[2])
        default:
            appData.ABuf.WriteRune(inputRune)
        }
    }
    switch string(inputRune) {
    case "q":
        fmt.Println("<q>")
    case "h", "j", "k", "l":
        editorMoveCursor(appData, keyRead[0])
    default: fmt.Print(string(inputRune))
    }
}

func editorMoveCursor(appData *data.EditorConfig, inputByte byte){
    switch string(inputByte) {
        case "h"://Left
            appData.CursorPosX--
        case "j"://Down
            appData.CursorPosY++
        case "k"://Up
            appData.CursorPosY--
        case "l"://Right
            appData.CursorPosX++
    }
    switch inputByte {
        case LEFT_ARROW:
            appData.CursorPosX--
        case RIGHT_ARROW:
            appData.CursorPosX++
        case DOWN_ARROW:
            appData.CursorPosY++
        case UP_ARROW:
            appData.CursorPosY--
    }
}

func editorRefreshScreen(appData *data.EditorConfig){
    appData.ABuf.WriteString("\033[?25l")
    appData.ABuf.WriteString("\033[H")
    editorDrawRows(appData)
    
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
        if y == editorData.ScreenRows/3 {
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

        editorData.ABuf.WriteString("\033[K")
        if y < editorData.ScreenRows -1 {
            editorData.ABuf.WriteString("\r\n")
        }
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
    initCursorX, initCursorY := 2,2
    //fmt.Printf("\033[%d;%dH", initCursorY, initCursorX) // Set cursor position    
    var newBuf strings.Builder
    newBuf.Reset()
    return data.EditorConfig{
        OldTerminalState: *oldState, 
        ScreenRows: height,
        ScreenColumns: width,
        CursorPosX: initCursorX,
        CursorPosY: initCursorY,
        ABuf: newBuf,
    }
}

//TODO: Check if setting the cursor works!
func getCursorPosition(data *data.EditorConfig) (int, int){
    fmt.Printf("\033[%d;%dH", data.CursorPosY, data.CursorPosX) // Set cursor position    
    return data.CursorPosX, data.CursorPosY
}
func setCursorPosition(data *data.EditorConfig) {
    cursorPos := fmt.Sprintf("\033[%d;%dH", data.CursorPosY, data.CursorPosX)
    data.ABuf.WriteString(cursorPos)
}
