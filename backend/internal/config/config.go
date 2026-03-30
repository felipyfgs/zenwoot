package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	App   AppConfig
	DB    DBConfig
	NATS  NATSConfig
	JWT   JWTConfig
	MinIO MinIOConfig
}

type AppConfig struct {
	Host string `env:"APP_HOST" env-default:"0.0.0.0"`
	Port string `env:"APP_PORT" env-default:"3000"`
	Env  string `env:"APP_ENV"  env-default:"development"`
}

type DBConfig struct {
	DSN string `env:"DATABASE_URL" env-required:"true"`
}

type NATSConfig struct {
	URL string `env:"NATS_URL" env-default:"nats://localhost:4222"`
}

type JWTConfig struct {
	Secret      string `env:"JWT_SECRET"       env-required:"true"`
	ExpiryHours int    `env:"JWT_EXPIRY_HOURS" env-default:"24"`
}

type MinIOConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT"  env-default:"localhost:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" env-default:"zenwoot-admin"`
	SecretKey string `env:"MINIO_SECRET_KEY" env-default:"zenwoot-secret123"`
	Bucket    string `env:"MINIO_BUCKET"     env-default:"zenwoot"`
	UseSSL    bool   `env:"MINIO_USE_SSL"    env-default:"false"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
