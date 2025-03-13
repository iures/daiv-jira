package jira

import (
	"time"

	extJira "github.com/andygrunwald/go-jira"
	plugin "github.com/iures/daivplug"
)

// ActivityService handles the processing of Jira data into domain models
type ActivityService struct {
	client *JiraClient
}

// NewActivityService creates a new activity service
func NewActivityService(client *JiraClient) *ActivityService {
	return &ActivityService{
		client: client,
	}
}

// GetActivityReport retrieves and processes Jira activity data for the given time range
func (s *ActivityService) GetActivityReport(pluginTimeRange plugin.TimeRange, user *extJira.User) (*ActivityReport, error) {
	// Convert plugin.TimeRange to our domain TimeRange
	timeRange := TimeRange{
		Start: pluginTimeRange.Start,
		End:   pluginTimeRange.End,
	}

	// Convert external User to our domain User
	domainUser := User{
		AccountID:   user.AccountID,
		DisplayName: user.DisplayName,
		Email:       user.EmailAddress,
	}

	// Fetch issues from Jira
	issues, err := s.client.fetchUpdatedIssues(pluginTimeRange)
	if err != nil {
		return nil, err
	}

	// Process issues into domain model
	domainIssues := s.processIssues(issues, timeRange, domainUser)

	// Create and return the activity report
	return &ActivityReport{
		TimeRange: timeRange,
		User:      domainUser,
		Issues:    domainIssues,
	}, nil
}

// processIssues converts external Jira issues to domain model issues
func (s *ActivityService) processIssues(issues []extJira.Issue, timeRange TimeRange, user User) []Issue {
	result := make([]Issue, 0, len(issues))

	for _, issue := range issues {
		domainIssue := Issue{
			Key:     issue.Key,
			Summary: issue.Fields.Summary,
			Status:  issue.Fields.Status.Name,
		}

		// Process comments
		if issue.Fields.Comments != nil {
			domainIssue.Comments = s.processComments(issue.Fields.Comments.Comments, timeRange)
		}

		// Process changelog
		if issue.Changelog != nil {
			domainIssue.Changes = s.processChangelog(issue.Changelog.Histories, timeRange, user.AccountID)
		}

		result = append(result, domainIssue)
	}

	return result
}

// processComments converts external Jira comments to domain model comments
func (s *ActivityService) processComments(comments []*extJira.Comment, timeRange TimeRange) []Comment {
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
func (s *ActivityService) processChangelog(histories []extJira.ChangelogHistory, timeRange TimeRange, userAccountID string) []Change {
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
