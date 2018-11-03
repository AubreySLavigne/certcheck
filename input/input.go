package input

import (
	"bufio"
	"io"
)

// consume reads each entry from infile to the provided channel
func consume(inputChan chan<- string, infile io.Reader) error {
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		inputChan <- scanner.Text()
	}

	return nil
}
