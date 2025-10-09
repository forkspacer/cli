package validation

import (
	"fmt"
	"regexp"
)

var (
	// DNS-1123 subdomain regex (max 253 chars)
	dns1123SubdomainRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
)

// ValidateDNS1123Subdomain validates a string as a DNS-1123 subdomain
func ValidateDNS1123Subdomain(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) > 253 {
		return fmt.Errorf("name must be 253 characters or less (got %d)", len(name))
	}

	if !dns1123SubdomainRegex.MatchString(name) {
		return fmt.Errorf("name must be lowercase alphanumeric with '-' or '.' only")
	}

	return nil
}

// DNS1123Examples returns example valid names
func DNS1123Examples() []string {
	return []string{
		"dev-env",
		"staging.cluster",
		"test-ws-123",
		"prod",
	}
}

// DNS1123Requirements returns human-readable requirements
func DNS1123Requirements() []string {
	return []string{
		"Lowercase letters (a-z)",
		"Numbers (0-9)",
		"Hyphens (-) or dots (.)",
		"Max 253 characters",
		"Must start and end with alphanumeric character",
	}
}
