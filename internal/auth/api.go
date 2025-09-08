package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no api key found")
	}

	apiKeyParts := strings.Split(apiKey, " ")
	if len(apiKeyParts) < 2 || apiKeyParts[0] != "ApiKey" {
		return "", fmt.Errorf("invalid api key")
	}

	return strings.TrimSpace(apiKeyParts[1]), nil
}