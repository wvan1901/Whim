package app

import (
	"os"
	"strings"
	"time"
	"wicho/whim/app/consts"
	"wicho/whim/app/fileio"
	"wicho/whim/app/mode"
	"wicho/whim/app/output"
	"wicho/whim/app/terminal"

	"golang.org/x/term"
)

// A Solution to errors with no file, is check for input once it get ones then we add a row
// This could be a mode
func RunApp() {
	oldState := terminal.EnableRawMode()
	defer terminal.DisableRawMode(&oldState)
	AppData := initEditor(&oldState)
	argsWithProg := os.Args
	if len(argsWithProg) >= 2 {
		fileio.EditorOpen(&AppData, argsWithProg[1])
	}
	output.EditorSetStatusMessage(&AppData, "HELP: Ctrl-s = save | Ctrl-C = Quit | Ctrl-F = Find")
	for {
		output.EditorRefreshScreen(&AppData)
		//editorProcessKeyPress(&AppData)
		AppData.Mode.EditorProcessKeyPress(&AppData)
	}
}

func initEditor(oldState *term.State) consts.EditorConfig {
	width, height := terminal.GetWindowSize()
	height -= 2 //Making space for status bar
	initCursorX, initCursorY := 0, 0
	var newBuf strings.Builder
	newBuf.Reset()
	newRowSlice := make([]*consts.EditorRow, 0)
	mode := mode.Normal{}
	features := consts.Features{
		LineNumberOn:       true,
		RelativeLineNumber: true,
		LineNumberIndent:   0,
	}
	return consts.EditorConfig{
		OldTerminalState:  *oldState,
		ScreenRows:        height,
		RowOffSet:         0,
		ColOffSet:         0,
		RendorIndexX:      0,
		ScreenColumns:     width,
		RawScreenColumns:  width,
		CursorPosX:        initCursorX,
		CursorPosY:        initCursorY,
		ABuf:              newBuf,
		Row:               newRowSlice,
		NumRows:           0,
		FileName:          nil,
		StatusMessage:     "",
		StatusMessageTime: time.Now(),
		StringFindData:    nil,
		Dirty:             0,
		EditorSyntax:      nil,
		Features:          &features,
		Mode:              &mode,
	}
}
