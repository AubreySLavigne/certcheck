package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	certLookup "github.com/genkiroid/cert"
	"github.com/olekukonko/tablewriter"
)

// result is the returned value of the certificate lookup
type result struct {
	domain string
	status string
	err    error
}

const numWorkers = 4

func main() {
	domains := []string{
		"neverssl.com",
		"www.httpvshttps.com",
		"google.com",
		"reddit.com",
	}

	resultChan := make(chan result)
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
		for _, domain := range domains {
			inputChan <- domain
		}
		close(inputChan)
	}()

	statusTable := tablewriter.NewWriter(os.Stdout)
	statusTable.SetHeader([]string{"Name", "Status"})

	errorTable := tablewriter.NewWriter(os.Stdout)
	errorTable.SetHeader([]string{"Name", "Error"})

	for res := range resultChan {
		if res.err != nil {
			errorTable.Append([]string{
				res.domain,
				fmt.Sprintf("%s", res.err),
			})
		}
		if res.status != "" {
			statusTable.Append([]string{
				res.domain,
				res.status,
			})
		}
	}

	statusTable.Render()
	fmt.Println("")
	errorTable.Render()
}

func certStatusWorker(inputChan <-chan string, resultChan chan<- result, wg *sync.WaitGroup) {
	for domain := range inputChan {
		resultChan <- getCertStatus(domain)
	}
	wg.Done()
}

func getCertStatus(domain string) result {
	certRange, err := getCertValidRange(domain)
	res := result{
		domain: domain,
	}

	if err != nil {
		res.err = err
	} else if certRange.contains(time.Now()) {
		res.status = fmt.Sprintf("Valid   - %s", certRange.End)
	} else {
		res.status = fmt.Sprintf("Expired - %s", certRange.End)
	}

	return res
}

type dateRange struct {
	Start time.Time
	End   time.Time
}

func (d *dateRange) contains(t time.Time) bool {
	return d.Start.Before(t) && t.Before(d.End)
}

func getCertValidRange(domain string) (dateRange, error) {
	cert := certLookup.NewCert(domain)
	if cert.Error != "" {
		return dateRange{}, fmt.Errorf("SSL Lookup Error       - %s", cert.Error)
	}

	startTime, err := time.Parse("2006-01-02 15:04:05 Z0700 MST", cert.NotBefore)
	if err != nil {
		return dateRange{}, fmt.Errorf("Start time parse error - %s", err)
	}

	endTime, err := time.Parse("2006-01-02 15:04:05 Z0700 MST", cert.NotAfter)
	if err != nil {
		return dateRange{}, fmt.Errorf("End time parse error   - %s", err)
	}

	return dateRange{
		Start: startTime,
		End:   endTime,
	}, nil
}
