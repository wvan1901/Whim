package mode

import (
	"wicho/whim/app/consts"
	"wicho/whim/app/data"
	"wicho/whim/app/fileio"
	"wicho/whim/app/find"
	"wicho/whim/app/terminal"
)

type Normal struct{}
type Insert struct{}

func (n *Normal) EditorProcessKeyPress(c *consts.EditorConfig) {
	keyReadRune := terminal.EditorReadKey()
	switch keyReadRune {
	case '\r':
		editorMoveCursor(c, consts.DOWN_ARROW)
		break
	case consts.BACKSPACE, consts.DEL_KEY:
		editorMoveCursor(c, consts.LEFT_ARROW)
		break
	case consts.ESC:
		break
	case consts.CONTROLCASCII:
		//TODO: Add a warning message if file is dirty and ask User to Confirm
		terminal.Quit(&c.OldTerminalState)
		break
	case consts.CONTROL_S:
		fileio.EditorSave(c)
		break
	case consts.LEFT_ARROW, consts.RIGHT_ARROW, consts.UP_ARROW, consts.DOWN_ARROW:
		editorMoveCursor(c, keyReadRune)
		break
	case consts.PAGE_UP, consts.PAGE_DOWN:
		if keyReadRune == consts.PAGE_UP {
			c.CursorPosY = c.RowOffSet
		} else if keyReadRune == consts.PAGE_DOWN {
			c.CursorPosY = c.RowOffSet + c.ScreenRows - 1
			if c.CursorPosY > c.NumRows {
				c.CursorPosY = c.NumRows
			}
		}
		times := c.ScreenRows
		for times > 0 {
			directionRune := consts.DOWN_ARROW
			if keyReadRune == consts.PAGE_UP {
				directionRune = consts.UP_ARROW
			}
			editorMoveCursor(c, rune(directionRune))
			times--
		}
		break
	case consts.HOME_KEY:
		c.CursorPosX = 0
		break
	case consts.END_KEY:
		if c.CursorPosY < c.NumRows {
			c.CursorPosX = c.Row[c.CursorPosY].Size
		}
		break
	case consts.CONTROL_F:
		find.EditorFind(c)
		break
	case consts.NOTHINGKEY:
		break
	}
	switch string(keyReadRune) {
	case "h":
		editorMoveCursor(c, consts.LEFT_ARROW)
	case "j":
		editorMoveCursor(c, consts.DOWN_ARROW)
	case "k":
		editorMoveCursor(c, consts.UP_ARROW)
	case "l":
		editorMoveCursor(c, consts.RIGHT_ARROW)
	case "i":
		mode := Insert{}
		c.Mode = &mode
		return
	case ":":
		mode := Command{}
		c.Mode = &mode
	}
}

func (n *Normal) ShortString() string {
	return "N"
}

func (n *Insert) EditorProcessKeyPress(c *consts.EditorConfig) {
	keyReadRune := terminal.EditorReadKey()
	switch keyReadRune {
	case '\r':
		data.EditorInsertNewLine(c)
		break
	case consts.BACKSPACE, consts.DEL_KEY:
		data.EditorDelChar(c)
		break
	case consts.ESC, consts.CONTROLCASCII:
		mode := Normal{}
		c.Mode = &mode
		break
	case consts.LEFT_ARROW, consts.RIGHT_ARROW, consts.UP_ARROW, consts.DOWN_ARROW:
		editorMoveCursor(c, keyReadRune)
		break
	case consts.NOTHINGKEY:
		break
	default:
		data.EditorInsertChar(c, keyReadRune)
		break
	}
}

func editorMoveCursor(appData *consts.EditorConfig, inputRune rune) {
	// We are not (for now) implementing moving at new line by hitting left on the end of the line
	// Or right at the end of the line
	var pointerRow *consts.EditorRow
	if appData.CursorPosY >= appData.NumRows {
		pointerRow = nil
	} else {
		pointerRow = appData.Row[appData.CursorPosY]
	}

	switch inputRune {
	case consts.LEFT_ARROW:
		if appData.CursorPosX != 0 {
			appData.CursorPosX--
		}
		break
	case consts.RIGHT_ARROW:
		if pointerRow != nil && appData.CursorPosX < pointerRow.Size {
			appData.CursorPosX++
		}
		break
	case consts.DOWN_ARROW:
		if appData.CursorPosY < appData.NumRows-1 {
			appData.CursorPosY++
		}
		break
	case consts.UP_ARROW:
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

func (n *Insert) ShortString() string {
	return "I"
}
