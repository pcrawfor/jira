/*
Package jira provides a simple wrapper around calls to the JIRA REST API

http://godoc.org/github.com/pcrawfor/jira

// basic example searching for issues with a status of reviewed
c := jira.NewJiraClient(siteURL, username, password, 100)
issues, ierr := c.Search("status=reviewed")
if ierr != nil {
	fmt.Println("Error: ", ierr)
	os.Exit(0)
}

for _, issue := range issues {
	fmt.Println("Issue:", issue)
}

*/
package jira
