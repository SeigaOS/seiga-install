package main

import (
	"encoding/json"
	"fmt"
	"strings"

	// "log"
	"os"

	// "github.com/BurntSushi/toml"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/davecgh/go-spew/spew"
	// "github.com/davecgh/go-spew/spew"
)

// var sources = map[string] func() ([]huh.Option[T], error) {
//   "getKeymaps": func[T comparable]() ([]huh.Option[string], error) {
//   }
// }
// type SourceMap struct {
// 	Map map[string]Source[any]
// }
// type Source[T comparable] struct {
// 	Function    func() ([]huh.Option[T], error)
// 	Description string
// }

// var sources = SourceMap{
// 	Map: map[string]Source{
// 		"getKeymaps": {
// 			Function: func() ([]huh.Option[string], error) {
// 				return nil, nil
// 			},
// 			Description: "Get keymaps from /usr/share/kbd/keymaps",
// 		},
// 	},
// }

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
				// opt.Selected(false)
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
		// opts[0].Selected(tru/* e) */
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

var F *os.File

type Config struct {
	// The map keys are ids
	Groups map[string]Group `json:"groups"`
}

type Group struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Reference   *string          `json:"reference"` // Optional
	Inputs      map[string]Input `json:"inputs"`    // We have to loosely type because go has no support for multiple possible types in one field.
}

// TODO: Make this support more than just strings. How? No idea (probably codegen/macros?)
type Input struct {
	Title       string  `json:"title"`
	Description *string `json:"description"` // Optional... for now
	Type        string  `json:"type"`
	// We only care about choices for 'select' fields
	Options       []string `json:"options"`       // Optional
	OptionsSource *string  `json:"optionsSource"` // Optional

	Value  *string `json:"value"`  // This kinda serves a default value
	SkipIf *SkipIf `json:"skipIf"` // Optional, if nil then we automatically enable the field
}

type SkipIf struct {
	Field string  `json:"field"`
	Value *string `json:"value"`
}

func (c *Config) GenerateForm() (*huh.Form, error) {
	groups := make([]*huh.Group, 0, len(c.Groups))
	for _, group := range c.Groups {
		fields := make([]huh.Field, 0, len(group.Inputs))
		for _, input := range group.Inputs {
			if input.Value == nil {
				return nil, fmt.Errorf("input '%s' must have default value", input.Title)
			}
			var field huh.Field

			var desc string
			if input.Description != nil {
				desc += *input.Description + " "
			}
			desc += fmt.Sprintf("Default is %v", *input.Value)

			switch input.Type {
			case "boolean":
				confirm := huh.NewConfirm().Title(input.Title).Description(desc).Value(input.Value)
				field = confirm
			default:
				return nil, fmt.Errorf("'%s' is not a valid input type", input.Type)
			}
			if input.SkipIf != nil {
				jdir := strings.Split(input.SkipIf.Field, ".")
				field.SetSkipFunction(func() bool {
					other := c.Groups[jdir[0]].Inputs[jdir[1]].Value
					// fmt.Fprintf(F, "other: %v\n", other)
					return *other == *input.SkipIf.Value
				})
			}
			fields = append(fields, field)
		}
		group := huh.NewGroup(fields...).Title(group.Title).Description(group.Description)
		groups = append(groups, group)
	}
	return huh.NewForm(groups...), nil
}

// func (in *Input[T]) CreateField() (huh.Field, error) {
// 	// Everything must have a default value
// 	var desc string
// 	if in.Description != nil {
// 		desc += *in.Description + " "
// 	}
// 	desc += fmt.Sprintf("Default is %v", in.Value)

// 	switch in.Type {
// 	case "select":
// 		sel := huh.NewSelect[string]().Title(in.Title).Description(desc).Value(&in.Value)

// 		if in.OptionsSource == nil {
// 			if in.Options == nil {
// 				return nil, fmt.Errorf("'options' or 'optionsSource' required for a 'select' element")
// 			} else if len(in.Options) == 0 {
// 				return nil, fmt.Errorf("'options' must have at least one element if 'optionsSource' is not provided")
// 			}
// 		}
// 		var options []huh.Option[string]
// 		options = huh.NewOptions(in.Options...)
// 		if in.OptionsSource != nil {
// 			if f, ok := funcs[*in.OptionsSource]; ok {
// 				opts, err := f()
// 				if err != nil {
// 					return nil, err
// 				}
// 				if len(opts) == 0 {
// 					return nil, fmt.Errorf("no options returned from '%s'", *in.OptionsSource)
// 				}
// 				options = append(options, opts...)
// 			}
// 		}
// 		sel.Options(options...)
// 		sel.SetSelected(0) // Apparently setting the option to selected doesn't work? We also know there is at least one option

// 		if in.SkipIf != nil {
// 			jdir := strings.Split(in.SkipIf.Field, ".")
// 			other :=

// 				sel.SkipFunc(func() bool {})
// 		}
// 		return sel, nil
// 	case "boolean":
// 		conf := huh.NewConfirm().Title(in.Title).Description(desc).Value(&in.Value)

// 	default:
// 		return nil, fmt.Errorf("'%s' is not a valid input type", in.Type)
// 	}

// }

// func (c *Config) Form() (*huh.Form, error) {
// 	groups := make([]*huh.Group, 0, len(c.Groups))
// 	for _, group := range c.Groups {
// 		inputs := make([]huh.Field, 0, len(group.Inputs))
// 		for id, input := range group.Inputs {
// 			// https://github.com/charmbracelet/huh?tab=readme-ov-file#field-reference
// 			// field, err := input.CreateField()
// 			// if err != nil {
// 			// 	return nil, fmt.Errorf("field '%s' failed: '%v'", id, err)
// 			// }	// Everything must have a default value
// 			var desc string
// 			if input.Description != nil {
// 				desc += *input.Description
// 			}
// 			desc += fmt.Sprintf(" Default is %v", input.Value)

// 			switch input.Type {
// 			case "select":
// 				sel := huh.NewSelect[string]().Title(in.Title).Description(desc).Value(&input.Value)

// 				if input.OptionsSource == nil {
// 					if input.Options == nil {
// 						return nil, fmt.Errorf("'options' or 'optionsSource' required for a 'select' element")
// 					} else if len(input.Options) == 0 {
// 						return nil, fmt.Errorf("'options' must have at least one element if 'optionsSource' is not provided")
// 					}
// 				}
// 				var options []huh.Option[string]
// 				options = huh.NewOptions(input.Options...)
// 				if input.OptionsSource != nil {
// 					if f, ok := funcs[*input.OptionsSource]; ok {
// 						opts, err := f()
// 						if err != nil {
// 							return nil, err
// 						}
// 						if len(opts) == 0 {
// 							return nil, fmt.Errorf("no options returned from '%s'", *input.OptionsSource)
// 						}
// 						options = append(options, opts...)
// 					}
// 				}
// 				sel.Options(options...)
// 				sel.SetSelected(0) // Apparently setting the option to selected doesn't work? We also know there is at least one option

// 				if input.SkipIf != nil {
// 					jdir := strings.Split(input.SkipIf.Field, ".")

// 					sel.SkipFunc(func() bool {
// 						other := c.Groups[jdir[0]].Inputs[jdir[1]].Value
// 						return other == input.SkipIf.Value
// 					})
// 				}
// 				inputs = append(inputs, sel)
// 			case "boolean":
// 				conf := huh.NewConfirm().Title(input.Title).Description(desc).Value(&input.Value)

// 			default:
// 				return nil, fmt.Errorf("'%s' is not a valid input type", input.Type)
// 			}

// 		}
// 		groups = append(groups, huh.NewGroup(inputs...))
// 	}
// 	return huh.NewForm(groups...), nil
// }

func main() {
	logf, e := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		fmt.Println(e)
		return
	}
	F = logf
	f := "test.json"
	var config Config
	data, err := os.ReadFile(f)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(config)

	form, err := config.GenerateForm()
	// var test []string
	// test = nil
	// // test = []string{"test2"}
	// arr := []string{"test"}
	// arr = append(test, arr...)
	// // // spew.Dump(config)
	// fmt.Println(arr)
	// spew.Dump(config)
	// form, err := config.Form()
	// e := true
	// form := huh.NewForm(
	// 	huh.NewGroup(
	// 		huh.NewConfirm().Title("Enable?").Value(&e),
	// 		huh.NewSelect[string]().Title("Invis").Options(huh.NewOptions("a", "b", "c")...).SkipFunc(func() bool { return !e }),
	// 	),
	// 	huh.NewGroup(
	// 		huh.NewInput().Title("Secret Input"),
	// 	).WithHideFunc(func() bool { return !e }))
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
	spew.Dump(config)
}
