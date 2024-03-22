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
        //fmt.Println(line)
        appData.ABuf.WriteString(line)
    }
    // scanner.Scan doesnt include the newlines ?!
    // line = strings.TrimSuffix(line, "\n")
    appData.Row.Size = totalLines
    appData.Row.Runes = &line

    newRow := data.EditorRow{
        Size: len(line),
        Runes: &line,
    }
    appData.Row = newRow
    appData.NumRows = 1
}
