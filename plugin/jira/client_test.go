package jira

import (
	"testing"
)

func TestNewJiraClient(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		config      *JiraConfig
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: &JiraConfig{
				Username: "test",
				Token:    "test",
				URL:      "https://test.atlassian.net",
				Project:  "TEST",
				QueryOptions: DefaultQueryOptions(),
			},
			expectError: false,
		},
		{
			name: "Invalid URL",
			config: &JiraConfig{
				Username: "test",
				Token:    "test",
				URL:      "://invalid-url",
				Project:  "TEST",
				QueryOptions: DefaultQueryOptions(),
			},
			expectError: true,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the method being tested
			client, err := NewJiraClient(tc.config)

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If no error is expected, check the client
			if !tc.expectError && err == nil {
				if client == nil {
					t.Errorf("Expected a non-nil client but got nil")
				}
				if client.config != tc.config {
					t.Errorf("Expected client config to be %v, got %v", tc.config, client.config)
				}
				if client.client == nil {
					t.Errorf("Expected a non-nil extJira.Client but got nil")
				}
				if client.repository == nil {
					t.Errorf("Expected a non-nil JiraRepository but got nil")
				}
			}
		})
	}
}

func TestJiraClient_GetRepository(t *testing.T) {
	// Create a client
	config := &JiraConfig{
		Username: "test",
		Token:    "test",
		URL:      "https://test.atlassian.net",
		Project:  "TEST",
		QueryOptions: DefaultQueryOptions(),
	}
	client, err := NewJiraClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the method being tested
	repo := client.GetRepository()

	// Check the repository
	if repo == nil {
		t.Errorf("Expected a non-nil repository but got nil")
	}
} 
