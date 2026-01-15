package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
	"steamedeo.dev/gitlingo/internal/github"
)

func NewLanguageCommand() *cli.Command {
	return &cli.Command{
		Name:  "languages",
		Usage: "show your top programming languages statistics",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "all",
				Usage:    "show all languages instead of top 10",
				Required: false,
			},
		},
		Action: showLanguages,
	}
}

func showLanguages(ctx context.Context, cmd *cli.Command) error {
	showAll := cmd.Bool("all")
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Fetching GitHub data..."
	s.Color("magenta", "bold")
	s.Start()
	stats, err := github.ProcessGithubData(ctx)
	if err != nil {
		s.Stop()
		return err
	}
	s.Stop()
	fmt.Println() // Add spacing
	renderStats(stats, showAll)
	return nil
}

func renderStats(stats *github.GithubStats, showAll bool) {
	// Primary color: bright pink (205)
	// Secondary color: purple/lavender (141)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	userLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("223"))

	userValueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("223"))

	// Top 3 styles (primary pink)
	top3NumberStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	top3LangStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	top3BarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("218"))

	top3SizeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	// Rest styles (white)
	restNumberStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	restLangStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	restBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	restSizeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	// Header with border
	header := titleStyle.Render("✨ GitHub Language Stats ✨")
	fmt.Println(header)
	userLine := fmt.Sprintf("   %s%s",
		userLabelStyle.Render("@"),
		userValueStyle.Render(stats.Username))
	fmt.Println(userLine)
	fmt.Println()

	// Top 10 languages
	topLanguages := stats.LanguageStats
	if len(topLanguages) > 10 && !showAll {
		topLanguages = topLanguages[:10]
	}

	if len(topLanguages) == 0 {
		return
	}

	maxBytes := topLanguages[0].Bytes

	for i, lang := range topLanguages {
		// Bar width (max 35 chars for better proportions)
		barWidth := int(float64(lang.Bytes) / float64(maxBytes) * 35)
		if barWidth < 1 && lang.Bytes > 0 {
			barWidth = 1
		}
		bar := strings.Repeat("█", barWidth)

		sizeStr := formatBytes(lang.Bytes)

		// Calculate percentage
		percentage := float64(lang.Bytes) / float64(stats.TotalBytes) * 100
		sizeWithPercent := fmt.Sprintf("%s (%.1f%%)", sizeStr, percentage)

		// Use different colors for top 3 vs rest
		var numberStyle, langStyle, barStyle, sizeStyle lipgloss.Style
		if i < 3 {
			numberStyle = top3NumberStyle
			langStyle = top3LangStyle
			barStyle = top3BarStyle
			sizeStyle = top3SizeStyle
		} else {
			numberStyle = restNumberStyle
			langStyle = restLangStyle
			barStyle = restBarStyle
			sizeStyle = restSizeStyle
		}

		fmt.Printf("%s %s %s %s\n",
			numberStyle.Render(fmt.Sprintf("%2d.", i+1)),
			langStyle.Render(fmt.Sprintf("%-15s", lang.Name)),
			barStyle.Render(bar),
			sizeStyle.Render(sizeWithPercent))
	}

	fmt.Println()

	// Display total
	totalStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("141"))

	totalBytes := formatBytes(stats.TotalBytes)
	fmt.Printf("%s %s\n", totalStyle.Render("Total Code Scanned:"), totalStyle.Render(totalBytes))
	fmt.Println()
}

func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
