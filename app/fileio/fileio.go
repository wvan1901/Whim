package fileio

import (
	"bufio"
	"fmt"
	"os"
	"wicho/whim/app/data"
)

func EditorOpen(appData *data.EditorConfig, fileName string){
    file, err := os.Open(fileName)
    if err != nil {
        appData.Die()
    }
    defer file.Close()
    // appData.FileName = file.Name()
    appData.FileName = &fileName

    line := ""
    totalLines := 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        totalLines++
        line = scanner.Text()
        appData.EditorAppendRow(line)
    }
    appData.Dirty = 0
}

func EditorSave(appData *data.EditorConfig){
    if appData.FileName == nil {
        return
    }

    fileIntoString := editorRowsToString(appData)

    file, err := os.OpenFile(*appData.FileName, os.O_RDWR | os.O_CREATE, 0644)
    defer file.Close()
    if err != nil{
        appData.EditorSetStatusMessage("Can't save! I/O error: File Open")
        return
    }
    err = file.Truncate(0)
    if err != nil{
        appData.EditorSetStatusMessage("Can't save! I/O error: File Truncate")
        return
    }
    _, err = file.Seek(0,0)
    if err != nil{
        appData.EditorSetStatusMessage("Can't save! I/O error: File Seek")
        return
    }
    bytesWritten, err := file.WriteString(fileIntoString)
    if err != nil{
        appData.EditorSetStatusMessage("Can't save! I/O error: File Write")
        return
    }
    msg := fmt.Sprintf("%d bytes written to disk", bytesWritten)
    appData.Dirty = 0
    appData.EditorSetStatusMessage(msg)
}

func editorRowsToString(appData *data.EditorConfig) string {
    fileAsString := ""
    for _, row := range appData.Row{
        fileAsString += *row.Runes+string('\n')
    }
    return fileAsString
}

