package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Split the string
func spiltString(str string) []string {
	return strings.Split(str, ",")
}

// Handle error
func handleError(err string, exitOnError bool) {
	fmt.Println("\n" + err)
	if exitOnError {
		os.Exit(1)
	}
}

// Prompt for confirmation
func yesOrNoConfirmation(appName string) string {
	var YesOrNo = map[string]string{"y": "y", "ye": "y", "yes": "y", "n": "n", "no": "n"}

	// Start the new scanner to get the user input
	fmt.Printf("Are you sure you want to delete these apps (%s), do you wish to continue (Yy/Nn)?: ", appName)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		// The choice entered
		choiceEntered := input.Text()

		// If its a valid value move on
		if YesOrNo[strings.ToLower(choiceEntered)] == "y" { // Is it Yes
			return choiceEntered
		} else if YesOrNo[strings.ToLower(choiceEntered)] == "n" { // Is it No
			handleError("Canceling as per user request", true)
		} else { // Invalid choice, ask to re-enter
			fmt.Println("Invalid Choice: Please enter Yy/Nn, try again.")
			return yesOrNoConfirmation(appName)
		}
	}
	return ""
}

// Check if manifest file exists in the current working directory
func readManifest() []byte {
	manifestFileName := "manifest.yml"
	currentDirectory, err := os.Getwd()
	fullPath := fmt.Sprintf("%s/%s", currentDirectory, manifestFileName)
	if err != nil {
		handleError(fmt.Sprintf("Unable to get the current working directory, err: %v", err), true)
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		handleError(fmt.Sprintf("Unable to file the manifest file \"%s\"", fullPath), true)
	} else {
		b, err := ioutil.ReadFile(fullPath) // b has type []byte
		if err != nil {
			handleError(fmt.Sprintf("Error when reading the manifest \"%s\", err: %v",fullPath, err), true)
		}
		return b
	}
	return []byte{}
}