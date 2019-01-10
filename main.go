package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"certcheck/certificate"
	"certcheck/input"

	"github.com/olekukonko/tablewriter"
)

func main() {

	// Command Line arguments
	filename := flag.String("filename", "", "The name of the file containing the target domain names, one per line.")
	numWorkers := flag.Int("num-routines", 4, "The number of routines that will process this data concurrently.")
	flag.Parse()

	// Main body
	inputChan := make(chan string)
	resultChan := make(chan result)

	wg := &sync.WaitGroup{}

	// Start all workers
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(inputChan, resultChan, wg)
	}

	// When all workers are done, close the result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Send domains to the workers
	go processInput(inputChan, *filename)

	write(resultChan)
}

func write(resultChan <-chan result) {
	statusTable := tablewriter.NewWriter(os.Stdout)
	statusTable.SetHeader([]string{"Name", "Status", "Details"})

	errorTable := tablewriter.NewWriter(os.Stderr)
	errorTable.SetHeader([]string{"Name", "Status", "Error"})

	for res := range resultChan {
		if res.cert.Error != nil {
			errorTable.Append([]string{
				res.cert.Domain,
				res.cert.Status,
				fmt.Sprintf("%s", res.cert.Error),
			})
		} else if res.cert.Status != "" {
			statusTable.Append([]string{
				res.cert.Domain,
				res.cert.Status,
				res.cert.Details,
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

type result struct {
	cert certificate.Certificate
	err  error
}

func worker(inputChan <-chan string, resultChan chan<- result, wg *sync.WaitGroup) {
	for domain := range inputChan {
		cert, err := certificate.Load(domain)
		resultChan <- result{
			cert: cert,
			err:  err,
		}
	}
	wg.Done()
}
