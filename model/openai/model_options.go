package openai

type options struct {
	ApiKey string
	Model  string
}

type ModelOptions func(*options)

func WithModel(modelName string) ModelOptions {
	return func(c *options) {
		c.Model = modelName
	}
}

func WithAPIKey(key string) ModelOptions {
	return func(c *options) {
		c.ApiKey = key
	}
}
