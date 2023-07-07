package main

import (
	"fmt"
	"log"

	"github.com/maticnetwork/polygon-cli/cmd"
)

// Directory in which the documentation will be generated.
var (
	docDir   = "doc"
	startTag = "<generated>"
	endTag   = "</generated>"
)

func main() {
	polycli := cmd.NewPolycliCommand()

	// Generate documentation for the `polycli` command.
	if err := genMarkdownDoc(polycli, docDir); err != nil {
		fmt.Println("Unable to generate documentation.")
		log.Fatal(err)
	}
	fmt.Println("Documentation generated!")

	// Update the summary of commands in the `README.md`
	if err := updateReadmeCommands(polycli, startTag, endTag); err != nil {
		fmt.Println("Unable to update `README.md`.")
		log.Fatal(err)
	}
	fmt.Println("`README.md` updated!")
}
