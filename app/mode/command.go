package mode

import (
	"wicho/whim/app/consts"
	"wicho/whim/app/fileio"
	"wicho/whim/app/input"
	"wicho/whim/app/terminal"
)

type Command struct{}

func (n *Command) EditorProcessKeyPress(c *consts.EditorConfig) {
	getCommand(c)
}

func (n *Command) ShortString() string {
	return "C"
}

func getCommand(appData *consts.EditorConfig) {
	savedCursorPosX := appData.CursorPosX
	savedCursorPosY := appData.CursorPosY
	savedColOffset := appData.ColOffSet
	savedRowOffset := appData.RowOffSet

	queryPrompt := input.EditorPrompt(appData, ":", nil)
	if queryPrompt == nil {
		appData.Mode = &Normal{}
		return
	} else {
		// It returns the buffer
		appData.CursorPosX = savedCursorPosX
		appData.CursorPosY = savedCursorPosY
		appData.ColOffSet = savedColOffset
		appData.RowOffSet = savedRowOffset
		executeCommand(*queryPrompt, appData)
	}
	appData.Mode = &Normal{}
}

func executeCommand(cmd string, c *consts.EditorConfig) {
	switch cmd[0] {
	case 'w': // Write File
		fileio.EditorSave(c)
		return
	case 'q': // Quit App
		terminal.Quit(&c.OldTerminalState)
		return
	}
}
