package consts

import (
    "fmt"
    "os"
    "slices"
    "strings"
	"time"

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
    CONTROL_F = 1014
    WHIM_VERSION = "0.0.1"
    HL_NORMAL = 0
    HL_NUMBER = 1
    HL_MATCH = 2
)

type EditorConfig struct {
    OldTerminalState term.State
    ScreenRows int
    ScreenColumns int
    RowOffSet int
    ColOffSet int
    CursorPosX int
    CursorPosY int
    RendorIndexX int
    NumRows int
    Row []*EditorRow
    FileName *string
    ABuf strings.Builder
    StatusMessage string
    StatusMessageTime time.Time
    StringFindData *FindData
    Dirty int
}

// Below is a struct that will be used to hold temp search data
type FindData struct{
    LastMatch int    
    Direction int
    SavedHlLine int
    SavedHighlights []int
}

type EditorRow struct {
    Size int
    Runes *string
    Render *string
    RenderSize int
    Highlights []int
}

func (appData *EditorConfig) Die(){
    fmt.Print("\033[2J")
    fmt.Print("\033[H")
    //Add os exit method! 
    term.Restore(0, &appData.OldTerminalState)
    defer os.Exit(1)
}

func RuneIsCtrlKey(aRune rune) bool {
    var sliceOfRunes = []rune{
        NOTHINGKEY, CONTROLCASCII, CONTROLFIRSTBYTE, LEFT_ARROW,
        RIGHT_ARROW, DOWN_ARROW, UP_ARROW, PAGE_UP, PAGE_DOWN, 
        HOME_KEY, END_KEY, DEL_KEY, ESC, BACKSPACE, CONTROL_S, CONTROL_F,
    }
    return slices.Contains(sliceOfRunes, aRune)
}
