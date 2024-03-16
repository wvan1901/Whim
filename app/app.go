package app

import (
	"os"
	"wicho/whim/input"

	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
    CONTROLCASCII = 3
)

func RunApp(){
    oldState := enableRawMode()
    defer disableRawMode(&oldState)
    for {
        editorRefreshScreen()
        editorProcessKeyPress(&oldState)
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

func editorReadKey() byte {
    return input.ReaderByte()//inputByte
}

func editorProcessKeyPress(oldState *term.State){
    inputRune, _ := utf8.DecodeRune([]byte(string(editorReadKey())))
    if unicode.IsControl(inputRune){
        switch inputRune {
        case CONTROLCASCII:
            fmt.Println("<C>")
            disableRawMode(oldState)
            fmt.Print("\033[2J")
            fmt.Print("\033[H")
            defer os.Exit(0)
            break
        }
    }
    switch string(inputRune) {
    case "q":
        fmt.Println("<q>")
    default: fmt.Print(string(inputRune))
    }
}

func editorRefreshScreen(){
    fmt.Print("\033[2J")
    fmt.Print("\033[H")
    editorDrawRows()
    fmt.Print("\033[H")
}

func die(){
    fmt.Print("\033[2J")
    fmt.Print("\033[H")
    //Add os exit method!
    defer os.Exit(1)
}

func editorDrawRows(){
    y:= 0
    for y=0;y<24;y++ {
        fmt.Print("~\r\n")
    }
}

