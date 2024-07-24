package app

import (
	"log/slog"
)

type Handler struct {
	// LLMModel *llmmodel.Handler
}

func NewHandler(_ *Config, _ *slog.Logger, _ *Services, _ *Presenter) (*Handler, error) {
	return &Handler{
		//LLMModel: llmmodel.NewHandler(services.LLMModel),
	}, nil
}
