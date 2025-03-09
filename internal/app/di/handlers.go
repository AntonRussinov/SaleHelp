package di

import (
	"errors"
)

type Handlers struct {
	// LLMModel *llmmodel.Handlers
}

var ErrInitHandlerss = errors.New("failed to initialize Handlerss")

func (prov *Provider) InitHandlers() error {
	// authHandlers, err := auth.NewHandlers(prov.cfg.Auth, prov.Services.User)
	// if err != nil {
	// 	return fmt.Errorf("%w: %w", ErrInitHandlerss, err)
	// }

	// prov.Handlerss.Auth = authHandlers

	// prov.Handlerss.UserQuota = quota.NewHandlers(prov.Services.Quota, prov.Presenters.Quota)
	return nil
}
