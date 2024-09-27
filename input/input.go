package input

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func InputByte() {
	//Note: this does stop after reciving one char because terminal is in
	// Standard mode, we will have to deactive it later

	fmt.Println("Input:")
	var char rune
	_, err := fmt.Scanf("%c", &char)
	if err != nil {
		log.Fatal("Error getting input Char")
	}
}

func InputChar() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return char
}

func ReaderByte() byte {
	STDINFILE := os.Stdin
	var charValue byte
	var err error
	reader := bufio.NewReader(STDINFILE)
	// TODO: Refactor this to take an array of bytes
	testByte := make([]byte, 4)
	reader.Read(testByte)
	charValue, err = reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			// fmt.Println("END OF FILE")
		}
	}
	return charValue
}

// NEW Func that gets the input []byte
func ReaderBytes() []byte {
	STDINFILE := os.Stdin
	reader := bufio.NewReader(STDINFILE)
	inputBytes := make([]byte, 3)
	reader.Read(inputBytes)
	return inputBytes
}
