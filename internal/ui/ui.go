package ui

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"serverfans/internal/config"
)

const (
	MinFanSpeed = 10
	MaxFanSpeed = 100
)

// Catppuccin Mocha palette
var (
	colorMauve  = lipgloss.Color("#cba6f7")
	colorBlue   = lipgloss.Color("#89b4fa")
	colorGreen  = lipgloss.Color("#a6e3a1")
	colorRed    = lipgloss.Color("#f38ba8")
	colorYellow = lipgloss.Color("#f9e2af")
	colorBase   = lipgloss.Color("#1e1e2e")

	bannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMauve).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBlue).
			Padding(0, 4).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Italic(true)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBase).
			Background(colorGreen).
			Padding(0, 2)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBase).
			Background(colorRed).
			Padding(0, 2)

	infoStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Padding(0, 1)

	ipRegex = regexp.MustCompile(
		`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])$`,
	)
)

// PrintBanner prints the application header to stdout.
func PrintBanner() {
	title := bannerStyle.Render("ServerFans\n" + subtitleStyle.Render("IPMI Fan Speed Controller"))
	fmt.Println()
	fmt.Println(title)
	fmt.Println()
}

// RunForm displays a single grouped form for all connection and speed inputs.
// Pre-filled values from cfg are shown as defaults and skipped by the user.
func RunForm(cfg *config.Config) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server IP").
				Description("IPMI IP address of the target server").
				Placeholder("192.168.1.100").
				Value(&cfg.IP).
				Validate(func(s string) error {
					if !ipRegex.MatchString(s) {
						return fmt.Errorf("invalid IP address")
					}
					return nil
				}),

			huh.NewInput().
				Title("Username").
				Description("IPMI user with power user or admin privileges").
				Value(&cfg.User).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("username cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Description("IPMI user password").
				Value(&cfg.Password).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("password cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Fan Speed (%)").
				Description(fmt.Sprintf("Target speed percentage (%d–%d)", MinFanSpeed, MaxFanSpeed)).
				Placeholder("50").
				Value(&cfg.Speed).
				Validate(func(s string) error {
					n, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("must be a whole number")
					}
					if n < MinFanSpeed || n > MaxFanSpeed {
						return fmt.Errorf("must be between %d and %d", MinFanSpeed, MaxFanSpeed)
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
}

// PrintSuccess prints a styled success message.
func PrintSuccess(speed string) {
	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf(" Fan speed set to %s%% ", speed)))
	fmt.Println()
}

// PrintError prints a styled error message.
func PrintError(err error) {
	fmt.Println()
	fmt.Println(errorStyle.Render(fmt.Sprintf(" Error: %s ", err.Error())))
	fmt.Println()
}

// PrintInfo prints a styled informational message.
func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render(msg))
}

// AskContinue shows a confirm dialog and returns true if the user wants to quit.
func AskContinue() (bool, error) {
	var quit bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("What would you like to do next?").
				Affirmative("Quit").
				Negative("Set another speed").
				Value(&quit),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	return quit, err
}
