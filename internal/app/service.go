package app

import (
	"log/slog"
)

type Services struct {
	//LLMModel      *llmmodel.Service
}

//nolint:revive
func NewServices(_ *slog.Logger, _ *Repository) *Services {
	return &Services{
		//LLMModel:      llmmodel.NewService(log, repository.LLMModel),
	}
}
