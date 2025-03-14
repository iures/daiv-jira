package jira

import (
	"errors"
	"fmt"
	"testing"
	"time"

	extJira "github.com/andygrunwald/go-jira"
)

// MockJiraRepository is a mock implementation of JiraRepository for testing
type MockJiraRepository struct {
	MockGetUser   func() (*User, error)
	MockGetIssues func(timeRange TimeRange, userAccountID string) ([]Issue, error)
}

// GetUser implements the JiraRepository interface
func (m *MockJiraRepository) GetUser() (*User, error) {
	return m.MockGetUser()
}

// GetIssues implements the JiraRepository interface
func (m *MockJiraRepository) GetIssues(timeRange TimeRange, userAccountID string) ([]Issue, error) {
	return m.MockGetIssues(timeRange, userAccountID)
}

func TestActivityService_GetActivityReport(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name          string
		mockRepo      *MockJiraRepository
		timeRange     TimeRange
		expectError   bool
		expectedIssues int
	}{
		{
			name: "Successful report generation",
			mockRepo: &MockJiraRepository{
				MockGetUser: func() (*User, error) {
					return &User{
						AccountID:   "user123",
						DisplayName: "Test User",
						Email:       "test@example.com",
					}, nil
				},
				MockGetIssues: func(timeRange TimeRange, userAccountID string) ([]Issue, error) {
					return []Issue{
						{
							Key:     "JIRA-123",
							Summary: "Test Issue",
							Status:  "In Progress",
						},
					}, nil
				},
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expectError:   false,
			expectedIssues: 1,
		},
		{
			name: "Error getting user",
			mockRepo: &MockJiraRepository{
				MockGetUser: func() (*User, error) {
					return nil, errors.New("failed to get user")
				},
				MockGetIssues: func(timeRange TimeRange, userAccountID string) ([]Issue, error) {
					return []Issue{}, nil
				},
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expectError:   true,
			expectedIssues: 0,
		},
		{
			name: "Error getting issues",
			mockRepo: &MockJiraRepository{
				MockGetUser: func() (*User, error) {
					return &User{
						AccountID:   "user123",
						DisplayName: "Test User",
						Email:       "test@example.com",
					}, nil
				},
				MockGetIssues: func(timeRange TimeRange, userAccountID string) ([]Issue, error) {
					return nil, errors.New("failed to get issues")
				},
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expectError:   true,
			expectedIssues: 0,
		},
		{
			name: "Empty issues list",
			mockRepo: &MockJiraRepository{
				MockGetUser: func() (*User, error) {
					return &User{
						AccountID:   "user123",
						DisplayName: "Test User",
						Email:       "test@example.com",
					}, nil
				},
				MockGetIssues: func(timeRange TimeRange, userAccountID string) ([]Issue, error) {
					return []Issue{}, nil
				},
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expectError:   false,
			expectedIssues: 0,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with mock repository
			service := NewActivityService(tc.mockRepo)

			// Convert TimeRange to plugin.TimeRange
			pluginTimeRange := struct {
				Start time.Time
				End   time.Time
			}{
				Start: tc.timeRange.Start,
				End:   tc.timeRange.End,
			}

			// Call the method being tested
			report, err := service.GetActivityReport(pluginTimeRange)

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If no error is expected, check the report
			if !tc.expectError && err == nil {
				// Check time range
				if !report.TimeRange.Start.Equal(tc.timeRange.Start) {
					t.Errorf("Expected start time %v, got %v", tc.timeRange.Start, report.TimeRange.Start)
				}
				if !report.TimeRange.End.Equal(tc.timeRange.End) {
					t.Errorf("Expected end time %v, got %v", tc.timeRange.End, report.TimeRange.End)
				}

				// Check issues count
				if len(report.Issues) != tc.expectedIssues {
					t.Errorf("Expected %d issues, got %d", tc.expectedIssues, len(report.Issues))
				}

				// Check user info if issues were returned
				if tc.expectedIssues > 0 {
					expectedUser, _ := tc.mockRepo.GetUser()
					if report.User.AccountID != expectedUser.AccountID {
						t.Errorf("Expected user account ID %s, got %s", expectedUser.AccountID, report.User.AccountID)
					}
					if report.User.DisplayName != expectedUser.DisplayName {
						t.Errorf("Expected user display name %s, got %s", expectedUser.DisplayName, report.User.DisplayName)
					}
					if report.User.Email != expectedUser.Email {
						t.Errorf("Expected user email %s, got %s", expectedUser.Email, report.User.Email)
					}
				}
			}
		})
	}
}

func TestActivityService_ProcessIssuesConcurrently(t *testing.T) {
	// Create a large number of test issues to demonstrate concurrency benefits
	const numIssues = 50
	testIssues := make([]extJira.Issue, numIssues)
	
	for i := 0; i < numIssues; i++ {
		testIssues[i] = extJira.Issue{
			Key: fmt.Sprintf("JIRA-%d", i+1),
			Fields: &extJira.IssueFields{
				Summary: fmt.Sprintf("Test Issue %d", i+1),
				Status: &extJira.Status{
					Name: "In Progress",
				},
				Comments: &extJira.Comments{
					Comments: []*extJira.Comment{
						{
							Created: "2023-01-01T12:00:00.000-0700",
							Author: extJira.User{
								DisplayName: "Test User",
							},
							Body: fmt.Sprintf("Comment for issue %d", i+1),
						},
					},
				},
			},
			Changelog: &extJira.Changelog{
				Histories: []extJira.ChangelogHistory{
					{
						Created: "2023-01-01T10:00:00.000-0700",
						Author: extJira.User{
							AccountID:   "user123",
							DisplayName: "Test User",
						},
						Items: []extJira.ChangelogItems{
							{
								Field:      "status",
								FromString: "Open",
								ToString:   "In Progress",
							},
						},
					},
				},
			},
		}
	}
	
	// Create a mock repository
	mockRepo := &MockJiraRepository{}
	
	// Create the service
	service := NewActivityService(mockRepo)
	
	// Set up the test time range and user
	timeRange := TimeRange{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	user := User{
		AccountID:   "user123",
		DisplayName: "Test User",
		Email:       "test@example.com",
	}
	
	// Measure the time it takes to process issues concurrently
	startConcurrent := time.Now()
	resultConcurrent := service.processIssues(testIssues, timeRange, user)
	durationConcurrent := time.Since(startConcurrent)
	
	// Verify the results
	if len(resultConcurrent) != numIssues {
		t.Errorf("Expected %d issues, got %d", numIssues, len(resultConcurrent))
	}
	
	// Check a few issues to ensure they were processed correctly
	for i := 0; i < numIssues; i++ {
		found := false
		expectedKey := fmt.Sprintf("JIRA-%d", i+1)
		
		for _, issue := range resultConcurrent {
			if issue.Key == expectedKey {
				found = true
				
				// Check that comments were processed
				if len(issue.Comments) != 1 {
					t.Errorf("Expected 1 comment for issue %s, got %d", issue.Key, len(issue.Comments))
				}
				
				// Check that changes were processed
				if len(issue.Changes) != 1 {
					t.Errorf("Expected 1 change for issue %s, got %d", issue.Key, len(issue.Changes))
				}
				
				break
			}
		}
		
		if !found {
			t.Errorf("Issue with key %s not found in results", expectedKey)
		}
	}
	
	// Log the processing time for information
	t.Logf("Processed %d issues concurrently in %v", numIssues, durationConcurrent)
} 
