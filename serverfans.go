package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"serverfans/config"
	"strconv"

	"github.com/charmbracelet/huh"
)

func main() {

	configuration, err := config.LoadConfig("/home/osegura/Software/Golang/serverfans/config.json")
	if err != nil {
		log.Fatalf("Mamense un guebo los que inventaron json")
	}

	// Const Necesary to logical use
	const minFanSpeed = 10
	const maxFanSpeed = 100

	// Confirm ipmitool is installed
	if exec.Command("ipmitool", "-V").Run() != nil {
		log.Fatal("Ipmi is no installed on system!")
	}

	// Validate Ipmi is installed on system
	cmdValidate := exec.Command("whereis", "ipmitool")
	if err := cmdValidate.Run(); err != nil {
		log.Fatal(err)
	}

	// // Get ip and verify is valid with regex
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server IP").
				EchoMode(huh.EchoModeNormal).
				Description("Enter the server IP: ").
				Value(&configuration.IP),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	pattern := "^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])$"
	match, err := regexp.MatchString(pattern, configuration.IP)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		fmt.Errorf("Insert a Valid IP")
	}

	// Get the user and save in enviroment variable
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server User").
				EchoMode(huh.EchoModeNormal).
				Description("Enter server user: ").
				Value(&configuration.User),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get the Password
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server Password").
				EchoMode(huh.EchoModePassword).
				Description("Enter server password: ").
				Value(&configuration.Password),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get user Fan Speed in percent
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Fan Speed").
				EchoMode(huh.EchoModeNormal).
				Description("Insert Fan Speed 10 to 100 percent: ").
				Value(&configuration.Speed),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Convert to Int
	fanSpeed, err := strconv.Atoi(configuration.Speed)
	if err != nil {
		log.Fatal("Se murio la velocidad!")
	}

	// validate Fan Speed is Correct
	if fanSpeed < minFanSpeed || fanSpeed > maxFanSpeed {
		log.Fatal("Insert Valid Fan Speed")
	}
	// Print Valid Fan Speed on second Command
	hexString := fmt.Sprintf("0x%x", fanSpeed)

	// Primero, establecer el control manual del ventilador
	cmd1 := exec.Command("ipmitool", "-I", "lanplus", "-H", configuration.IP, "-U", configuration.User, "-P", configuration.Password, "raw", "0x30", "0x30", "0x01", "0x00")
	if err := cmd1.Run(); err != nil {
		log.Fatal(err)
	}

	// Segundo, establecer los ventiladores al 20%
	cmd2 := exec.Command("ipmitool", "-I", "lanplus", "-H", configuration.IP, "-U", configuration.User, "-P", configuration.Password, "raw", "0x30", "0x30", "0x02", "0xff", hexString)
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Fan Speed set to " + configuration.Speed + " percent").
				Affirmative("quit").
				Negative("start again"),
		),
	).WithTheme(huh.ThemeCatppuccin())

	for {
		if err := form.Run(); err != nil {
			log.Fatal(err)
		}

		buttons := form.Get("")

		if buttons == true {
			break
		} else {
			main()
		}
	}
}
