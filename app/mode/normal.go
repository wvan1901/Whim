package mode

import (
	"wicho/whim/app/consts"
	"wicho/whim/app/input"
	"wicho/whim/app/terminal"
)

func normalExit(c *consts.EditorConfig) {
	if c.Dirty == 0 {
		//Exit App
		terminal.Quit(&c.OldTerminalState)
		return
	}
	exitApp := confirmExit(c)
	if exitApp {
		terminal.Quit(&c.OldTerminalState)
		return
	}
}

func confirmExit(c *consts.EditorConfig) bool {
	queryPrompt := input.EditorPrompt(c, "File is dirty! exit? (y/n)", nil)
	if queryPrompt == nil {
		return false
	}
	switch *queryPrompt {
	case "y", "Y":
		return true
	}
	return false
}
