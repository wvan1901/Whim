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
	NOTHINGKEY       = 999 //This key will do nothing when pressed
	CONTROLCASCII    = 1000
	CONTROLFIRSTBYTE = 1001
	LEFT_ARROW       = 1002
	RIGHT_ARROW      = 1003
	DOWN_ARROW       = 1004
	UP_ARROW         = 1005
	PAGE_UP          = 1006
	PAGE_DOWN        = 1007
	HOME_KEY         = 1008
	END_KEY          = 1009
	DEL_KEY          = 1010
	ESC              = 1011
	BACKSPACE        = 1012
	CONTROL_S        = 1013
	CONTROL_F        = 1014
	WHIM_VERSION     = "0.0.1"
	HL_NORMAL        = 0
	HL_NUMBER        = 1
	HL_MATCH         = 2
	HL_STRING        = 3
	HL_COMMENT       = 4
	HL_KEYWORD1      = 5
	HL_KEYWORD2      = 6
	HL_MLCOMMENT     = 7
)

type EditorConfig struct {
	OldTerminalState  term.State
	ScreenRows        int
	ScreenColumns     int
	RawScreenColumns  int
	RowOffSet         int
	ColOffSet         int
	CursorPosX        int
	CursorPosY        int
	RendorIndexX      int
	NumRows           int
	Row               []*EditorRow
	Dirty             int
	FileName          *string
	ABuf              strings.Builder
	StatusMessage     string
	StatusMessageTime time.Time
	StringFindData    *FindData
	EditorSyntax      *EditorSyntax
	Features          *Features
	Mode              Mode
}

// Below is a struct that will be used to hold temp search data
type FindData struct {
	LastMatch       int
	Direction       int
	SavedHlLine     int
	SavedHighlights []int
}

type EditorRow struct {
	Idx           int
	Size          int
	Runes         *string
	Render        *string
	RenderSize    int
	Highlights    []int
	HlOpenComment bool
}

type EditorSyntax struct {
	Filetype               string
	Filematch              []string
	Keywords               []string
	SinglelineCommentStart string
	MultilineCommentStart  string
	MultilineCommentEnd    string
	Flags                  []string
}

type Features struct {
	LineNumberOn       bool
	RelativeLineNumber bool
	LineNumberIndent   int
}

func (appData *EditorConfig) Die() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	term.Restore(0, &appData.OldTerminalState)
	defer os.Exit(1)
}

func (appData *EditorConfig) Quit() {
	//TODO: When terminal module is refactored call func here
}

func RuneIsCtrlKey(aRune rune) bool {
	var sliceOfRunes = []rune{
		NOTHINGKEY, CONTROLCASCII, CONTROLFIRSTBYTE, LEFT_ARROW,
		RIGHT_ARROW, DOWN_ARROW, UP_ARROW, PAGE_UP, PAGE_DOWN,
		HOME_KEY, END_KEY, DEL_KEY, ESC, BACKSPACE, CONTROL_S, CONTROL_F,
	}
	return slices.Contains(sliceOfRunes, aRune)
}

func HLDB() []EditorSyntax {
	cInfo := EditorSyntax{
		Filetype:  "c",
		Filematch: []string{".c", ".h", ".cpp"},
		Keywords: []string{"swtich", "if", "while", "for", "break", "continue", "return", "else",
			"struct", "union", "typedef", "static", "enum", "class", "case",
			"int|", "long|", "double|", "float|", "char|", "unsigned|", "signed|", "void|",
		},
		SinglelineCommentStart: "//",
		MultilineCommentStart:  "/*",
		MultilineCommentEnd:    "*/",
		Flags:                  []string{"HL_HIGHLIGHT_NUMBERS", "HL_HIGHLIGHT_STRINGS"},
	}

	var bd []EditorSyntax
	bd = append(bd, cInfo)
	return bd
}

type Mode interface {
	EditorProcessKeyPress(*EditorConfig)
	ShortString() string
}
