package main

import (
	"fmt"
	"os"
	"strconv"

	"serverfans/internal/config"
	"serverfans/internal/ipmi"
	"serverfans/internal/ui"
)

func main() {
	ui.PrintBanner()

	if err := ipmi.CheckInstalled(); err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}

	cfg := config.Load()
	if cfg.IP != "" {
		ui.PrintInfo(fmt.Sprintf("Config loaded: %s@%s", cfg.User, cfg.IP))
	}

	for {
		if err := ui.RunForm(cfg); err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		speed, _ := strconv.Atoi(cfg.Speed)
		hexSpeed := fmt.Sprintf("0x%x", speed)

		if err := ipmi.SetFanSpeed(cfg.IP, cfg.User, cfg.Password, hexSpeed); err != nil {
			ui.PrintError(err)
			// Stay in loop — let user retry or fix the issue
			cfg.Speed = ""
			continue
		}

		ui.PrintSuccess(cfg.Speed)

		quit, err := ui.AskContinue()
		if err != nil || quit {
			break
		}
		cfg.Speed = ""
	}
}
