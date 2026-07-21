package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestOcrEngine_Value(t *testing.T) {
	valid := []OcrEngine{TESSERACT, EASY_OCR, TESSERACT_NEW, EASY_OCR_NEW}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted and normalized to "".
	assertValuerValid(t, "empty", OcrEngine(""), "")
}

func TestOcrEngine_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", OcrEngine("bogus"))
}

func TestOcrEngine_Scan(t *testing.T) {
	var engine OcrEngine
	err := engine.Scan("tesseract")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if engine != TESSERACT {
		utils.PrintTestError(t, engine, TESSERACT)
	}
}
