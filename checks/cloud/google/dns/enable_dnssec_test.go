package dns

import (
	"testing"

	defsecTypes "github.com/aquasecurity/trivy/pkg/iac/types"

	"github.com/aquasecurity/trivy/pkg/iac/state"

	"github.com/aquasecurity/trivy/pkg/iac/providers/google/dns"
	"github.com/aquasecurity/trivy/pkg/iac/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnableDnssec(t *testing.T) {
	tests := []struct {
		name     string
		input    dns.DNS
		expected bool
	}{
		{
			name: "DNSSec disabled and required when visibility explicitly public",
			input: dns.DNS{
				ManagedZones: []dns.ManagedZone{
					{
						Metadata:   defsecTypes.NewTestMetadata(),
						Visibility: defsecTypes.String("public", defsecTypes.NewTestMetadata()),
						DNSSec: dns.DNSSec{
							Metadata: defsecTypes.NewTestMetadata(),
							Enabled:  defsecTypes.Bool(false, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "DNSSec enabled",
			input: dns.DNS{
				ManagedZones: []dns.ManagedZone{
					{
						Metadata:   defsecTypes.NewTestMetadata(),
						Visibility: defsecTypes.String("public", defsecTypes.NewTestMetadata()),
						DNSSec: dns.DNSSec{
							Metadata: defsecTypes.NewTestMetadata(),
							Enabled:  defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "DNSSec not required when private",
			input: dns.DNS{
				ManagedZones: []dns.ManagedZone{
					{
						Metadata:   defsecTypes.NewTestMetadata(),
						Visibility: defsecTypes.String("private", defsecTypes.NewTestMetadata()),
						DNSSec: dns.DNSSec{
							Metadata: defsecTypes.NewTestMetadata(),
							Enabled:  defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.Google.DNS = test.input
			results := CheckEnableDnssec.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableDnssec.LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
