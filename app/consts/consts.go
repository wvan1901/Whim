package consts

import "slices"

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
)

func RuneIsCtrlKey(aRune rune) bool {
    var sliceOfRunes = []rune{
        NOTHINGKEY, CONTROLCASCII, CONTROLFIRSTBYTE, LEFT_ARROW,
        RIGHT_ARROW, DOWN_ARROW, UP_ARROW, PAGE_UP, PAGE_DOWN, 
        HOME_KEY, END_KEY, DEL_KEY, ESC, BACKSPACE, CONTROL_S, CONTROL_F,
    }
    return slices.Contains(sliceOfRunes, aRune)
}
