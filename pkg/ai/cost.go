package ai

import "math"

// Курс конвертации: 1 USD = 90 RUB = 9000 копеек (фиксированный для оценки).
const usdToKopecks = 9000

type modelPricing struct {
	InputPerMillion  int // копеек за 1M input-токенов
	OutputPerMillion int // копеек за 1M output-токенов
}

// knownModelPricing — таблица цен OpenAI в копейках за 1M токенов.
// Обновлять при смене моделей или тарифов OpenAI.
var knownModelPricing = map[string]modelPricing{
	// gpt-4o-mini: $0.15/$0.60 per 1M
	"gpt-4o-mini":              {InputPerMillion: 1350, OutputPerMillion: 5400},
	"gpt-4o-mini-2024-07-18":   {InputPerMillion: 1350, OutputPerMillion: 5400},
	"gpt-4o-mini-2024-11-05":   {InputPerMillion: 1350, OutputPerMillion: 5400},
	// gpt-4o: $2.50/$10.00 per 1M
	"gpt-4o":              {InputPerMillion: 22500, OutputPerMillion: 90000},
	"gpt-4o-2024-08-06":   {InputPerMillion: 22500, OutputPerMillion: 90000},
	"gpt-4o-2024-11-20":   {InputPerMillion: 22500, OutputPerMillion: 90000},
}

// EstimateCostKopeks возвращает оценочную стоимость вызова в копейках.
// Возвращает 0 для неизвестных моделей.
func EstimateCostKopeks(model string, inputTokens, outputTokens int) int {
	p, ok := knownModelPricing[model]
	if !ok {
		return 0
	}
	cost := float64(p.InputPerMillion)*float64(inputTokens)/1_000_000 +
		float64(p.OutputPerMillion)*float64(outputTokens)/1_000_000
	return int(math.Round(cost))
}
