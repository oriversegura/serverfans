package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/charmbracelet/huh"
)

func main() {

	//Const Necesary to logical use
	const minFanSpeed = 10
	const maxFanSpeed = 100
	/*ipmiToolPath := []string {
		"/usr/local/bin/ipmitool",
		"/usr/bin/ipmitol",
		"C:\\ipmitool",
	}*/

	// declare variables to use
	var user, ip string
	var fanSpeed int

	//Validate Ipmi is installed on system

	// Valite ipmitool is installed
	cmdValidate := exec.Command("whereis", "ipmitool")
	if err := cmdValidate.Run(); err != nil {
		log.Fatal(err)
	}

	// Get ip and verify is valid with regex
	fmt.Printf("Insert server ip: ")
	fmt.Scanf("%s", &ip)
	pattern := "^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])$"
	match, err := regexp.MatchString(pattern, ip)
	if err != nil {
		log.Fatal(err)
	}
	if match != true {
		fmt.Errorf("Insert a Valid IP")
	}

	// Get the user and save in enviroment variable
	fmt.Printf("%s", "Insert server user: ")
	fmt.Scan(&user)

	// Get the password
	var password string

	// Get the Password and save in enviroment variable
	// New Form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server Password").
				EchoMode(huh.EchoModePassword).
				Description("Enter server Password: ").
				Value(&password),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get user Fan Speed in percent
	fmt.Printf("Insert Fan Speed %d to %d percent: ", minFanSpeed, maxFanSpeed)
	fmt.Scanf("%d", &fanSpeed)
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

	fmt.Printf("Ventiladores al %d%", fanSpeed)

}
