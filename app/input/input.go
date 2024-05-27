package input

import (
	"wicho/whim/app/consts"
	"wicho/whim/app/data"
	"wicho/whim/app/output"
	"wicho/whim/app/terminal"
)

func EditorPrompt(appData *data.EditorConfig, prompt string, aFunc func(data *data.EditorConfig, query string, b rune)) *string {
    buf := ""

    for {
        appData.EditorSetStatusMessage(prompt, buf)
        output.EditorRefreshScreen(appData)

        inputRune := terminal.EditorReadKey()
        if inputRune == consts.BACKSPACE {
            if len(buf) != 0 {
                buf = buf[:len(buf)-1]
            }
        } else if inputRune == consts.ESC {
            appData.EditorSetStatusMessage("")
            if aFunc != nil {
                aFunc(appData, buf, inputRune)
            }
            buf = ""
            return nil
        } else if inputRune == '\r' {
            if len(buf) != 0 {
                appData.EditorSetStatusMessage("")
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

