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

func (f fakeModel) Generate(ctx context.Context, m []model.Message) (*model.Generation, error) {
	if len(m) == 0 {
		return nil, errors.New("empty message provided")
	}
	if f.Response.Contents == nil {
		return nil, errors.New("no content set in response")
	}
	return &model.Generation{
		Results: []*model.Message{&f.Response},
	}, nil
}
