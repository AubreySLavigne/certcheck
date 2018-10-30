package cert

import (
	"fmt"
	"time"

	certLookup "github.com/genkiroid/cert"
)

// Certificate contains information about the TLS Certificate lookup
type Certificate struct {
	Domain string
	cert   certLookup.Cert

	Status  string
	Error   error
	Details string
}

// LoadCertificate loads the certificate for the provided domain
func LoadCertificate(domain string) Certificate {
	res := Certificate{
		Domain: domain,
		cert:   *certLookup.NewCert(domain),
	}

	if res.cert.Error != "" {
		res.Status = "SSL Lookup Error"
		res.Error = fmt.Errorf(res.cert.Error)
		return res
	}

	certRange, err := getCertValidRange(res.cert)
	if err != nil {
		res.Error = err
	} else if certRange.contains(time.Now()) {
		res.Status = "Valid"
	} else {
		res.Status = "Expired"
	}

	res.Details = certRange.End.String()

	return res
}

func getCertValidRange(cert certLookup.Cert) (dateRange, error) {

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

type dateRange struct {
	Start time.Time
	End   time.Time
}

func (d *dateRange) contains(t time.Time) bool {
	return d.Start.Before(t) && t.Before(d.End)
}
