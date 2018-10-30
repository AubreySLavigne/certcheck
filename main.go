package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"certcheck/cert"

	"github.com/olekukonko/tablewriter"
)

const numWorkers = 4

func main() {

	resultChan := make(chan cert.Certificate)
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
		scanner := bufio.NewScanner(os.Stdin)
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

func certStatusWorker(inputChan <-chan string, resultChan chan<- cert.Certificate, wg *sync.WaitGroup) {
	for domain := range inputChan {
		resultChan <- cert.LoadCertificate(domain)
	}
	wg.Done()
}
