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
	if err := godotenv.Load(); err != nil {
		return err
	}
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return errors.New("github token not set!")
	}

	AppConfig = &Config{
		GithubToken: githubToken,
	}
	return nil
}
