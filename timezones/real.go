package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	_ "time/tzdata"

	"github.com/charmbracelet/huh"
	"github.com/pelletier/go-toml"
)

func setSystemTimezone(timezone string) {
	zoneinfoPath := "/usr/share/zoneinfo/" + timezone
	localtimePath := "/etc/localtime"
	if err := os.Remove(localtimePath); err != nil {
		log.Fatalf("Failed to remove old timezone link: %v", err)
	}
	if err := os.Symlink(zoneinfoPath, localtimePath); err != nil {
		log.Fatalf("Failed to create symlink for new timezone: %v", err)
	}
}

// /usr/share/zoneinfo
type Timeconf struct {
	Place string `toml:"place"`
}

func main() {

	entries, err := os.ReadDir("/usr/share/zoneinfo")

	if err != nil {
		log.Fatal(err)
	}
	var timezonesplace []string

	options := make([]huh.Option[string], 0, len(timezonesplace))

	for _, e := range entries {
		if strings.Contains(e.Name(), "+") {
			continue
		} else {

			fmt.Println(e.Name())
			timezonesplace = append(timezonesplace, e.Name())
			options = append(options, huh.NewOption[string](e.Name(), e.Name()))
		}

	}

	var timezone string
	huh.NewSelect[string]().
		Title("Pick a timezone.").
		Options(options...).
		Value(&timezone).
		Run()

	fmt.Println("great ?!", timezone)
	newone, err := os.ReadDir("/usr/share/zoneinfo/" + timezone)
	if err != nil {
		log.Fatal(err)
	}
	var secondtz string
	newoptions := make([]huh.Option[string], 0, len(timezonesplace))

	for _, e := range newone {
		fmt.Println(e.Name())
		newoptions = append(newoptions, huh.NewOption[string](e.Name(), e.Name()))
	}
	huh.NewSelect[string]().
		Title("Pick a timezone.").
		Options(newoptions...).
		Value(&secondtz).
		Run()
	fmt.Println("great ?!", timezone)
	f, err := os.Create("config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	configFile, err := os.ReadFile("config.toml")
	var config Timeconf
	if err := toml.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	var realtimezone = timezone + "/" + secondtz
	config.Place = realtimezone

	newConfig, err := toml.Marshal(config)
	if err := os.WriteFile("config.toml", newConfig, 0644); err != nil {
		panic(err)
	}
	fmt.Println(realtimezone)
	if err != nil {

	}
	setSystemTimezone(realtimezone)
}
