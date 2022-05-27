package validation

import (
	"errors"
	"testing"

	v1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/externaldns/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestVerifyDNSRecord_ErrorsOnInvalidRecordType(t *testing.T) {
	t.Parallel()
	err := verifyDNSRecordType("B")
	if err == nil {
		t.Fatal(err)
	}
	if err != nil {
		var fieldErr *field.Error
		if !errors.As(err, &fieldErr) {
			t.Fatal(err)
		}
	}
}

func TestVerifyTargets_ErrorsOnInvalidIP(t *testing.T) {
	t.Parallel()
	invalidTargets := v1.Targets{"10.12.34.1111"}
	err := verifyTargets(invalidTargets)
	if err == nil {
		t.Fatal(err)
	}
	if err != nil {
		var fieldErr *field.Error
		if !errors.As(err, &fieldErr) {
			t.Fatal(err)
		}
	}
}

func TestVerifyDNSname_ErrorsOnInvalidName(t *testing.T) {
	t.Parallel()
	invalidName := "abc.example..."
	err := verifyDNSName(invalidName)
	if err == nil {
		t.Fatal(err)
	}
	if err != nil {
		var fieldErr *field.Error
		if !errors.As(err, &fieldErr) {
			t.Fatal(err)
		}
	}
}

func TestVerifyEndpoint(t *testing.T) {
	tt := []struct {
		name  string
		input v1.Endpoint
	}{
		{
			name: "Returns error on invalid endpoint targets",
			input: v1.Endpoint{
				DNSName:    "",
				Targets:    []string{"1000.1.1.1"},
				RecordType: "A",
				RecordTTL:  3600,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if err := verifyEndpoint(&tc.input); err == nil {
				t.Fatal(err)
			}
		})
	}
}

func TestValidateDNSEndpoint(t *testing.T) {
	t.Parallel()
	tt := []struct {
		endpoint *v1.DNSEndpoint
		name     string
	}{
		{
			name:     "Return error on empty DNSEndpoint struct",
			endpoint: &v1.DNSEndpoint{},
		},
		{
			name: "Return error on empty DNSEndpointSpec struct",
			endpoint: &v1.DNSEndpoint{
				Spec: v1.DNSEndpointSpec{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateDNSEndpoint(tc.endpoint); err == nil {
				t.Fatal(err)
			}
		})
	}
}
