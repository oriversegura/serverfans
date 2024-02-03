package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	ip := "10.11.11.4" // Reemplaza esto con tu dirección IP

	// Primero, establecer el control manual del ventilador
	cmd1 := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", "root", "-P", "CSC2023@@", "raw", "0x30", "0x30", "0x01", "0x00")
	if err := cmd1.Run(); err != nil {
		log.Fatal(err)
	}

	// Segundo, establecer los ventiladores al 20%
	cmd2 := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", "root", "-P", "CSC2023@@", "raw", "0x30", "0x30", "0x02", "0xff", "0x14")
	if err := cmd2.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ventiladores al 20%")
}
