package input

import (
	"fmt"
	"os"
)

// LoadFromPipe reads each line from stdin and sends them to the provided channel.
//
// This function will only read when data is piped to the program, but will
// not work if the data is manually entered (in the terminal)
func LoadFromPipe(inputChan chan<- string) error {
	input, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if (input.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("Terminal Input not supported")
	}

	return consume(inputChan, os.Stdin)
}
