package terminal

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"

	"wicho/whim/app/consts"
	"wicho/whim/input"
)

//Die
func Die(){
    fmt.Print("\033[2J")
    fmt.Print("\033[H")
    //Add os exit method!
    defer os.Exit(1)
}

//DisableRawMode
func DisableRawMode(oldState *term.State) {
    term.Restore(0, oldState)
}

//EnableRowMode
func EnableRawMode() term.State{
    oldState, err := term.MakeRaw(0)
    if err != nil {
        panic(err)
    }
    return *oldState
}

//EditorReadKey
func EditorReadKey() rune {
    inputBytes := input.ReaderBytes()
    inputRune, _ := utf8.DecodeRune(inputBytes)
    if unicode.IsControl(inputRune){
        //fmt.Println("--bytes-",inputBytes, "-Rune-", inputRune, "-String-", string(inputRune))
        //os.Exit(1)
        switch inputRune {
        case 3://CTRL-C
            return consts.CONTROLCASCII
        case 27://First byte is A CTRL byte
            // Should add another switch statement that deals with this?
            returnRune, _ := utf8.DecodeRune(inputBytes[2:])
            switch returnRune{
            case 53://PAGE UP
                return consts.PAGE_UP
            case 54://PAGE DOWN
                return consts.PAGE_DOWN
            case 68://LEFT ARROW
                return consts.LEFT_ARROW
            case 67://RIGHT ARROW
                return consts.RIGHT_ARROW
            case 66://DOWN ARROW
                return consts.DOWN_ARROW
            case 65://UP ARROW
                return consts.UP_ARROW
            case 72://HOME KEY
                return consts.HOME_KEY
            case 70://END KEY
                return consts.END_KEY
            case 51://DEL KEY
                return consts.DEL_KEY
            case 0://ESC KEY
                return consts.ESC
            }
            return consts.NOTHINGKEY
        case 127://BACKSPACE
            return consts.BACKSPACE
        case 13://ENTER
            return '\r'
        case 19://Ctrl-s
            return consts.CONTROL_S
        default:
            return consts.NOTHINGKEY
        }
    }
    switch string(inputRune) {
    case "q":
        return consts.CONTROLCASCII
    case "h":
        return consts.LEFT_ARROW
    case "j":
        return consts.DOWN_ARROW
    case "k":
        return consts.UP_ARROW
    case "l":
        return consts.RIGHT_ARROW
    }
    return inputRune
}

//GetCursorPosition

//GetWindowSize
func GetWindowSize()(int, int){
    width, height, err := term.GetSize(0)
    if err != nil {
        Die()
    }
    return width, height
}

