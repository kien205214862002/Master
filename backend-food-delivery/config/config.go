package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	ProductionEnv = "production"
	TestEnv       = "testing"

	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 5 * time.Minute
)

type Schema struct {
	Environment   string `env:"environment"`
	Port          int    `env:"port"`
	AuthSecret    string `env:"auth_secret"`
	DatabaseURI   string `env:"database_uri"`
	RedisURI      string `env:"redis_uri"`
	RedisPassword string `env:"redis_password"`
	RedisDB       int    `env:"redis_db"`
	S3BucketName  string `env:"s3_bucket_name"`
	S3Region      string `env:"s3_region"`
	S3APIKey      string `env:"s3_api_key"`
	S3SecretKey   string `env:"s3_secret_key"`
	S3Domain      string `env:"s3_domain"`
	SystemSecret  string `env:"system_secret"`
}

var (
	cfg Schema
)

func init() {
	environment := os.Getenv("environment")
	err := godotenv.Load("config/config.yaml")
	if err != nil && environment != TestEnv {
		log.Fatalf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}
}

func GetConfig() *Schema {
	return &cfg
}
