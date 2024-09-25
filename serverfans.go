package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/charmbracelet/huh"
)

func main() {

	//Const Necesary to logical use
	const minFanSpeed = 10
	const maxFanSpeed = 100

	//Confirm ipmitool is installed
	if exec.Command("ipmitool", "-V").Run() != nil {
		log.Fatal("Ipmi is no installed on system!")
	}

	// declare variables to use
	var user, ip, fanSpeedT string

	//Validate Ipmi is installed on system
	cmdValidate := exec.Command("whereis", "ipmitool")
	if err := cmdValidate.Run(); err != nil {
		log.Fatal(err)
	}

	// // Get ip and verify is valid with regex
	// fmt.Printf("Insert server ip: ")
	// fmt.Scanf("%s", &ip)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server IP").
				EchoMode(huh.EchoModeNormal).
				Description("Enter the server IP: ").
				Value(&ip),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	pattern := "^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])$"
	match, err := regexp.MatchString(pattern, ip)
	if err != nil {
		log.Fatal(err)
	}
	if match != true {
		fmt.Errorf("Insert a Valid IP")
	}

	// Get the user and save in enviroment variable
	// fmt.Printf("%s", "Insert server user: ")
	// fmt.Scan(&user)
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server User").
				EchoMode(huh.EchoModeNormal).
				Description("Enter server user: ").
				Value(&user),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get the password
	var password string

	// Get the Password and save in enviroment variable
	// New Form
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server Password").
				EchoMode(huh.EchoModePassword).
				Description("Enter server password: ").
				Value(&password),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get user Fan Speed in percent
	//
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Fan Speed").
				EchoMode(huh.EchoModeNormal).
				Description("Insert Fan Speed 10 to 100 percent: ").
				Value(&fanSpeedT),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Convert to Int
	fanSpeed, err := strconv.Atoi(fanSpeedT)
	if err != nil {
		log.Fatal(err)
	}

	//validate Fan Speed is Correct
	if fanSpeed < minFanSpeed || fanSpeed > maxFanSpeed {
		log.Fatal("Insert Valid Fan Speed")
	}
	// Print Valid Fan Speed on second Command
	hexString := fmt.Sprintf("0x%x", fanSpeed)

	// Primero, establecer el control manual del ventilador
	cmd1 := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", user, "-P", password, "raw", "0x30", "0x30", "0x01", "0x00")
	if err := cmd1.Run(); err != nil {
		log.Fatal(err)
	}

	// Segundo, establecer los ventiladores al 20%
	cmd2 := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", user, "-P", password, "raw", "0x30", "0x30", "0x02", "0xff", hexString)
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Fan Speed set to " + fanSpeedT + " percent").
				Affirmative("Quit").
				Negative("Quit"),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("Fan Speed set to %d% \n", fanSpeed)

}
