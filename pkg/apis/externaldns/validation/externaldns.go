package validation

import (
	"errors"
	"fmt"

	v1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/externaldns/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// verifyDNSRecordType checks if provided record is a valid DNS record type.
func verifyDNSRecordType(record string) error {
	validRecords := map[string]bool{
		"A":     true,
		"CNAME": true,
		"TXT":   true,
		"SRV":   true,
		"NS":    true,
		"PTR":   true,
	}
	_, ok := validRecords[record]
	if !ok {
		return fmt.Errorf("invalid DNS record: %s", record)
	}
	return nil
}

// verifyDNSName checks if provided string represents a valid DNS name.
func verifyDNSName(s string) error {
	return nil
}

// vaerifyTargets checks if targets represent valid IP adresses.
// It returns an error if any of the provided targets is not an IP address.
func verifyTargets(targets v1.Targets) error {
	for _, target := range targets {
		result := validation.IsValidIP(target)
		if len(result) == 0 {
			continue
		}
		return errors.New(result[0])
	}
	return nil
}

// verifyEndpoint
func verifyEndpoint(e *v1.Endpoint) error {
	if err := verifyTargets(e.Targets); err != nil {
		return err
	}
	if err := verifyDNSRecordType(e.RecordType); err != nil {
		return err
	}

	return nil
}

func verifyDNSEndpointSpec(es *v1.DNSEndpointSpec) error {
	if len(es.Endpoints) == 0 {
		return errors.New("endpoints not provided")
	}
	return nil
}

func ValidateDNSEndpoint(dnsendpoint *v1.DNSEndpoint) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error validating DNSEndpoint: %w", err)
		}
	}()

	if err := verifyDNSEndpointSpec(&dnsendpoint.Spec); err != nil {
		return err
	}

	return nil
}
