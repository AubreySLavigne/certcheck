package input

import (
	"strings"
	"testing"
)

func TestConsume(t *testing.T) {

	// Create an io.Reader that will return the following strings
	data := []string{"First", "Second", "Third"}
	input := strings.Join(data, "\n")
	ioReader := strings.NewReader(input)
	inputChan := make(chan string)

	// Process all strings in a separate goroutine
	go func() {
		err := consume(inputChan, ioReader)
		if err != nil {
			t.Errorf("Consume expects return of nil")
		}
		close(inputChan)
	}()

	// Check that the data is received as expected and in order
	i := 0
	for res := range inputChan {
		if expected := data[i]; res != expected {
			t.Errorf("Unexpected Input: Expected %v, got %v", expected, res)
		}

		i++
	}

	// Check that all strings were processed
	if count := len(data); i != count {
		t.Errorf("Did not process all domains. Expected %d, got %d", count, i)
	}
}
