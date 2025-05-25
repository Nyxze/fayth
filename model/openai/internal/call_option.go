package internal

import (
	"fmt"
	"net/url"
	"strings"
)

// CallOption represents a functional option that modifies the behavior of an OpenAI API call.
//
// These options can be applied at different levels in the call hierarchy:
// - Client level: sets global defaults.
// - Service level: overrides client defaults for specific services.
// - Method level: overrides both client and service defaults for a specific call.
//
// The final configuration is resolved in the following precedence order (from lowest to highest):
//
//	Client -> Service -> Method
//
// This allows for flexible and composable customization of API behavior.
type CallOption func(*CallConfig) error

type CallConfig struct {
	BaseUrl      *url.URL
	Organization string
	Project      string
	APIKey       string
}

const (
	defaultBaseURL = "https://api.openai.com/v1"
)

func WithBaseURL(base string) CallOption {
	return func(cc *CallConfig) error {
		u, err := url.Parse(base)
		if err != nil {
			return fmt.Errorf("call options: WithBaseURL failed to parse url %s\n", err)
		}
		if u.Path != "" && !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		cc.BaseUrl = u
		return nil
	}
}
func WithOrganization(org string) CallOption {
	return func(cc *CallConfig) error {
		cc.Organization = org
		return nil
	}
}

func WithProject(project string) CallOption {
	return func(cc *CallConfig) error {
		cc.Project = project
		return nil
	}
}
func WithAPIKey(key string) CallOption {
	return func(cc *CallConfig) error {
		cc.APIKey = key
		return nil
	}
}
