package fileio

import (
	"bufio"
	"os"
	"wicho/whim/app/data"
)

func EditorOpen(appData *data.EditorConfig, fileName string){
    file, err := os.Open(fileName)
    if err != nil {
        appData.Die()
    }
    defer file.Close()

    line := ""
    totalLines := 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        totalLines++
        line = scanner.Text()
        appData.EditorAppendRow(line)
    }
}
