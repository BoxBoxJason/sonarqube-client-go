package helpers

import (
	"context"
	"strings"

	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// GetEdition returns the SonarQube edition reported by the connected server
// (e.g. "community", "developer", "enterprise", "datacenter"). It relies on
// the System Info endpoint, which does not require an active license.
func GetEdition(client *sonar.Client) (string, error) {
	info, _, err := client.System.Info(context.Background())
	if err != nil {
		return "", err
	}

	return info.System.Edition, nil
}

// IsEnterpriseOrAbove reports whether the given edition string is Enterprise
// or Data Center edition (Data Center being a superset of Enterprise).
func IsEnterpriseOrAbove(edition string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(edition, " ", ""))

	return strings.Contains(normalized, "enterprise") || strings.Contains(normalized, "datacenter")
}

// HasActiveLicense reports whether the connected server currently has a
// valid, supported Enterprise Edition license installed. It never returns
// an error: any failure to determine license state is treated as "no
// license", since license activation is optional for most enterprise-only
// endpoints covered by this suite.
func HasActiveLicense(client *sonar.Client) bool {
	result, _, err := client.V2.Entitlements.GetLicense(context.Background())
	if err != nil || result == nil {
		return false
	}

	return result.Supported && !result.Expired
}
