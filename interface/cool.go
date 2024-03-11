package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

// problemnetstat -i
// list all interfaces for user select
// netstat -i  also a thing

func main() {
	cmd := exec.Command("netstat", "-i")
	output, _ := cmd.CombinedOutput()
	//fmt.Println(string(output))
	temp := strings.Split(string(output), "\n")
	fmt.Println("   " + temp[0])
	var timezonesplace []string

	options := make([]huh.Option[string], 0, len(timezonesplace))

	for _, e := range temp {
		if strings.Contains(e, "Name") || e == "" {
			continue
		} else {
			timezonesplace = append(timezonesplace, e)
			options = append(options, huh.NewOption[string](e, e))
		}
	}
	var timezone string
	huh.NewSelect[string]().
		Title("Pick a interface.").
		Options(options...).
		Value(&timezone).
		Run()

	fmt.Println("great ?!", timezone)

}
