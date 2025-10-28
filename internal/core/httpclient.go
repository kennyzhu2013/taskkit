package core

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

// NewHTTPClient provides a preconfigured Resty client with sane defaults
// - 30s timeout
// - optional InsecureSkipVerify for trusted environments
func NewHTTPClient(skipTLS bool) *resty.Client {
	c := resty.New()
	c.SetTimeout(30 * time.Second)
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: skipTLS})
	return c
}

// GetGitHubAuthHeaders returns standard headers for GitHub API requests
// Includes Authorization only if a token is provided
func GetGitHubAuthHeaders(cliToken string) map[string]string {
	token := ResolveGitHubToken(cliToken)
	headers := map[string]string{
		"Accept": "application/vnd.github+json",
	}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}
	return headers
}
