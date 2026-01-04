package domain

import "encoding/json"

// JSONRaw помогает безопасно конвертировать разные JSON-типы в json.RawMessage.
func JSONRaw(b []byte) json.RawMessage {
	if b == nil {
		return json.RawMessage([]byte(`{}`))
	}
	return json.RawMessage(b)
}
