package agent

import "nyxze/fayth/model"

type Agent struct {
	name  string
	model model.Model // The underlying model used for performing task

}

func NewAgent(name string, model model.Model) *Agent {
	return &Agent{
		name:  name,
		model: model,
	}
}
