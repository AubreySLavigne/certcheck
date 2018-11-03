package input

import (
	"bufio"
	"os"
)

// LoadFile opens the provided file and sends each line to the channel.
func LoadFile(inputChan chan<- string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return consume(inputChan, bufio.NewReader(file))
}
