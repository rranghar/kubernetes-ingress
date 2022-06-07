package validation

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	v1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/externaldns/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// ValidateDNSEnpoint validates if all DNSEndpoint fields are valid.
func ValidateDNSEndpoint(dnsendpoint *v1.DNSEndpoint) error {
	if err := validateDNSEndpointSpec(&dnsendpoint.Spec); err != nil {
		return err
	}
	return nil
}

// validateDNSEndpointSpec checks if endpoints are provided.
func validateDNSEndpointSpec(es *v1.DNSEndpointSpec) error {
	if len(es.Endpoints) == 0 {
		return fmt.Errorf("%w: no endpoints supplied, expected a list of endpoints", ErrTypeRequired)
	}
	for _, endpoint := range es.Endpoints {
		if err := validateEndpoint(endpoint); err != nil {
			return err
		}
	}
	return nil
}

// validateEndpoint checks if all Endpoint fields are valid.
func validateEndpoint(e *v1.Endpoint) error {
	if err := validateDNSName(e.DNSName); err != nil {
		return err
	}
	if err := validateTargets(e.Targets); err != nil {
		return err
	}
	if err := validateDNSRecordType(e.RecordType); err != nil {
		return err
	}
	if err := validateTTL(e.RecordTTL); err != nil {
		return err
	}
	return nil
}

// validateDNSName checks if provided string represents a valid DNS name.
func validateDNSName(name string) error {
	if issues := validation.IsDNS1123Subdomain(name); len(issues) > 0 {
		return fmt.Errorf("%w: name %s, %s", ErrTypeInvalid, name, strings.Join(issues, ", "))
	}
	return nil
}

// validateTargets checks if targets represent valid FQDN entries.
// It returns an error if any of the provided targets is not an IP address.
func validateTargets(targets v1.Targets) error {
	for _, target := range targets {
		if err := isFullyQualifiedDomainName(target); err != nil {
			return fmt.Errorf("%w: target %q is invalid, it should be a valid IP address or hostname", ErrTypeInvalid, target)
		}
	}
	return isUnique(targets)
}

// isUnique checks if targets are not duplicated.
// It returns error if targets has a duplicated entry.
func isUnique(targets v1.Targets) error {
	occured := make(map[string]bool)
	for _, target := range targets {
		if occured[target] {
			return fmt.Errorf("%w: target %s, expected unique targets", ErrTypeDuplicated, target)
		}
		occured[target] = true
	}
	return nil
}

// validateDNSRecordType checks if provided record is a valid DNS record type.
// Valid records match the list of records implemented by the external-dns project.
func validateDNSRecordType(record string) error {
	if !slices.Contains(validRecords, record) {
		return fmt.Errorf("%w: record %s, %s", ErrTypeNotSupported, record, strings.Join(validRecords, ","))
	}
	return nil
}

// validateTTL checks if TTL value is > 0.
func validateTTL(ttl v1.TTL) error {
	if ttl <= 0 {
		return fmt.Errorf("%w: ttl %d, ttl value should be > 0", ErrTypeNotInRange, ttl)
	}
	return nil
}

// isFullyQualifiedDomainName checks if the domain name is fully qualified.
// It requires a minimum of 2 segments and accepts a trailing . as valid.
func isFullyQualifiedDomainName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("%w: name not provided", ErrTypeInvalid)
	}
	if strings.HasSuffix(name, ".") {
		name = name[:len(name)-1]
	}
	if issues := validation.IsDNS1123Subdomain(name); len(issues) > 0 {
		return fmt.Errorf("%w: name %s is not valid subdomain, %s", ErrTypeInvalid, name, strings.Join(issues, ", "))
	}
	if len(strings.Split(name, ".")) < 2 {
		return fmt.Errorf("%w: name %s should be a domain with at least two segments separated by dots", ErrTypeInvalid, name)
	}
	for _, label := range strings.Split(name, ".") {
		if issues := validation.IsDNS1123Label(label); len(issues) > 0 {
			return fmt.Errorf("%w: label %s should conform to the definition of label in DNS (RFC1123), %s", ErrTypeInvalid, label, strings.Join(issues, ", "))
		}
	}
	return nil
}

var (
	// validRecords represents allowed DNS record names
	//
	// NGINX ingress controller at the moment supports
	// a subset of DNS record types listed in the external-dns project.
	validRecords = []string{"A", "CNAME"}

	// validation error types based on k8s validators
	ErrTypeNotSupported = errors.New("type not supported")
	ErrTypeInvalid      = errors.New("type invalid")
	ErrTypeDuplicated   = errors.New("type duplicated")
	ErrTypeRequired     = errors.New("type required")
	ErrTypeNotInRange   = errors.New("type not in range")
)
