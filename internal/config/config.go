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
	godotenv.Load()
	githubToken := os.Getenv("github_token")
	if githubToken == "" {
		return errors.New("Github token not set!")
	}

	AppConfig = &Config{
		GithubToken: githubToken,
	}
	return nil
}
