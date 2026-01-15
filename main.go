package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v3"
	"steamedeo.dev/gitlingo/cmd"
	"steamedeo.dev/gitlingo/internal/config"
)

// Version information (injected at build time via ldflags)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	app := &cli.Command{
		Name:    "gitlingo",
		Usage:   "Fetch and visualize statistics about your github repos",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version-info",
				Aliases: []string{"V"},
				Usage:   "Show detailed version information",
				Action: func(ctx context.Context, cmd *cli.Command, b bool) error {
					if b {
						fmt.Printf("Version:    %s\n", version)
						fmt.Printf("Git Commit: %s\n", commit)
						fmt.Printf("Build Date: %s\n", date)
						os.Exit(0)
					}
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			cmd.NewGeneratecommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
