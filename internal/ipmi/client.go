package ipmi

import (
	"fmt"
	"os/exec"
)

// CheckInstalled verifies ipmitool is available on the system PATH.
func CheckInstalled() error {
	if err := exec.Command("ipmitool", "-V").Run(); err != nil {
		return fmt.Errorf("ipmitool is not installed or not found in PATH")
	}
	return nil
}

// SetFanSpeed enables manual fan control on the target host and sets all fans
// to the requested speed. speedHex must be a hex string like "0x32" (50%).
func SetFanSpeed(host, user, password, speedHex string) error {
	base := []string{"-I", "lanplus", "-H", host, "-U", user, "-P", password, "raw", "0x30", "0x30"}

	// Enable manual fan control mode
	if err := exec.Command("ipmitool", append(base, "0x01", "0x00")...).Run(); err != nil {
		return fmt.Errorf("failed to enable manual fan control: %w", err)
	}

	// Set all fans to the requested speed
	if err := exec.Command("ipmitool", append(base, "0x02", "0xff", speedHex)...).Run(); err != nil {
		return fmt.Errorf("failed to set fan speed: %w", err)
	}

	return nil
}
