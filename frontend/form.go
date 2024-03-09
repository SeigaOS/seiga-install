package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

type Config struct {
}

func main() {
	// form := huh.NewForm(
	// 	huh.NewGroup(
	// 		// Ask the user for a base burger and toppings.
	// 		huh.NewSelect[string]().
	// 			Title("Choose your burger").
	// 			Options(
	// 				huh.NewOption("Charmburger Classic", "classic"),
	// 				huh.NewOption("Chickwich", "chickwich"),
	// 				huh.NewOption("Fishburger", "fishburger"),
	// 				huh.NewOption("Charmpossible™ Burger", "charmpossible"),
	// 			).
	// 			Value(&burger), // store the chosen option in the "burger" variable

	// 		// Let the user select multiple toppings.
	// 		huh.NewMultiSelect[string]().
	// 			Title("Toppings").
	// 			Options(
	// 				huh.NewOption("Lettuce", "lettuce").Selected(true),
	// 				huh.NewOption("Tomatoes", "tomatoes").Selected(true),
	// 				huh.NewOption("Jalapeños", "jalapeños"),
	// 				huh.NewOption("Cheese", "cheese"),
	// 				huh.NewOption("Vegan Cheese", "vegan cheese"),
	// 				huh.NewOption("Nutella", "nutella"),
	// 			).
	// 			Limit(4). // there’s a 4 topping limit!
	// 			Value(&toppings),

	// 		// Option values in selects and multi selects can be any type you
	// 		// want. We’ve been recording strings above, but here we’ll store
	// 		// answers as integers. Note the generic "[int]" directive below.
	// 		huh.NewSelect[int]().
	// 			Title("How much Charm Sauce do you want?").
	// 			Options(
	// 				huh.NewOption("None", 0),
	// 				huh.NewOption("A little", 1),
	// 				huh.NewOption("A lot", 2),
	// 			).
	// 			Value(&sauceLevel),
	// 	),

	// 	// Gather some final details about the order.
	// 	huh.NewGroup(
	// 		huh.NewInput().
	// 			Title("What's your name?").
	// 			Value(&name).
	// 			// Validating fields is easy. The form will mark erroneous fields
	// 			// and display error messages accordingly.
	// 			Validate(func(str string) error {
	// 				if str == "Frank" {
	// 					return errors.New("Sorry, we don’t serve customers named Frank.")
	// 				}
	// 				return nil
	// 			}),

	// 		huh.NewText().
	// 			Title("Special Instructions").
	// 			CharLimit(400).
	// 			Value(&instructions),

	// 		huh.NewConfirm().
	// 			Title("Would you like 15% off?").
	// 			Value(&discount),
	// 	),
	// ) // Create a new form
	// err := form.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if !discount {
	// 	fmt.Println("What? You didn’t take the discount?!")
	// }

}
