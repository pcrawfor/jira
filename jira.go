package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const restPath = "/rest/api/2/"
const defaultMaxResults = 200

const mPost = "POST"
const mGet = "GET"

// Jira is a client object with functions to make reuqests to the jira api
type Jira struct {
	client        *http.Client
	baseurl       string
	auth          Auth
	maxResults    int
	IssuesService *IssueService
}

// Auth contains username and password attributes used for api request authentication
type Auth struct {
	Username string
	Password string
}

func (i *Issue) String() string {
	return "Id: " + i.ID + " Key: " + i.Key + " self: " + i.Self
}

// NewJiraClient returns an instance of the Jira api client
func NewJiraClient(baseurl, username, password string, maxResults int) *Jira {
	if maxResults == -1 {
		maxResults = defaultMaxResults
	}
	j := &Jira{client: &http.Client{}, baseurl: baseurl, auth: Auth{username, password}, maxResults: maxResults}
	j.IssuesService = &IssueService{j}

	return j
}

// "search?jql=status=reviewed OR status=released OR status='ready for release' OR status='qa review'&validateQuery=true&fields=id,summary"

// Search runs an arbitrary search request against the Jira API for Issues
func (j *Jira) Search(query string) ([]*Issue, error) {
	return j.SearchWithFields(query, nil)
}

// SearchWithFields runs an arbitrary search request and builds the set of fields to be returned by the response as defined in the fields param
func (j *Jira) SearchWithFields(query string, fields []string) ([]*Issue, error) {
	max := strconv.Itoa(j.maxResults)

	useFields := "id,summary"

	if nil != fields && len(fields) > 0 {
		useFields = flatten(fields)
	}

	params := map[string]string{
		"jql":           query,
		"validateQuery": "true",
		"fields":        useFields,
		"maxResults":    max,
	}

	urlStr := j.buildURL("search", params)

	issueData, err := j.execRequest(mGet, urlStr, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	// parse the issue data and return the issues
	issueList := IssueList{}
	uerr := json.Unmarshal(issueData, &issueList)
	if uerr != nil {
		fmt.Println("Issue search error: ", uerr)
		return nil, uerr
	}

	return issueList.Issues, nil
}

// Issue loads the jira data for a single jira issue key, with the specified issue fields if the fields param is set
func (j *Jira) Issue(key string, fields []string) (*Issue, error) {
	useFields := "id,summary"

	if nil != fields && len(fields) > 0 {
		useFields = flatten(fields)
	}

	params := map[string]string{
		"fields": useFields,
	}

	urlStr := j.buildURL("issue/"+key, params)
	issueData, err := j.execRequest(mGet, urlStr, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	issue := Issue{}
	uerr := json.Unmarshal(issueData, &issue)
	if uerr != nil {
		fmt.Println("Issue error: ", uerr)
		return nil, uerr
	}

	return &issue, nil
}

// Issues loads the jira data for all the issue keys provided specifying the fields to include if the fields param is set
func (j *Jira) Issues(keys []string, fields []string) ([]*Issue, error) {
	// build a query with all the issue keys
	qry := ""
	fmt.Println("keys: ", keys)
	for i := 0; i < len(keys); i++ {
		if i == len(keys)-1 {
			qry = qry + fmt.Sprintf("id = %s", keys[i])
		} else {
			qry = qry + fmt.Sprintf("id = %s or ", keys[i])
		}
	}

	fmt.Println("QRY: ", qry)

	return j.SearchWithFields(qry, fields)
}

// apiRequest builds a request for the jira API
func (j *Jira) apiRequest(method, path string, params map[string]interface{}) ([]byte, error) {
	url := j.baseurl + restPath + path
	return j.execRequest(method, url, nil)
}

//buildURL creates a url for the given path and url parameters
func (j *Jira) buildURL(path string, params map[string]string) string {
	var aURL *url.URL
	aURL, err := url.Parse(j.baseurl)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	aURL.Path += restPath + path
	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	aURL.RawQuery = parameters.Encode()
	return aURL.String()
}

// execRequest executes an arbitrary request for the given method and url returning the contents of the response in []byte or an error
func (j *Jira) execRequest(method, aURL string, params map[string]interface{}) ([]byte, error) {

	// json string encode the params for the POST body if there are any
	var body io.Reader
	if params != nil && method == mPost {
		b, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Json error: ", err)
		}
		body = bytes.NewBuffer(b)
		fmt.Println("BODY: ", string(b))
	}

	req, err := http.NewRequest(method, aURL, body)
	if err != nil {
		fmt.Println("execRequest error: ", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(j.auth.Username, j.auth.Password)

	fmt.Println("URL: ", req.URL)

	resp, rerr := j.client.Do(req)
	if rerr != nil {
		fmt.Println("req error: ", rerr)
		return nil, rerr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode > 399 {
		return nil, fmt.Errorf("HTTP Error Status returned: %d", resp.StatusCode)
	}

	data, derr := ioutil.ReadAll(resp.Body)
	if derr != nil {
		fmt.Println("Error reading response: ", derr)
		return nil, derr
	}

	return data, nil
}

func flatten(list []string) string {
	str := ""
	for _, v := range list {
		str = str + "," + v
	}
	return str
}
