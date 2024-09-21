package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/charmbracelet/huh"
)

func main() {
	// p := tea.NewProgram(model{}, tea.WithReportFocus())
	// New Form
	form := huh.NewForm()

	// declare variables to use
	var user, ip string
	var fanSpeed int

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
	os.Setenv("IPMI_USER", user)

	// Get the password
	var password string

	// Get the Password and save in enviroment variable
	form.Append(
		huh.NewInput().
			Title("Server Password").
			EchoMode(huh.EchoModePassword).
			Description("Enter server Password: ").
			Value(&password),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Insert fan speed in percent
	fmt.Printf("%s", "Insert Speed 10 to 100 percent: ")
	fmt.Scanf("%d", &fanSpeed)
	hexString := fmt.Sprintf("0x%x", fanSpeed)

	// Primero, establecer el control manual del ventilador
	cmd1 := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", os.Getenv("IPMI_USER"), "-P", os.Getenv("IPMI_PASS"), "raw", "0x30", "0x30", "0x01", "0x00")
	if err := cmd1.Run(); err != nil {
		log.Fatal(err)
	}

	// Segundo, establecer los ventiladores al 20%
	cmd2 := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", os.Getenv("IPMI_USER"), "-P", os.Getenv("IPMI_PASS"), "raw", "0x30", "0x30", "0x02", "0xff", hexString)
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ventiladores al %d%", fanSpeed)
}
