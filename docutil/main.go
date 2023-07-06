package main

import (
	"fmt"
	"log"

	"github.com/maticnetwork/polygon-cli/cmd"
)

// Directory in which the documentation will be generated.
var docDir = "doc"

func main() {
	// Generate documentation for the `polycli` command.
	polycli := cmd.NewPolycliCommand()
	if err := genMarkdownDoc(polycli, docDir); err != nil {
		fmt.Println("Unable to generate documentation.")
		log.Fatal(err)
	}
	fmt.Println("Documentation generated!")
}
