package internal

type ClientOptions func(*Client)

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
