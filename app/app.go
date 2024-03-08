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
    /*
    var inputByte byte
    for string(inputByte) != "q"{
        inputRune, _ := utf8.DecodeRune([]byte(string(inputByte)))
        //If the Char is a control char then it enter if statement
        if unicode.IsControl(inputRune) {
            if inputByte == CONTROLCASCII {
                break
            }
        }
        inputByte = input.ReaderByte()

        //fmt.Print(string(inputByte))
    }
    */
    for {
        editorProcessKeyPress(&oldState)
    }
}

func enableRawMode() term.State{
    /*
    STDINFILENO := 0
    raw, err := unix.IoctlGetTermios(STDINFILENO, unix.TCGETS)
    if err != nil {
        panic(err)
    }
    
    rawState := *raw
    rawState.Lflag &^= unix.ECHO err = unix.IoctlSetTermios(STDINFILENO, unix.TCSETA, &rawState) if err != nil {
        panic(err)
    }
    */
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
    // var inputByte byte
    // inputByte = input.ReaderByte()
    return input.ReaderByte()//inputByte
}

func editorProcessKeyPress(oldState *term.State){
    inputRune, _ := utf8.DecodeRune([]byte(string(editorReadKey())))
    if unicode.IsControl(inputRune){
        switch inputRune {
        case CONTROLCASCII:
            fmt.Println("<C>")
            disableRawMode(oldState)
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
