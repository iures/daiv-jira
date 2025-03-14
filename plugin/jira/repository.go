package jira

import (
	"fmt"
	"strings"
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
	
	// For testing purposes
	getUserFunc func() (*User, error)
	searchIssuesFunc func(jql string, options *extJira.SearchOptions) ([]extJira.Issue, error)
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
	// If a mock function is provided for testing, use it
	if r.getUserFunc != nil {
		return r.getUserFunc()
	}
	
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
	rawIssues, err := r.fetchUpdatedIssues(pluginTimeRange, userID)
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

// fetchUpdatedIssues retrieves issues from Jira based on the given time range and user ID
func (r *JiraAPIRepository) fetchUpdatedIssues(timeRange plugin.TimeRange, userID string) ([]extJira.Issue, error) {
	// If a mock function is provided for testing, use it
	if r.searchIssuesFunc != nil {
		// Convert the plugin.TimeRange to string format for the JQL query
		fromTime := timeRange.Start.Format("2006-01-02 15:04")
		toTime := timeRange.End.Format("2006-01-02 15:04")
		
		// Build the JQL query
		jql := r.buildJQLQuery(fromTime, toTime)
		
		// Create search options
		options := &extJira.SearchOptions{
			MaxResults: r.config.QueryOptions.MaxResults,
			Fields:     r.config.QueryOptions.Fields,
		}
		
		// If changelog should be expanded, add it to the expand options
		if r.config.QueryOptions.ExpandChangelog {
			options.Expand = "changelog"
		}
		
		return r.searchIssuesFunc(jql, options)
	}

	// Format time range for JQL query
	fromTime := timeRange.Start.Format("2006-01-02 15:04")
	toTime := timeRange.End.Format("2006-01-02 15:04")

	// Build the JQL query
	jql := r.buildJQLQuery(fromTime, toTime)

	// Create search options
	options := &extJira.SearchOptions{
		MaxResults: r.config.QueryOptions.MaxResults,
		Fields:     r.config.QueryOptions.Fields,
	}

	// If changelog should be expanded, add it to the expand options
	if r.config.QueryOptions.ExpandChangelog {
		options.Expand = "changelog"
	}

	// Search for issues
	issues, _, err := r.client.Issue.Search(jql, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues in Jira: %w", err)
	}

	return issues, nil
}

// buildJQLQuery builds a JQL query based on the query options
func (r *JiraAPIRepository) buildJQLQuery(fromTime, toTime string) string {
	var conditions []string
	opts := r.config.QueryOptions

	// Start with the base JQL template
	baseQuery := fmt.Sprintf(opts.JQLTemplate, opts.Project, fromTime, toTime)
	conditions = append(conditions, baseQuery)

	// Add assignee condition if needed
	if opts.AssigneeCurrentUser {
		conditions = append(conditions, "assignee = currentUser()")
	}

	// Add status filter if provided
	if opts.StatusFilter != "" {
		// Handle special case for "!Closed" which is not valid JQL
		if opts.StatusFilter == "!Closed" {
			conditions = append(conditions, "status != Closed")
		} else if opts.StatusFilter == "!= Closed" {
			conditions = append(conditions, "status != Closed")
		} else {
			conditions = append(conditions, fmt.Sprintf("status %s", opts.StatusFilter))
		}
	}

	// Add sprint condition if needed
	if opts.InOpenSprints {
		conditions = append(conditions, "sprint IN openSprints()")
	}

	// Join all conditions with AND
	return strings.Join(conditions, " AND ")
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
