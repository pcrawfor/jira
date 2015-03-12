package jira

import (
	"encoding/json"
	"fmt"
)

// IssueService is a jira client with functions for operating on issues
type IssueService struct {
	client *Jira
}

// Issue is the type representing of a Jira Issue
type Issue struct {
	ID     string                 `json:"id,omitempty"`
	Key    string                 `json:"key,omitempty"`
	Self   string                 `json:"self,omitempty"`
	Expand string                 `json:"expand,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

// IssueList is the type representing of a list of Jira Issues as defined by their API response structure
type IssueList struct {
	Expand     string   `json:"expand,omitempty"`
	StartAt    int      `json:"starts_at,omitempty"`
	MaxResults int      `json:"max_results,omitempty"`
	Total      int      `json:"total,omitempty"`
	Issues     []*Issue `json:"issues,omitempty"`
	//Pagination *Pagination
}

// TransitionList is the type representing of a list of Jira transitions as defined by their API response structure
type TransitionList struct {
	Expand      string       `json:"expand,omitempty"`
	Transitions []Transition `json:"transitions,omitempty"`
}

// Transition is the type representing of a Jira Transition
type Transition struct {
	ID     string                 `json:"id,omitempty"`
	Name   string                 `json:"name,omitempty"`
	To     map[string]interface{} `json:"to,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

var issueBasePath = restPath + "issue/"

// Transition executes a transition for the given issue key to the given transition ID or returns an error
func (i *IssueService) Transition(key, transitionID string) ([]byte, error) {
	url := issueBasePath + key + "/transitions"

	c := map[string]interface{}{
		"add": map[string]string{
			"body": "releasebot transition",
		},
	}

	params := map[string]interface{}{
		"transition": map[string]string{
			"id": transitionID,
		},
		"update": map[string]interface{}{
			"comment": []interface{}{c},
		},
	}

	return i.client.execRequest(MPOST, i.client.baseurl+url, params)
}

// GetTransitions loads the available transitions for a given issue key
func (i *IssueService) GetTransitions(key string) (*TransitionList, error) {
	url := "issue/" + key + "/transitions?expand=transitions.fields"
	b, e := i.client.apiRequest(MGET, url, nil)
	if e != nil {
		return nil, e
	}

	transitions := TransitionList{}
	terr := json.Unmarshal(b, &transitions)
	if terr != nil {
		fmt.Println("Transitions error: ", terr)
		return nil, terr
	}

	return &transitions, nil
}
