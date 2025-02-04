package fileio

import (
	"bufio"
	"fmt"
	"os"
	"wicho/whim/app/consts"
	"wicho/whim/app/highlight"
	"wicho/whim/app/input"
	"wicho/whim/app/output"
	"wicho/whim/app/row"
)

func EditorOpen(appData *consts.EditorConfig, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		appData.Die()
	}
	defer file.Close()
	appData.FileName = &fileName

	highlight.EditorSelectSyntaxHighlight(appData)

	line := ""
	totalLines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalLines++
		line = scanner.Text()
		row.EditorInsertRow(appData, appData.NumRows, line)
	}
	appData.Dirty = 0
}

func EditorSave(appData *consts.EditorConfig) {
	if appData.FileName == nil {
		appData.FileName = input.EditorPrompt(appData, "(ESC to cancel) Save as: ", nil)
		if appData.FileName == nil {
			output.EditorSetStatusMessage(appData, "Save Aborted")
			return
		}
		highlight.EditorSelectSyntaxHighlight(appData)
	}

	fileIntoString := editorRowsToString(appData)

	file, err := os.OpenFile(*appData.FileName, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		output.EditorSetStatusMessage(appData, "Can't save! I/O error: File Open")
		return
	}
	err = file.Truncate(0)
	if err != nil {
		output.EditorSetStatusMessage(appData, "Can't save! I/O error: File Truncate")
		return
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		output.EditorSetStatusMessage(appData, "Can't save! I/O error: File Seek")
		return
	}
	bytesWritten, err := file.WriteString(fileIntoString)
	if err != nil {
		output.EditorSetStatusMessage(appData, "Can't save! I/O error: File Write")
		return
	}
	msg := fmt.Sprintf("%d bytes written to disk", bytesWritten)
	appData.Dirty = 0
	output.EditorSetStatusMessage(appData, msg)
}

func editorRowsToString(appData *consts.EditorConfig) string {
	fileAsString := ""
	for _, row := range appData.Row {
		fileAsString += *row.Runes + string('\n')
	}
	return fileAsString
}
