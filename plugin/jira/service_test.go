package jira

import (
	"errors"
	"testing"
	"time"
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
