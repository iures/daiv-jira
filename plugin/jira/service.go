package jira

import (
	"fmt"
	"time"

	extJira "github.com/andygrunwald/go-jira"
	plugin "github.com/iures/daivplug"
)

// ActivityService handles the processing of Jira data into domain models
type ActivityService struct {
	repository JiraRepository
}

// NewActivityService creates a new activity service
func NewActivityService(repository JiraRepository) *ActivityService {
	return &ActivityService{
		repository: repository,
	}
}

// GetActivityReport retrieves and processes Jira activity data for the given time range
func (s *ActivityService) GetActivityReport(pluginTimeRange plugin.TimeRange) (*ActivityReport, error) {
	// Convert plugin.TimeRange to our domain TimeRange
	timeRange := TimeRange{
		Start: pluginTimeRange.Start,
		End:   pluginTimeRange.End,
	}

	// Get the current user
	user, err := s.repository.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get issues for the user and time range
	issues, err := s.repository.GetIssues(timeRange, user.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	// Create and return the activity report
	return &ActivityReport{
		TimeRange: timeRange,
		User:      *user,
		Issues:    issues,
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
