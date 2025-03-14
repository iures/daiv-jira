package jira

import (
	"fmt"
	"sync"
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
	if len(issues) == 0 {
		return []Issue{}
	}

	// Create a channel to receive processed issues
	resultChan := make(chan Issue, len(issues))
	
	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	
	// Process each issue in a separate goroutine
	for _, issue := range issues {
		wg.Add(1)
		go func(issue extJira.Issue) {
			defer wg.Done()
			
			// Process the issue
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
			
			// Send the processed issue to the channel
			resultChan <- domainIssue
		}(issue)
	}
	
	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results from the channel
	result := make([]Issue, 0, len(issues))
	for domainIssue := range resultChan {
		result = append(result, domainIssue)
	}
	
	return result
}

// processComments converts external Jira comments to domain model comments
func (s *ActivityService) processComments(comments []*extJira.Comment, timeRange TimeRange) []Comment {
	if len(comments) == 0 {
		return []Comment{}
	}
	
	// For small number of comments, process sequentially
	if len(comments) < 5 {
		return s.processCommentsSequential(comments, timeRange)
	}
	
	// Create a channel to receive processed comments
	resultChan := make(chan Comment, len(comments))
	
	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	
	// Process each comment in a separate goroutine
	for _, comment := range comments {
		wg.Add(1)
		go func(comment *extJira.Comment) {
			defer wg.Done()
			
			createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", comment.Created)
			if err != nil {
				return
			}
			
			if timeRange.IsInRange(createdTime) {
				resultChan <- Comment{
					Timestamp: createdTime,
					Author:    comment.Author.DisplayName,
					Content:   comment.Body,
				}
			}
		}(comment)
	}
	
	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results from the channel
	result := make([]Comment, 0, len(comments))
	for comment := range resultChan {
		result = append(result, comment)
	}
	
	return result
}

// processCommentsSequential processes comments sequentially (for small number of comments)
func (s *ActivityService) processCommentsSequential(comments []*extJira.Comment, timeRange TimeRange) []Comment {
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
	if len(histories) == 0 {
		return []Change{}
	}
	
	// For small number of histories, process sequentially
	if len(histories) < 5 {
		return s.processChangelogSequential(histories, timeRange, userAccountID)
	}
	
	// Create a channel to receive processed changes
	resultChan := make(chan Change, len(histories)*2) // Assuming average of 2 changes per history
	
	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	
	// Process each history in a separate goroutine
	for _, history := range histories {
		wg.Add(1)
		go func(history extJira.ChangelogHistory) {
			defer wg.Done()
			
			createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
			if err != nil {
				return
			}
			
			if timeRange.IsInRange(createdTime) && history.Author.AccountID == userAccountID {
				for _, item := range history.Items {
					resultChan <- Change{
						Timestamp: createdTime,
						Author:    history.Author.DisplayName,
						Field:     item.Field,
						FromValue: item.FromString,
						ToValue:   item.ToString,
					}
				}
			}
		}(history)
	}
	
	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results from the channel
	result := make([]Change, 0)
	for change := range resultChan {
		result = append(result, change)
	}
	
	return result
}

// processChangelogSequential processes changelog sequentially (for small number of histories)
func (s *ActivityService) processChangelogSequential(histories []extJira.ChangelogHistory, timeRange TimeRange, userAccountID string) []Change {
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
