package aws

import "testing"

func TestExtractNameFromARN(t *testing.T) {
	testCases := []struct {
		testName      string
		inputARN      string
		expectedName  string
		expectedError bool
	}{
		{"ECS Cluster ARN", "arn:aws:ecs:us-west-2:123456789:cluster/cluster-name", "cluster-name", false},
		{"ECS Service ARN", "arn:aws:ecs:us-west-2:123456789:service/cluster-name/service-name", "service-name", false},
		{"Blank", "invalid-arn", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			resourceName, err := extractNameFromARN(tc.inputARN)

			if tc.expectedError && err == nil {
				t.Errorf("Expected an error, but no error occurred")
			} else if !tc.expectedError && err != nil {
				t.Errorf("An error occurred, but no error was expected: %v", err)
			}

			if resourceName != tc.expectedName {
				t.Errorf("Resource name does not match the expected value. Expected: %s, Actual: %s", tc.expectedName, resourceName)
			}
		})
	}
}
