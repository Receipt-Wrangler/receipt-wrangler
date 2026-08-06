package structs

import "receipt-wrangler/api/internal/models"

type EmailSettings struct {
	Host     string
	Port     int
	Username string
	Password string
}

type AiSettings struct {
	AiType     models.AiClientType `json:"type"`
	Key        string              `json:"key"`
	Url        string              `json:"url"`
	Model      string              `json:"model"`
	NumWorkers int                 `json:"numWorkers"`
	OcrEngine  models.OcrEngine    `json:"ocrEngine"`
}

type FeatureConfig struct {
	EnableLocalSignUp bool `json:"enableLocalSignUp"`
	AiPoweredReceipts bool `json:"aiPoweredReceipts"`
	// LoginQrUrl is the composed App/Universal Link the desktop login page
	// encodes as a QR. Empty unless an admin enabled the login QR and set a
	// mobile server URL. This is the only login-QR value exposed pre-auth.
	LoginQrUrl string `json:"loginQrUrl"`
}

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Engine   string `json:"engine"`
	Filename string `json:"filename"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DebugConfig struct {
	DebugOcr bool `json:"debugOcr"`
}

type Config struct {
	SecretKey            string          `json:"secretKey"`
	AiSettings           AiSettings      `json:"aiSettings"`
	EmailPollingInterval int             `json:"emailPollingInterval"`
	EmailSettings        []EmailSettings `json:"emailSettings"`
	Features             FeatureConfig   `json:"features"`
	Database             DatabaseConfig  `json:"database"`
	Debug                DebugConfig     `json:"debug"`
}
