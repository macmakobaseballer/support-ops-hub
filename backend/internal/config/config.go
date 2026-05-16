package config

import "github.com/kelseyhightower/envconfig"

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	DBDSN              string `envconfig:"DB_DSN" required:"true"`
	GatewayPort        string `envconfig:"GATEWAY_PORT" default:"8000"`
	JWTSecret          string `envconfig:"JWT_SECRET" required:"true"`
	JWTTTL             string `envconfig:"JWT_TTL" default:"15m"`
	CORSAllowedOrigins string `envconfig:"CORS_ALLOWED_ORIGINS" default:"http://localhost:3000"`
}

// Load reads environment variables into Config.
func Load() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return Config{}, err
	}
	return c, nil
}
