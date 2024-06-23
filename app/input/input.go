package input

import (
	"wicho/whim/app/consts"
	"wicho/whim/app/output"
	"wicho/whim/app/terminal"
)

func EditorPrompt(appData *consts.EditorConfig, prompt string, aFunc func(data *consts.EditorConfig, query string, b rune)) *string {
	buf := ""

	for {
		output.EditorSetStatusMessage(appData, prompt, buf)
		output.EditorRefreshScreen(appData)

		inputRune := terminal.EditorReadKey()
		if inputRune == consts.BACKSPACE {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
			}
		} else if inputRune == consts.ESC {
			output.EditorSetStatusMessage(appData, "")
			if aFunc != nil {
				aFunc(appData, buf, inputRune)
			}
			buf = ""
			return nil
		} else if inputRune == '\r' {
			if len(buf) != 0 {
				output.EditorSetStatusMessage(appData, "")
				if aFunc != nil {
					aFunc(appData, buf, inputRune)
				}
				return &buf
			}
		} else if !consts.RuneIsCtrlKey(inputRune) {
			buf += string(inputRune)
		}
		if aFunc != nil {
			aFunc(appData, buf, inputRune)
		}
	}
}
