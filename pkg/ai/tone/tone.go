package tone

import "salehelp/domain"

// Instruction возвращает строгую инструкцию по тону для промпта.
// Допустимые значения: standard | friendly | official | empathic | minimalist.
func Instruction(toneOfVoice string) string {
	s, ok := instructions[toneOfVoice]
	if !ok {
		s = instructions[domain.ToneStandard]
	}
	return s
}

// instructions — централизованная конфигурация tone of voice.
var instructions = map[string]string{
	domain.ToneStandard:  "Вежливо, нейтрально, как менеджер бренда. Без эмодзи.",
	domain.ToneFriendly: "Тёплый, заботливый тон. Допускается 1–2 уместных эмодзи.",
	domain.ToneOfficial: "Строго и коротко. Акцент на качестве и стандартах. Без эмодзи.",
	domain.ToneEmpathic: "Акцент на чувствах клиента, мягко. Особенно в негативных отзывах — признание и поддержка.",
	domain.ToneMinimalist: "Очень коротко и по делу. Минимум слов.",
}
