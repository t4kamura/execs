package aws

import "testing"

func TestExtractNameFromARN(t *testing.T) {
	testCases := []struct {
		inputARN      string
		expectedName  string
		expectedError bool
	}{
		{"arn:aws:s3:::my-bucket", "my-bucket", false},
		{"arn:aws:lambda:us-west-2:123456789:function/my-function", "my-function", false},
		{"invalid-arn", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.inputARN, func(t *testing.T) {
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
