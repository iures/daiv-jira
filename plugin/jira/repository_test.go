package jira

import (
	"errors"
	"testing"
	"time"

	extJira "github.com/andygrunwald/go-jira"
)

// MockJiraClient is a mock implementation of the external Jira client
type MockJiraClient struct {
	MockGetSelf       func() (*extJira.User, error)
	MockSearchIssues  func(jql string, options *extJira.SearchOptions) ([]extJira.Issue, *extJira.Response, error)
}

// We'll use a simpler approach to mocking by directly patching the repository methods
// instead of trying to mock the go-jira client which has complex interfaces

func TestJiraAPIRepository_GetUser(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		setupMock   func(*JiraAPIRepository)
		expectError bool
		expectedUser *User
	}{
		{
			name: "Successful user retrieval",
			setupMock: func(repo *JiraAPIRepository) {
				// Replace the client's User.GetSelf method with our mock
				repo.client = &extJira.Client{}
				repo.getUserFunc = func() (*User, error) {
					return &User{
						AccountID:   "user123",
						DisplayName: "Test User",
						Email:       "test@example.com",
					}, nil
				}
			},
			expectError: false,
			expectedUser: &User{
				AccountID:   "user123",
				DisplayName: "Test User",
				Email:       "test@example.com",
			},
		},
		{
			name: "Error getting user",
			setupMock: func(repo *JiraAPIRepository) {
				repo.client = &extJira.Client{}
				repo.getUserFunc = func() (*User, error) {
					return nil, errors.New("failed to get user")
				}
			},
			expectError: true,
			expectedUser: nil,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create repository with default config
			config := &JiraConfig{
				Username: "test",
				Token:    "test",
				URL:      "https://test.atlassian.net",
				Project:  "TEST",
				QueryOptions: DefaultQueryOptions(),
			}
			repo := NewJiraAPIRepository(&extJira.Client{}, config)
			
			// Setup the mock
			tc.setupMock(repo)

			// Call the method being tested
			user, err := repo.GetUser()

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If no error is expected, check the user
			if !tc.expectError && err == nil {
				if user.AccountID != tc.expectedUser.AccountID {
					t.Errorf("Expected user account ID %s, got %s", tc.expectedUser.AccountID, user.AccountID)
				}
				if user.DisplayName != tc.expectedUser.DisplayName {
					t.Errorf("Expected user display name %s, got %s", tc.expectedUser.DisplayName, user.DisplayName)
				}
				if user.Email != tc.expectedUser.Email {
					t.Errorf("Expected user email %s, got %s", tc.expectedUser.Email, user.Email)
				}
			}
		})
	}
}

func TestJiraAPIRepository_GetIssues(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name          string
		setupMock     func(*JiraAPIRepository)
		timeRange     TimeRange
		userID        string
		expectError   bool
		expectedIssues int
	}{
		{
			name: "Successful issues retrieval",
			setupMock: func(repo *JiraAPIRepository) {
				repo.client = &extJira.Client{}
				repo.searchIssuesFunc = func(jql string, options *extJira.SearchOptions) ([]extJira.Issue, error) {
					return []extJira.Issue{
						{
							Key: "JIRA-123",
							Fields: &extJira.IssueFields{
								Summary: "Test Issue",
								Status: &extJira.Status{
									Name: "In Progress",
								},
							},
						},
					}, nil
				}
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			userID:        "user123",
			expectError:   false,
			expectedIssues: 1,
		},
		{
			name: "Error searching issues",
			setupMock: func(repo *JiraAPIRepository) {
				repo.client = &extJira.Client{}
				repo.searchIssuesFunc = func(jql string, options *extJira.SearchOptions) ([]extJira.Issue, error) {
					return nil, errors.New("failed to search issues")
				}
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			userID:        "user123",
			expectError:   true,
			expectedIssues: 0,
		},
		{
			name: "Empty issues list",
			setupMock: func(repo *JiraAPIRepository) {
				repo.client = &extJira.Client{}
				repo.searchIssuesFunc = func(jql string, options *extJira.SearchOptions) ([]extJira.Issue, error) {
					return []extJira.Issue{}, nil
				}
			},
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			userID:        "user123",
			expectError:   false,
			expectedIssues: 0,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create repository with default config
			config := &JiraConfig{
				Username: "test",
				Token:    "test",
				URL:      "https://test.atlassian.net",
				Project:  "TEST",
				QueryOptions: DefaultQueryOptions(),
			}
			repo := NewJiraAPIRepository(&extJira.Client{}, config)
			
			// Setup the mock
			tc.setupMock(repo)

			// Call the method being tested
			issues, err := repo.GetIssues(tc.timeRange, tc.userID)

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If no error is expected, check the issues
			if !tc.expectError && err == nil {
				if len(issues) != tc.expectedIssues {
					t.Errorf("Expected %d issues, got %d", tc.expectedIssues, len(issues))
				}

				// Check issue details if issues were returned
				if tc.expectedIssues > 0 {
					if issues[0].Key != "JIRA-123" {
						t.Errorf("Expected issue key JIRA-123, got %s", issues[0].Key)
					}
					if issues[0].Summary != "Test Issue" {
						t.Errorf("Expected issue summary 'Test Issue', got %s", issues[0].Summary)
					}
					if issues[0].Status != "In Progress" {
						t.Errorf("Expected issue status 'In Progress', got %s", issues[0].Status)
					}
				}
			}
		})
	}
} 
