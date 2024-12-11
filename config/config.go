package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    Server struct {
        Port string
    }
    Elasticsearch struct {
        URLs     []string
        Username string
        Password string
    }
}

func Load() *Config {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    config := &Config{}

    // Server configuration
    config.Server.Port = getEnv("SERVER_PORT", "8080")

    // Elasticsearch configuration
    config.Elasticsearch.URLs = []string{getEnv("ES_URL", "http://localhost:9200")}
    config.Elasticsearch.Username = getEnv("ES_USERNAME", "")
    config.Elasticsearch.Password = getEnv("ES_PASSWORD", "")

    return config
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
