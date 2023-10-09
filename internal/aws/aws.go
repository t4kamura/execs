package aws

import (
	"errors"
	"regexp"
)

func extractNameFromARNs(arns []string) ([]string, error) {
	var names []string
	for _, arn := range arns {
		name, err := extractNameFromARN(arn)
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}
	return names, nil
}

func extractNameFromARN(arn string) (string, error) {
	arnPattern := `^arn:aws:[^:]+:[^:]+:[^:]+:[^:]+/([^/]+)(?:/([^/]+))?`

	re := regexp.MustCompile(arnPattern)
	match := re.FindStringSubmatch(arn)

	if match == nil {
		return "", errors.New("invalid ARN")
	}

	return match[1], nil
}
