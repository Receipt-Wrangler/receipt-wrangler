package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestAiClientType_Value(t *testing.T) {
	valid := []AiClientType{
		OPEN_AI_CUSTOM,
		OPEN_AI,
		GEMINI,
		OPEN_AI_CUSTOM_NEW,
		OPEN_AI_NEW,
		GEMINI_NEW,
		OLLAMA,
	}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted and normalized to "".
	assertValuerValid(t, "empty", AiClientType(""), "")
}

func TestAiClientType_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", AiClientType("bogus"))
}

func TestAiClientType_Scan(t *testing.T) {
	var clientType AiClientType
	err := clientType.Scan("openAi")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if clientType != OPEN_AI {
		utils.PrintTestError(t, clientType, OPEN_AI)
	}
}
