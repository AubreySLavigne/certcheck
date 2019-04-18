package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"certcheck/certificate"
	"certcheck/input"

	"github.com/olekukonko/tablewriter"
)

func main() {

	// Command Line arguments
	filename := flag.String("filename", "", "The name of the file containing the target domain names, one per line.")
	numWorkers := flag.Int("num-routines", runtime.NumCPU(), "The number of routines that will process this data concurrently.")
	flag.Parse()

	// Main body
	inputChan := make(chan string)
	resultChan := make(chan certificate.Certificate)

	wg := &sync.WaitGroup{}

	// Start all workers
	wg.Add(*numWorkers)
	for i := 0; i < *numWorkers; i++ {
		go worker(inputChan, resultChan, wg)
	}

	// When all workers are done, close the result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Send domains to the workers
	go processInput(inputChan, *filename)

	// Capture results
	results := make([]certificate.Certificate, 0)
	for res := range resultChan {
		results = append(results, res)
	}

	// Write results to output
	writeResultsTable(results)
}

func writeResultsTable(results []certificate.Certificate) {
	// Handle output
	statusTable := tablewriter.NewWriter(os.Stdout)
	statusTable.SetHeader([]string{"Name", "Status", "Details"})

	errorTable := tablewriter.NewWriter(os.Stderr)
	errorTable.SetHeader([]string{"Name", "Status", "Error"})

	for _, res := range results {
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

func processInput(inputChan chan<- string, filename string) {
	defer close(inputChan)

	var err error
	if filename != "" {
		err = input.LoadFile(inputChan, filename)
	} else {
		err = input.LoadFromPipe(inputChan)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func worker(inputChan <-chan string, resultChan chan<- certificate.Certificate, wg *sync.WaitGroup) {
	for domain := range inputChan {
		resultChan <- certificate.Load(domain)
	}
	wg.Done()
}
