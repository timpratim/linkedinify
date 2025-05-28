// internal/config/config.go
package config

import "os"

type Config struct {
	HTTPAddr      string
	DSN           string
	JWTSecret     []byte
	OpenAIToken   string
	TreblleToken  string
	TreblleAPIKey string
}

func Load() Config {
	return Config{
		HTTPAddr:      envDefault("HTTP_ADDR", ":8080"),
		DSN:           envDefault("DATABASE_DSN", "postgres:///linkedinify?sslmode=disable"),
		JWTSecret:     []byte(envDefault("JWT_SECRET", "supersecret")),
		OpenAIToken:   os.Getenv("OPENAI_TOKEN"),
		TreblleToken:  os.Getenv("TREBLLE_SDK_TOKEN"),
		TreblleAPIKey: os.Getenv("TREBLLE_API_KEY"),
	}
}

func envDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c Config) GetJWTSecret() []byte {
	return c.JWTSecret
}
