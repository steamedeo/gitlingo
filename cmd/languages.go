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

	// Header
	header := titleStyle.Render("✨ GitHub Language Stats ✨")
	fmt.Println(header)
	userLine := fmt.Sprintf("   %s%s",
		userLabelStyle.Render("@"),
		userValueStyle.Render(stats.Username))
	fmt.Println(userLine)
	fmt.Println()

	// Top 15 languages by default (unless --all is used)
	topLanguages := stats.LanguageStats
	if len(topLanguages) > 15 && !showAll {
		topLanguages = topLanguages[:15]
	}

	if len(topLanguages) == 0 {
		return
	}

	// Build the language list widget
	languageList := buildLanguageList(stats, topLanguages)

	// Create bordered widgets
	widgetStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("141")).
		Padding(1, 2)

	if showAll {
		// Only show the list widget when --all is used
		leftWidget := widgetStyle.Render(languageList)
		fmt.Println(leftWidget)
	} else {
		// Build the pie chart widget
		pieChart := renderPieChart(stats, topLanguages)

		// Count lines in each widget content to determine height
		listLines := strings.Count(languageList, "\n") + 1
		pieLines := strings.Count(pieChart, "\n") + 1

		// Pad the shorter content with empty lines to match heights
		if listLines > pieLines {
			diff := listLines - pieLines
			for i := 0; i < diff; i++ {
				pieChart += "\n"
			}
		} else if pieLines > listLines {
			diff := pieLines - listLines
			for i := 0; i < diff; i++ {
				languageList += "\n"
			}
		}

		leftWidget := widgetStyle.Render(languageList)
		rightWidget := widgetStyle.Render(pieChart)

		// Join widgets horizontally
		dashboard := lipgloss.JoinHorizontal(lipgloss.Top, leftWidget, "  ", rightWidget)
		fmt.Println(dashboard)
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

func buildLanguageList(stats *github.GithubStats, topLanguages []github.Language) string {
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

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Width(60)

	maxBytes := topLanguages[0].Bytes

	var lines []string
	lines = append(lines, headerStyle.Render("📊 Top Languages"))
	lines = append(lines, "")

	for i, lang := range topLanguages {
		// Bar width (max 30 chars for widget)
		barWidth := int(float64(lang.Bytes) / float64(maxBytes) * 30)
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

		line := fmt.Sprintf("%s %s %s %s",
			numberStyle.Render(fmt.Sprintf("%2d.", i+1)),
			langStyle.Render(fmt.Sprintf("%-15s", lang.Name)),
			barStyle.Render(bar),
			sizeStyle.Render(sizeWithPercent))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func renderPieChart(stats *github.GithubStats, topLanguages []github.Language) string {
	// Color palette for the pie chart
	colors := []string{"205", "218", "141", "213", "212", "183", "255", "251", "249", "247"}

	// Use top 5 languages for the pie chart
	pieLanguages := topLanguages
	if len(pieLanguages) > 5 {
		pieLanguages = pieLanguages[:5]
	}

	// Calculate "Others" bytes if there are more than 5 languages
	var othersBytes int
	if len(topLanguages) > 5 {
		for i := 5; i < len(topLanguages); i++ {
			othersBytes += topLanguages[i].Bytes
		}
	}

	var lines []string

	// Title - centered (grid is 19 chars wide with spacing, title is ~15 chars)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Width(19)

	lines = append(lines, headerStyle.Render("📈 Distribution"))
	lines = append(lines, "")

	// Simple ASCII pie chart using block characters
	chartLines := buildASCIIPie(pieLanguages, othersBytes, stats.TotalBytes, colors)
	lines = append(lines, chartLines...)

	lines = append(lines, "")

	// Legend
	for i, lang := range pieLanguages {
		percentage := float64(lang.Bytes) / float64(stats.TotalBytes) * 100
		colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i]))

		// Truncate language name if too long
		langName := lang.Name
		if len(langName) > 10 {
			langName = langName[:10]
		}

		legendLine := fmt.Sprintf("%s %-10s %5.1f%%",
			colorStyle.Render("●"),
			langName,
			percentage)
		lines = append(lines, legendLine)
	}

	// "Others" if there are more languages
	if othersBytes > 0 {
		percentage := float64(othersBytes) / float64(stats.TotalBytes) * 100
		grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		legendLine := fmt.Sprintf("%s %-10s %5.1f%%",
			grayStyle.Render("●"),
			"Others",
			percentage)
		lines = append(lines, legendLine)
	}

	return strings.Join(lines, "\n")
}

func buildASCIIPie(languages []github.Language, othersBytes int, totalBytes int, colors []string) []string {
	// Create a 10x10 rectangle (10 rows, 10 columns = 100 dots)
	const totalDots = 100
	const gridWidth = 10
	const gridHeight = 10

	// Calculate how many dots each language gets
	dots := make([]int, len(languages)+1) // +1 for "Others"
	remainingDots := totalDots

	for i, lang := range languages {
		if i >= len(colors) {
			break
		}
		percentage := float64(lang.Bytes) / float64(totalBytes)
		dots[i] = int(percentage * float64(totalDots))
		if dots[i] == 0 && percentage > 0 {
			dots[i] = 1 // Ensure at least 1 dot if there's any data
		}
		remainingDots -= dots[i]
	}

	// Add dots for "Others" category
	if othersBytes > 0 {
		percentage := float64(othersBytes) / float64(totalBytes)
		dots[len(languages)] = int(percentage * float64(totalDots))
		if dots[len(languages)] == 0 && percentage > 0 {
			dots[len(languages)] = 1 // Ensure at least 1 dot if there's any data
		}
		remainingDots -= dots[len(languages)]
	}

	// Assign remaining dots to the largest language to ensure we use all dots
	if len(dots) > 0 && remainingDots > 0 {
		dots[0] += remainingDots
	} else if remainingDots < 0 {
		// If we went over, take from the largest
		dots[0] += remainingDots
	}

	// Create a flat array of dot colors
	dotColors := make([]int, totalDots)
	dotIndex := 0
	for i, count := range dots {
		for j := 0; j < count && dotIndex < totalDots; j++ {
			// Use gray color (240) for "Others" category
			if i == len(languages) {
				dotColors[dotIndex] = -1 // Special marker for "Others"
			} else {
				dotColors[dotIndex] = i
			}
			dotIndex++
		}
	}

	// Render the grid
	var lines []string
	for row := 0; row < gridHeight; row++ {
		var line strings.Builder
		for col := 0; col < gridWidth; col++ {
			idx := row*gridWidth + col
			if idx < len(dotColors) {
				if dotColors[idx] == -1 {
					// "Others" category in gray
					colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
					line.WriteString(colorStyle.Render("●"))
				} else if dotColors[idx] < len(colors) {
					colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[dotColors[idx]]))
					line.WriteString(colorStyle.Render("●"))
				} else {
					line.WriteString(" ")
				}
			} else {
				line.WriteString(" ")
			}
			// Add spacing between dots
			if col < gridWidth-1 {
				line.WriteString(" ")
			}
		}
		lines = append(lines, line.String())
	}

	return lines
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
