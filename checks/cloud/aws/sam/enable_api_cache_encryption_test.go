package sam

import (
	"testing"

	defsecTypes "github.com/aquasecurity/trivy/pkg/iac/types"

	"github.com/aquasecurity/trivy/pkg/iac/state"

	"github.com/aquasecurity/trivy/pkg/iac/providers/aws/sam"
	"github.com/aquasecurity/trivy/pkg/iac/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnableApiCacheEncryption(t *testing.T) {
	tests := []struct {
		name     string
		input    sam.SAM
		expected bool
	}{
		{
			name: "API unencrypted cache data",
			input: sam.SAM{
				APIs: []sam.API{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						RESTMethodSettings: sam.RESTMethodSettings{
							Metadata:           defsecTypes.NewTestMetadata(),
							CacheDataEncrypted: defsecTypes.Bool(false, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "API encrypted cache data",
			input: sam.SAM{
				APIs: []sam.API{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						RESTMethodSettings: sam.RESTMethodSettings{
							Metadata:           defsecTypes.NewTestMetadata(),
							CacheDataEncrypted: defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
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
			testState.AWS.SAM = test.input
			results := CheckEnableApiCacheEncryption.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableApiCacheEncryption.LongID() {
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
