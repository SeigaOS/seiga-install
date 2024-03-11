package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	_ "time/tzdata"

	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/pelletier/go-toml"
)

//test both
// func setSystemTimezone(timezone string) {
// 	zoneinfoPath := "/usr/share/zoneinfo/" + timezone
// 	localtimePath := "/etc/localtime"
// 	if err := os.Remove(localtimePath); err != nil {
// 		log.Fatalf("Failed to remove old timezone link: %v", err)
// 	}
// 	if err := os.Symlink(zoneinfoPath, localtimePath); err != nil {
// 		log.Fatalf("Failed to create symlink for new timezone: %v", err)
// 	}
// }

func setSystemTimezone(timezone string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("sudo", "ln", "-sf", "/usr/share/zoneinfo/"+timezone, "/etc/localtime")
	case "darwin":
		cmd = exec.Command("sudo", "systemsetup", "-settimezone", timezone)
	default:
		return fmt.Errorf("setting timezone not supported on this OS")
	}

	//run
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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

			//fmt.Println(e.Name())
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
	elpath := "/usr/share/zoneinfo/" + timezone

	if stat, err := os.Stat(elpath); err == nil && stat.IsDir() {
		fmt.Println("hey")
		configFile, err := os.ReadFile("config.toml")
		var config Timeconf
		newConfig, err := toml.Marshal(config)
		newone, err := os.ReadDir("/usr/share/zoneinfo/" + timezone)
		var secondtz string
		newoptions := make([]huh.Option[string], 0, len(timezonesplace))
		for _, e := range newone {
			//fmt.Println(e.Name())
			newoptions = append(newoptions, huh.NewOption[string](e.Name(), e.Name()))
		}
		huh.NewSelect[string]().
			Title("Pick a timezone.").
			Options(newoptions...).
			Value(&secondtz).
			Run()
		f, err := os.Create("config.toml")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err != nil {
			//fmt.Println(timezone)
			log.Fatal(err)
		}

		if err := toml.Unmarshal(configFile, &config); err != nil {
			panic(err)
		}

		if err != nil {
			panic(err)
		}
		var realtimezone = timezone + "/" + secondtz
		config.Place = realtimezone
		if err := os.WriteFile("config.toml", newConfig, 0644); err != nil {
			panic(err)
		}
		fmt.Println(realtimezone)
		if err != nil {
			fmt.Println("ok great theres an error")
		}
		fmt.Println("great ?! we've set your timezone as", realtimezone)
		setSystemTimezone(realtimezone)

	} else {

		var config Timeconf
		fmt.Println("great?! we've set your path as " + timezone)
		setSystemTimezone(timezone)

		configFile, err := os.ReadFile("config.toml")
		config.Place = timezone
		newConfig, err := toml.Marshal(config)
		if err != nil {
			panic(err)
		}
		if err := toml.Unmarshal(configFile, &config); err != nil {
			panic(err)
		}
		if err := os.WriteFile("config.toml", newConfig, 0644); err != nil {
			panic(err)
		}

	}

}
