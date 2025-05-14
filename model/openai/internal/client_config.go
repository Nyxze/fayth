package internal

type ClientOptions func(*Client)

const (
	defaultBaseURL = "https://api.openai.com/v1"
)

func WithOrganisation(org string) ClientOptions {
	return func(c *Client) {
		c.organization = org
	}
}
func WithProject(project string) ClientOptions {
	return func(c *Client) {
		c.projectId = project
	}
}
