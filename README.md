# Jira

Simple Go convenience wrapper for Jira API functionality

To use it import: _github.com/pcrawfor/jira_

## Example Usage:

    // create a client
    client := jira.NewJiraClient("https://yoursubdomain.atlassian.net", your_username, your_password, 1500)

    // load issues matching the given id's
    issues, err := client.Issues([]string{"WEB-1", "DEV-4"}, nil)
    if err != nil {
        fmt.Println("Error loading issues: ", err)
    }

    // the issues slice contains Id, Key, Self, Expand and Fields attributes
    for i, _ := range issues {
        fmt.Println("Issue Id: ", i.Id)
    }

## Fields

The fields included by default are 'id' and 'summary' - you can provide a custom set of fields by specifying them on the Issue, Issues and SearchWithFields calls

