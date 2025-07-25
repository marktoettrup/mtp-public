package flags

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	PWDBEndpoint = "endpoint"
	PWDBAPIKey   = "api-key"
)

func GetPWDBEndpoint(cmd *cobra.Command) (string, error) {
	pwdbEndpoint, err := cmd.Flags().GetString(PWDBEndpoint)
	if err != nil {
		return "", fmt.Errorf("could not get pwdb endpoint flag: %w", err)
	}
	if pwdbEndpoint == "" {
		return "", errors.New("pwdb endpoint is required")
	}
	return pwdbEndpoint, nil
}

func GetPWDBAPIKey(cmd *cobra.Command) (string, error) {
	apiKey := os.Getenv("IKS_PWDB_APIKEY")
	if apiKey == "" {
		apiKey, err := cmd.Flags().GetString(PWDBAPIKey)
		if err != nil {
			return "", fmt.Errorf("could not get pwdb api-key flag: %w", err)
		}
		if apiKey == "" {
			return "", errors.New("pwdb api-key is required")
		}
	}
	return apiKey, nil
}
