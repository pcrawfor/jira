package main

import (
	"fmt"
	"os"

	"github.com/pcrawfor/jira"
)

const username = ""
const password = ""
const siteURL = ""

func main() {
	c := jira.NewJiraClient(siteURL, username, password, 100)
	issues, ierr := c.Search("status=reviewed")
	if ierr != nil {
		fmt.Println("Error: ", ierr)
		os.Exit(0)
	}

	for _, issue := range issues {
		fmt.Println("Issue:", issue)
	}
}
