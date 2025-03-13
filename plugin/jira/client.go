package jira

import (
	"fmt"

	extJira "github.com/andygrunwald/go-jira"
	plugin "github.com/iures/daivplug"
)

type JiraConfig struct {
	Username string
	Token    string
	URL      string
	Project  string
}

type JiraClient struct {
	client *extJira.Client
	config *JiraConfig
	user   *extJira.User
}

func NewJiraClient(config *JiraConfig) (*JiraClient, error) {
	tp := extJira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Token,
	}

	client, err := extJira.NewClient(tp.Client(), config.URL)
	if err != nil {
		return nil, err
	}

	return &JiraClient{
		client: client,
		config: config,
	}, nil
}

func (j *JiraClient) GetSelf() (*extJira.User, error) {
	user, _, err := j.client.User.GetSelf()
	if err != nil {
		return nil, err
	}

	j.user = user

	return user, nil
}

func (j *JiraClient) fetchUpdatedIssues(timeRange plugin.TimeRange) ([]extJira.Issue, error) {
	fromTime := timeRange.Start.Format("2006-01-02")
	toTime := timeRange.End.Format("2006-01-02")

	searchString := fmt.Sprintf(
		`assignee = currentUser() AND project = %s AND status != Closed AND sprint IN openSprints() AND (updatedDate >= %s AND updatedDate < %s)`,
		j.config.Project,
		fromTime,
		toTime,
	)

	opt := &extJira.SearchOptions{
		MaxResults: 100,
		Expand:     "changelog",
		Fields:     []string{"summary", "description", "status", "changelog", "comment"},
	}

	issues, _, err := j.client.Issue.Search(searchString, opt)

	if err != nil {
		return nil, err
	}

	return issues, nil
}
