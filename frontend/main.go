package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"net"
	"strconv"
)

type Config struct {
	NetworkConfig NetworkConfig
	PacmanConfig  PacmanConfig
}

type NetworkConfig struct {
	Interface string
}

type PacmanConfig struct {
	ParallelDownloads string
	Multilib          bool
	Color             bool
	PacmanDesign      bool
}

func DefaultConfig() Config {
	return Config{
		PacmanConfig: PacmanConfig{
			ParallelDownloads: "4",
			Multilib:          true,
			Color:             true,
			PacmanDesign:      false,
		},
		NetworkConfig: NetworkConfig{
			Interface: "wlan0",
		},
	}
}

// Helper funcs ----
func isInt(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("'%s' is not a number", s)
	}
	return nil
}

// Fixing huh ----
func Group(title string, fields ...huh.Field) *huh.Group {
	return huh.NewGroup(fields...).Title(title)
}
func StrInput(title string, val *string) *huh.Input {
	return huh.NewInput().Title(title).Value(val)
}

func IntInput(title string, val *string) *huh.Input {
	return StrInput(title, val).Validate(isInt)
}

func Confirm(title string, val *bool) *huh.Confirm {
	return huh.NewConfirm().Title(title).Value(val)
}

func Select[T comparable](title string, opts []huh.Option[T], val *T) *huh.Select[T] {
	return huh.NewSelect[T]().Title(title).Options(opts...).Value(val)
}

func main() {
	defaultConfig := DefaultConfig()
	config := Config{}
	// Create a form with huh that sets timezone in  Arch linux live environment
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error getting network interfaces: %v", err)
	}
	iface_opts := make([]huh.Option[string], 0, len(ifaces))
	for _, iface := range ifaces {
		iface_opts = append(iface_opts, huh.NewOption(iface.Name, iface.Name))
	}

	form := huh.NewForm(
		Group("Network Configuration", Select("Interface", iface_opts, &config.NetworkConfig.Interface)),
		Group("Pacman Configuration",
			IntInput("Parallel Downloads", &config.PacmanConfig.ParallelDownloads),
			Confirm("Enable Multilib", &config.PacmanConfig.Multilib),
			Confirm("Enable Color", &config.PacmanConfig.Color),
			Confirm("Enable Pacman Design", &config.PacmanConfig.PacmanDesign),
		),
	)
	form.Run()
	fmt.Println(config)

}
