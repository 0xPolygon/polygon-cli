package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

// updateReadme will update the list of `polycli` commands.
// The section is identified by the HTML tags `<startTag></endTag>“.
func updateReadmeCommands(cmd *cobra.Command, startTag, endTag string) error {
	// Generate the list of commands.
	buf := new(bytes.Buffer)
	name := cmd.CommandPath()
	addDocPrefix := func(s string) string { return fmt.Sprintf("%s/%s", docDir, s) }
	printSeeAlso(buf, cmd, name, addDocPrefix)

	// Update the `README.md``
	data, err := ioutil.ReadFile("README.md")
	if err != nil {
		return err
	}
	var newData string
	newData, err = updateContent(string(data), startTag, endTag, buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("README.md", []byte(newData), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Take a piece of data and update the content between the start and end tags with new content.
func updateContent(originalContent, startTag, endTag string, newContent *bytes.Buffer) (string, error) {
	startIndex := strings.Index(originalContent, startTag)
	endIndex := strings.Index(originalContent, endTag)
	if startIndex == -1 || endIndex == -1 || endIndex <= startIndex {
		return "", fmt.Errorf("Unable to find start and end tags or they are in the wrong order.")
	}
	startIndex += len(startTag)
	return originalContent[:startIndex] + "\n" + newContent.String() + originalContent[endIndex:], nil
}
