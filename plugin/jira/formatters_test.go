package jira

import (
	"strings"
	"testing"
	"time"
)

func TestXMLFormatter_Format(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		report      *ActivityReport
		expectedStr string
	}{
		{
			name: "Empty report",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{},
			},
			expectedStr: "<jira_report></jira_report>",
		},
		{
			name: "Report with issues",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{
					{
						Key:     "JIRA-123",
						Summary: "Test Issue",
						Status:  "In Progress",
						Changes: []Change{
							{
								Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Field:     "status",
								FromValue: "Open",
								ToValue:   "In Progress",
							},
						},
						Comments: []Comment{
							{
								Timestamp: time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Content:   "This is a test comment",
							},
						},
					},
				},
			},
			expectedStr: "JIRA-123",
		},
	}

	// Run tests
	formatter := NewXMLFormatter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formatter.Format(tc.report)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			if result.ContentType != "application/xml" {
				t.Errorf("Expected content type 'application/xml', got '%s'", result.ContentType)
			}
			
			if !strings.Contains(result.Content, tc.expectedStr) {
				t.Errorf("Expected content to contain '%s', got '%s'", tc.expectedStr, result.Content)
			}
		})
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		report      *ActivityReport
		expectedStr string
	}{
		{
			name: "Empty report",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{},
			},
			expectedStr: "{}",
		},
		{
			name: "Report with issues",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{
					{
						Key:     "JIRA-123",
						Summary: "Test Issue",
						Status:  "In Progress",
						Changes: []Change{
							{
								Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Field:     "status",
								FromValue: "Open",
								ToValue:   "In Progress",
							},
						},
						Comments: []Comment{
							{
								Timestamp: time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Content:   "This is a test comment",
							},
						},
					},
				},
			},
			expectedStr: "JIRA-123",
		},
	}

	// Run tests
	formatter := NewJSONFormatter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formatter.Format(tc.report)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			if result.ContentType != "application/json" {
				t.Errorf("Expected content type 'application/json', got '%s'", result.ContentType)
			}
			
			if !strings.Contains(result.Content, tc.expectedStr) {
				t.Errorf("Expected content to contain '%s', got '%s'", tc.expectedStr, result.Content)
			}
		})
	}
}

func TestMarkdownFormatter_Format(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		report      *ActivityReport
		expectedStr string
	}{
		{
			name: "Empty report",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{},
			},
			expectedStr: "No activity found",
		},
		{
			name: "Report with issues",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{
					{
						Key:     "JIRA-123",
						Summary: "Test Issue",
						Status:  "In Progress",
						Changes: []Change{
							{
								Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Field:     "status",
								FromValue: "Open",
								ToValue:   "In Progress",
							},
						},
						Comments: []Comment{
							{
								Timestamp: time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Content:   "This is a test comment",
							},
						},
					},
				},
			},
			expectedStr: "JIRA-123",
		},
	}

	// Run tests
	formatter := NewMarkdownFormatter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formatter.Format(tc.report)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			if result.ContentType != "text/markdown" {
				t.Errorf("Expected content type 'text/markdown', got '%s'", result.ContentType)
			}
			
			if !strings.Contains(result.Content, tc.expectedStr) {
				t.Errorf("Expected content to contain '%s', got '%s'", tc.expectedStr, result.Content)
			}
		})
	}
}

func TestHTMLFormatter_Format(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		report      *ActivityReport
		expectedStr string
	}{
		{
			name: "Empty report",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{},
			},
			expectedStr: "No activity found",
		},
		{
			name: "Report with issues",
			report: &ActivityReport{
				TimeRange: TimeRange{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				User: User{
					AccountID:   "user123",
					DisplayName: "Test User",
					Email:       "test@example.com",
				},
				Issues: []Issue{
					{
						Key:     "JIRA-123",
						Summary: "Test Issue",
						Status:  "In Progress",
						Changes: []Change{
							{
								Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Field:     "status",
								FromValue: "Open",
								ToValue:   "In Progress",
							},
						},
						Comments: []Comment{
							{
								Timestamp: time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
								Author:    "Test User",
								Content:   "This is a test comment",
							},
						},
					},
				},
			},
			expectedStr: "JIRA-123",
		},
	}

	// Run tests
	formatter := NewHTMLFormatter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formatter.Format(tc.report)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			if result.ContentType != "text/html" {
				t.Errorf("Expected content type 'text/html', got '%s'", result.ContentType)
			}
			
			if !strings.Contains(result.Content, tc.expectedStr) {
				t.Errorf("Expected content to contain '%s', got '%s'", tc.expectedStr, result.Content)
			}
		})
	}
} 
