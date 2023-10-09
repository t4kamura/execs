package aws

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

// GetProfiles returns AWS profile names.
// Get profile names from ~/.aws/config.
func GetProfiles() ([]string, error) {
	configFilename := config.DefaultSharedConfigFilename()
	f, err := ini.Load(configFilename)
	if err != nil {
		return nil, err
	}

	var profiles []string
	for _, section := range f.Sections() {
		if len(section.Keys()) != 0 {
			profiles = append(profiles, section.Name())
		}
	}

	return profiles, nil
}
