package input

import (
	"wicho/whim/app/consts"
	"wicho/whim/app/data"
	"wicho/whim/app/output"
	"wicho/whim/app/terminal"

	"unicode"
)

func EditorPrompt(appData *data.EditorConfig, prompt string) *string {
    // bufSize := 128
    buf := ""
    // bufLen := 0

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
            buf = ""
            return nil
        } else if inputRune == '\r' {
            if len(buf) != 0 {
                appData.EditorSetStatusMessage("")
                return &buf
            }
        } else if !unicode.IsControl(inputRune) {
            buf += string(inputRune)
        }
    }
}
