package main

import (
	"encoding/json"
	"fmt"
	"strings"

	// "log"
	"os"

	// "github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"path/filepath"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/davecgh/go-spew/spew"
)

var funcs = map[string]func() ([]huh.Option[string], error){
	"getKeymaps": func() ([]huh.Option[string], error) {
		opts := make([]huh.Option[string], 0)
		keymapDir := "/usr/share/kbd/keymaps"
		err := filepath.Walk(keymapDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				basename := info.Name()
				name := strings.TrimSuffix(basename, ".map.gz") // Ext is probably .map.gz
				opt := huh.NewOption(name, name)
				opts = append(opts, opt)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if len(opts) == 0 {
			return nil, fmt.Errorf("no keymaps found")
		}
		return opts, nil

	},
}

// Don't worry about it
// func AtoOpts[T comparable](a []T) []huh.Option[T] {
// 	opts := make([]huh.Option[T], 0, len(a))
// 	for _, v := range a {
// 		opts = append(opts, huh.Option[T]{
// 			Key:   fmt.Sprint(v),
// 			Value: v,
// 		})
// 	}
// 	return opts
// }

type Config struct {
	Groups []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Inputs      []struct {
			Title           string   `json:"title"`
			Name            string   `json:"name"`
			Description     string   `json:"description"`
			Type            string   `json:"type"`
			Choices         []string `json:"choices"`         // Optional
			ChoicesFunction *string  `json:"choicesFunction"` // Optional
			Value           string   `json:"value"`           // Optional?
		} `json:"inputs"`
	} `json:"groups"`
}

func (c *Config) Form() (*huh.Form, error) {
	groups := make([]*huh.Group, 0, len(c.Groups))
	for _, group := range c.Groups {
		inputs := make([]huh.Field, 0, len(group.Inputs))
		for _, input := range group.Inputs {
			// https://github.com/charmbracelet/huh?tab=readme-ov-file#field-reference
			switch input.Type {
			case "select":
				selectEl := huh.NewSelect[string]().Title(input.Title).Description(fmt.Sprintf("%s Default is: `%s`", input.Description, input.Value)).Value(&input.Value)
				if input.Choices == nil && input.ChoicesFunction == nil {
					return nil, fmt.Errorf("`choices` or `choicesFunction` is required for select input")
				} else if input.Choices != nil {
					selectEl.Options(huh.NewOptions(input.Choices...)...)
				} else {
					if f, ok := funcs[*input.ChoicesFunction]; ok {
						opts, err := f()
						if err != nil {
							return nil, err
						}
						selectEl.Options(opts...)
						selectEl.SetOptions(opts...) // This is a test to see if the patch works
						selectEl.SetSelected(0)      // This framework kinda sucks
					} else {
						return nil, fmt.Errorf("choices function `%s` not found", *input.ChoicesFunction)
					}
				}
				inputs = append(inputs, selectEl)
			case "multiselect":
				return nil, fmt.Errorf("multiselect not implemented")
			case "input":
				return nil, fmt.Errorf("input not implemented")
			case "text":
				return nil, fmt.Errorf("text not implemented")
			case "confirm":
				return nil, fmt.Errorf("confirm not implemented")
			default:
				return nil, fmt.Errorf("unknown input type: %s", input.Type)
			}

		}
		groups = append(groups, huh.NewGroup(inputs...))
	}
	return huh.NewForm(groups...), nil
}

// var SEL = huh.NewSelect[string]().Title("Second").Options(huh.NewOptions("b", "c")...)

func main() {
	f := "schema.json"
	var config Config
	data, err := os.ReadFile(f)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	// spew.Dump(config)
	form, err := config.Form()
	// spew.Dump(form)
	// spew.Dump(form)
	// inputs := []huh.Field{
	// 	huh.NewSelect[string]().Title("Select").Description("Select").Options(AtoOpts([]string{"a", "b", "c"})...),
	// }
	// groups := []*huh.Group{huh.NewGroup(inputs...)}
	// form := huh.NewForm(groups...)
	// spew.Dump(form)

	// var sel *huh.Select[string]
	// sel = huh.NewSelect[string]().Title("First").Options(huh.NewOptions("default", "a", "aasdfas", "sdfs", "sdfs")...)
	// form := huh.NewForm(
	// 	huh.NewGroup(
	// 		huh.NewSelect[string]().Title("First").Validate(func(s string) error {
	// 			switch s {
	// 			case "a":
	// 				sel.SetOptions(huh.NewOptions("it", "was", "a"))
	// 			default:
	// 				sel.SetOptions(huh.NewOptions("it", "wasn't", "a"))
	// 			}
	// 			// spew.Dump(sel)
	// 			return nil
	// 		}).Options(huh.NewOptions("a", "b", "c")...),
	// 		sel,
	// 	),

	// )
	if err != nil {
		fmt.Println(err)
		return
	}
	err = form.Run()
	if err != nil {
		fmt.Println(err)
	}
}
