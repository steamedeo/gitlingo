package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GithubToken string
}

var AppConfig *Config

func Load() error {
	// Try to load from environment variable first
	githubToken := os.Getenv("GITHUB_TOKEN")

	// If not found in environment, try loading from .env file
	if githubToken == "" {
		// Ignore error if .env file doesn't exist - environment variable takes precedence
		_ = godotenv.Load()
		githubToken = os.Getenv("GITHUB_TOKEN")
	}

	if githubToken == "" {
		return errors.New("github token not set! Set GITHUB_TOKEN environment variable or create a .env file with github_token=your_token")
	}

	AppConfig = &Config{
		GithubToken: githubToken,
	}
	return nil
}
