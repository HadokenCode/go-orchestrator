package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/colorstring"
)

const (
	//
	// Orchestrator commands
	//

	// CmdRelease To start a release
	CmdRelease = "release"
	// CmdInfo to display informations about a release (catalog)
	CmdInfo = "info"
	// CmdDisplayUserConfiguration Display User Configuration
	CmdDisplayUserConfiguration = "show-user-config"
	// CmdDisplayCatalog To display the catalog of all available projects
	CmdDisplayCatalog = "show-catalog"
	// CmdDisplayRunningContainers display all running containers
	CmdDisplayRunningContainers = "show-running-containers"
	// CmdDisplayAllContainers display all containers
	CmdDisplayAllContainers = "show-all-containers"
    // CmdDisplayReleasesStatus Show status for all unit releases
    CmdDisplayReleasesStatus = "show-releases-status"
    
    // LabelAll all available labels
    LabelAll = "all"
)

// PrintError Display a error message with red color
func PrintError(message string, err error) {
	if message != "" {
		fmt.Fprintf(os.Stderr, colorstring.Color("[red]"+message))
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}

// UsageAndExit Display commands usage and exit
func UsageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}

// AskForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

// You might want to put the following two functions in a separate utility package.

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
