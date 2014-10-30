package jira

import (
	"encoding/json"
	"fmt"
)

type IssueService struct {
	client *Jira
}

type Issue struct {
	Id     string                 `json:"id,omitempty"`
	Key    string                 `json:"key,omitempty"`
	Self   string                 `json:"self,omitempty"`
	Expand string                 `json:"expand,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

type IssueList struct {
	Expand     string   `json:"expand,omitempty"`
	StartAt    int      `json:"starts_at,omitempty"`
	MaxResults int      `json:"max_results,omitempty"`
	Total      int      `json:"total,omitempty"`
	Issues     []*Issue `json:"issues,omitempty"`
	//Pagination *Pagination
}

type TransitionList struct {
	Expand      string       `json:"expand,omitempty"`
	Transitions []Transition `json:"transitions,omitempty"`
}

type Transition struct {
	Id     string                 `json:"id,omitempty"`
	Name   string                 `json:"name,omitempty"`
	To     map[string]interface{} `json:"to,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

var issueBasePath = REST_PATH + "issue/"

func (i *IssueService) Transition(key, transitionId string) ([]byte, error) {
	url := issueBasePath + key + "/transitions"

	c := map[string]interface{}{
		"add": map[string]string{
			"body": "releasebot transition",
		},
	}

	params := map[string]interface{}{
		"transition": map[string]string{
			"id": transitionId,
		},
		"update": map[string]interface{}{
			"comment": []interface{}{c},
		},
	}

	return i.client.execRequest(MPOST, i.client.baseurl+url, params)
}

func (i *IssueService) GetTransitions(key string) (*TransitionList, error) {
	url := "issue/" + key + "/transitions?expand=transitions.fields"
	b, e := i.client.ApiRequest(MGET, url, nil)
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
