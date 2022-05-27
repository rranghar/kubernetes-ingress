package validation

import (
	"fmt"
	"strings"

	v1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/externaldns/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// verifyDNSRecordType checks if provided record is a valid DNS record type.
// Valid records match the list of records implemented by the external-dns project.
func verifyDNSRecordType(record string) error {
	validRecords := []string{"A", "CNAME", "TXT", "SRV", "NS", "PTR"}
	records := make(map[string]bool, len(validRecords))
	for _, r := range validRecords {
		records[r] = true
	}
	_, ok := records[record]
	if !ok {
		return &field.Error{
			Type:     field.ErrorTypeNotSupported,
			Field:    "RecordType",
			BadValue: record,
			Detail:   fmt.Sprintf("supported values: %s", strings.Join(validRecords, ", ")),
		}
	}
	return nil
}

// verifyDNSName checks if provided string represents a valid DNS name.
func verifyDNSName(s string) error {
	result := validation.IsDNS1123Subdomain(s)
	if len(result) == 0 {
		return nil
	}
	return &field.Error{
		Type:     field.ErrorTypeInvalid,
		Field:    "DNSName",
		BadValue: s,
		Detail:   strings.Join(result, ", "),
	}
}

// vaerifyTargets checks if targets represent valid IP adresses.
// It returns an error if any of the provided targets is not an IP address.
func verifyTargets(targets v1.Targets) error {
	for _, target := range targets {
		result := validation.IsValidIP(target)
		if len(result) == 0 {
			continue
		}
		return &field.Error{
			Type:     field.ErrorTypeInvalid,
			Field:    "Targets",
			BadValue: target,
			Detail:   result[0],
		}
	}
	return nil
}

// verifyEndpoint checks if all Endpoint fields are valid.
func verifyEndpoint(e *v1.Endpoint) error {
	if err := verifyTargets(e.Targets); err != nil {
		return err
	}
	if err := verifyDNSRecordType(e.RecordType); err != nil {
		return err
	}
	return nil
}

// verif
func verifyDNSEndpointSpec(es *v1.DNSEndpointSpec) error {
	if len(es.Endpoints) == 0 {
		return &field.Error{
			Type:  field.ErrorTypeRequired,
			Field: "Endpoints",
		}
		//return errors.New("endpoints not provided")
	}
	return nil
}

// ValidateDNSEnpoints takes dnsendpoint and validates its fiels.
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
