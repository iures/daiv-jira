package jira

import (
	"fmt"
	"time"

	extJira "github.com/andygrunwald/go-jira"
	plugin "github.com/iures/daivplug"
)

// JiraRepository defines the interface for accessing Jira data
type JiraRepository interface {
	GetUser() (*User, error)
	GetIssues(timeRange TimeRange, userID string) ([]Issue, error)
}

// JiraAPIRepository implements JiraRepository using the Jira API
type JiraAPIRepository struct {
	client *extJira.Client
	config *JiraConfig
}

// NewJiraAPIRepository creates a new JiraAPIRepository
func NewJiraAPIRepository(client *extJira.Client, config *JiraConfig) *JiraAPIRepository {
	return &JiraAPIRepository{
		client: client,
		config: config,
	}
}

// GetUser retrieves the current user from Jira
func (r *JiraAPIRepository) GetUser() (*User, error) {
	user, _, err := r.client.User.GetSelf()
	if err != nil {
		return nil, fmt.Errorf("failed to get user from Jira: %w", err)
	}

	return &User{
		AccountID:   user.AccountID,
		DisplayName: user.DisplayName,
		Email:       user.EmailAddress,
	}, nil
}

// GetIssues retrieves issues from Jira based on the given time range and user ID
func (r *JiraAPIRepository) GetIssues(timeRange TimeRange, userID string) ([]Issue, error) {
	// Convert domain TimeRange to plugin.TimeRange for the API call
	pluginTimeRange := plugin.TimeRange{
		Start: timeRange.Start,
		End:   timeRange.End,
	}

	// Fetch raw issues from Jira
	rawIssues, err := r.fetchUpdatedIssues(pluginTimeRange)
	if err != nil {
		return nil, err
	}

	// Convert raw issues to domain model
	issues := make([]Issue, 0, len(rawIssues))
	for _, rawIssue := range rawIssues {
		issue := Issue{
			Key:     rawIssue.Key,
			Summary: rawIssue.Fields.Summary,
			Status:  rawIssue.Fields.Status.Name,
		}

		// Process comments
		if rawIssue.Fields.Comments != nil {
			issue.Comments = r.processComments(rawIssue.Fields.Comments.Comments, timeRange)
		}

		// Process changelog
		if rawIssue.Changelog != nil {
			issue.Changes = r.processChangelog(rawIssue.Changelog.Histories, timeRange, userID)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// fetchUpdatedIssues retrieves issues from Jira based on the given time range
func (r *JiraAPIRepository) fetchUpdatedIssues(timeRange plugin.TimeRange) ([]extJira.Issue, error) {
	fromTime := timeRange.Start.Format("2006-01-02")
	toTime := timeRange.End.Format("2006-01-02")

	searchString := fmt.Sprintf(
		`assignee = currentUser() AND project = %s AND status != Closed AND sprint IN openSprints() AND (updatedDate >= %s AND updatedDate < %s)`,
		r.config.Project,
		fromTime,
		toTime,
	)

	opt := &extJira.SearchOptions{
		MaxResults: 100,
		Expand:     "changelog",
		Fields:     []string{"summary", "description", "status", "changelog", "comment"},
	}

	issues, _, err := r.client.Issue.Search(searchString, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues in Jira: %w", err)
	}

	return issues, nil
}

// processComments converts external Jira comments to domain model comments
func (r *JiraAPIRepository) processComments(comments []*extJira.Comment, timeRange TimeRange) []Comment {
	result := make([]Comment, 0)

	for _, comment := range comments {
		createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", comment.Created)
		if err != nil {
			continue
		}

		if timeRange.IsInRange(createdTime) {
			result = append(result, Comment{
				Timestamp: createdTime,
				Author:    comment.Author.DisplayName,
				Content:   comment.Body,
			})
		}
	}

	return result
}

// processChangelog converts external Jira changelog to domain model changes
func (r *JiraAPIRepository) processChangelog(histories []extJira.ChangelogHistory, timeRange TimeRange, userAccountID string) []Change {
	result := make([]Change, 0)

	for _, history := range histories {
		createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
		if err != nil {
			continue
		}

		if timeRange.IsInRange(createdTime) && history.Author.AccountID == userAccountID {
			for _, item := range history.Items {
				result = append(result, Change{
					Timestamp: createdTime,
					Author:    history.Author.DisplayName,
					Field:     item.Field,
					FromValue: item.FromString,
					ToValue:   item.ToString,
				})
			}
		}
	}

	return result
} 
