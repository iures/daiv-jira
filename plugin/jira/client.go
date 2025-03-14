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
	QueryOptions QueryOptions
}

// JiraClient provides a client for interacting with Jira
type JiraClient struct {
	client     *extJira.Client
	config     *JiraConfig
	repository JiraRepository
}

// NewJiraClient creates a new JiraClient
func NewJiraClient(config *JiraConfig) (*JiraClient, error) {
	tp := extJira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Token,
	}

	client, err := extJira.NewClient(tp.Client(), config.URL)
	if err != nil {
		return nil, err
	}

	// Set project in query options if not already set
	if config.QueryOptions.Project == "" {
		config.QueryOptions.Project = config.Project
	}

	jiraClient := &JiraClient{
		client: client,
		config: config,
	}

	// Create the repository
	repository := NewJiraAPIRepository(client, config)
	jiraClient.repository = repository

	return jiraClient, nil
}

// GetRepository returns the Jira repository
func (j *JiraClient) GetRepository() JiraRepository {
	return j.repository
}

func (j *JiraClient) GetSelf() (*extJira.User, error) {
	user, _, err := j.client.User.GetSelf()
	if err != nil {
		return nil, err
	}

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
