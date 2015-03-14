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
		checkAuth(t, r)
		val := `{"expand":"renderedFields,names,schema,transitions,operations,editmeta,changelog","id":"1234","self":"https://test.atlassian.net/rest/api/2/issue/1234","key":"ABC-01","fields":{"summary":"This is a test"}}`
		fmt.Fprintln(w, val)
	}))
	defer ts.Close()

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

func TestIssues(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkAuth(t, r)
		val := `{"expand":"schema,names","startAt":0,"maxResults":100,"total":3,"issues":[
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1234","self":"https://test.atlassian.net/rest/api/2/issue/1234","key":"ABC-01","fields":{"summary":"This is a test"}},
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1235","self":"https://test.atlassian.net/rest/api/2/issue/1235","key":"ABC-02","fields":{"summary":"This is another test"}},
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1236","self":"https://test.atlassian.net/rest/api/2/issue/1236","key":"ABC-03","fields":{"summary":"This is also test"}}]}`
		fmt.Fprintln(w, val)
	}))
	defer ts.Close()

	c := NewJiraClient(ts.URL, "foo", "bar", 100)
	issues, err := c.Issues([]string{"ABC-01", "ABC-02", "ABC-03"}, nil)
	if err != nil {
		t.Error("Error loading issue:", err)
	}

	if len(issues) != 3 {
		t.Error("Recv'd the wrong number of issues in return")
	}

	if issues[0].Key != "ABC-01" {
		t.Error("Error expected first key to be ABC-01 got:", issues[0].Key)
	}

}

func TestSearch(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkAuth(t, r)
		val := `{"expand":"schema,names","startAt":0,"maxResults":100,"total":3,"issues":[
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1234","self":"https://test.atlassian.net/rest/api/2/issue/1234","key":"ABC-01","fields":{"summary":"This is a test"}},
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1235","self":"https://test.atlassian.net/rest/api/2/issue/1235","key":"ABC-02","fields":{"summary":"This is another test"}},
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1236","self":"https://test.atlassian.net/rest/api/2/issue/1236","key":"ABC-03","fields":{"summary":"This is also test"}},
		{"expand":"operations,editmeta,changelog,transitions,renderedFields","id":"1237","self":"https://test.atlassian.net/rest/api/2/issue/1237","key":"ABC-04","fields":{"summary":"This is more tests"}}]}`
		fmt.Fprintln(w, val)
	}))
	defer ts.Close()

	c := NewJiraClient(ts.URL, "foo", "bar", 100)
	issues, err := c.Search("status=reviewed")
	if err != nil {
		t.Error("Error loading issue:", err)
	}

	if len(issues) != 4 {
		t.Error("Recv'd the wrong number of issues in return")
	}
}

func checkAuth(t *testing.T, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok {
		t.Error("Error loading basic auth")
	}
	if u != "foo" || p != "bar" {
		t.Error("Auth creds invalid")
	}
}
