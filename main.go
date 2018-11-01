package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"certcheck/certificate"

	"github.com/olekukonko/tablewriter"
)

const numWorkers = 4

func main() {

	// Command Line arguments
	var filename string
	flag.StringVar(&filename, "filename", "", "The name of the file containing the target domain names, one per line.")
	flag.Parse()

	// Main body
	resultChan := make(chan certificate.Certificate)
	inputChan := make(chan string)

	wg := &sync.WaitGroup{}

	// Start all workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go certStatusWorker(inputChan, resultChan, wg)
	}

	// When all workers are done, close the result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Send domains to the workers
	go func() {
		var infile io.Reader
		if filename != "" {
			// Open filestream
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			infile = bufio.NewReader(file)
		} else {
			input, err := os.Stdin.Stat()
			if err != nil {
				log.Fatal(err)
			}
			if (input.Mode() & os.ModeCharDevice) != 0 {
				log.Fatal("Terminal Input not supported")
			}
			infile = os.Stdin
		}

		scanner := bufio.NewScanner(infile)
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}
		close(inputChan)
	}()

	statusTable := tablewriter.NewWriter(os.Stdout)
	statusTable.SetHeader([]string{"Name", "Status", "Details"})

	errorTable := tablewriter.NewWriter(os.Stderr)
	errorTable.SetHeader([]string{"Name", "Status", "Error"})

	for res := range resultChan {
		if res.Error != nil {
			errorTable.Append([]string{
				res.Domain,
				res.Status,
				fmt.Sprintf("%s", res.Error),
			})
		} else if res.Status != "" {
			statusTable.Append([]string{
				res.Domain,
				res.Status,
				res.Details,
			})
		}
	}

	statusTable.Render()
	fmt.Println("")
	errorTable.Render()
}

func certStatusWorker(inputChan <-chan string, resultChan chan<- certificate.Certificate, wg *sync.WaitGroup) {
	for domain := range inputChan {
		resultChan <- certificate.Load(domain)
	}
	wg.Done()
}
