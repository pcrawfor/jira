package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildUrl(t *testing.T) {
	c := NewJiraClient("http://test.jira.com", "foo", "bar", 100)
	params := map[string]string{
		"a": "b",
	}
	url := c.buildURL("test", params)
	if url != "http://test.jira.com/rest/api/2/test?a=b" {
		t.Error("Generated URL is invalid expected http://test.jira.com/rest/api/2/test?a=b got:", url)
	}
}

// TestIssue verifies that we can handle a expected response from the Jira API - all bets are off if they change what they send us though :)
func TestIssue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := `{"expand":"renderedFields,names,schema,transitions,operations,editmeta,changelog","id":"1234","self":"https://test.atlassian.net/rest/api/2/issue/1234","key":"ABC-01","fields":{"summary":"This is a test"}}`
		fmt.Fprintln(w, val)
	}))
	defer ts.Close()

	fmt.Println("ts.URL:", ts.URL)
	c := NewJiraClient(ts.URL, "foo", "bar", 100)
	i, err := c.Issue("abc", nil)
	if err != nil {
		t.Error("Error loading issue:", err)
	}

	if i == nil {
		t.Error("issue is nil")
	}

	if i.ID != "1234" {
		t.Error("Error expected issue ID to be 1234 got:", i.ID)
	}

	if i.Key != "ABC-01" {
		t.Error("Error expected issue Key to be ABC-01 got:", i.Key)
	}
}