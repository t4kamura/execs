package aws

import (
	"strings"

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
			// profile section name is like "[profile test]"
			// default profile section name is like "[default]"
			profileWords := strings.Split(section.Name(), " ")

			if len(profileWords) == 2 {
				profiles = append(profiles, profileWords[1])
				continue
			} else if len(profileWords) == 1 && profileWords[0] == "default" {
				profiles = append(profiles, profileWords[0])
			}
		}
	}

	return profiles, nil
}
