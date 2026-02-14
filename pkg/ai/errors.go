package ai

import "errors"

var (
	// ErrNoQuota — лимит бесплатных и платных ответов исчерпан.
	ErrNoQuota = errors.New("ai: no reply quota")
	// ErrAIDisabled — AI-ответы для магазина отключены (shop_ai_settings.is_active = false).
	ErrAIDisabled = errors.New("ai: disabled for shop")
	// ErrInvalidTone — неизвестный tone_of_voice.
	ErrInvalidTone = errors.New("ai: invalid tone of voice")
)
