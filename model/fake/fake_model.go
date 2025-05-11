package fake

import (
	"context"
	"errors"
	"nyxze/fayth/model"
)

type fakeModel struct {
	Name     string
	Response model.Message
}

// Create a fake model that respond with the given response
func NewModel(name string, resp model.Message) model.Model {
	return &fakeModel{
		Name:     name,
		Response: resp,
	}
}

func (f fakeModel) Generate(ctx context.Context, m []model.Message) (model.Message, error) {
	if len(m) == 0 {
		return model.Message{}, errors.New("empty message provided")
	}
	if f.Response.Contents == nil {
		return model.Message{}, errors.New("no content set in response")
	}
	return f.Response, nil
}
