package jira

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeRange_IsInRange(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		timeRange TimeRange
		testTime time.Time
		expected bool
	}{
		{
			name: "Time is in range",
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			testTime: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "Time is equal to start (inclusive)",
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			testTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "Time is equal to end (exclusive)",
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			testTime: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "Time is before range",
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			testTime: time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "Time is after range",
			timeRange: TimeRange{
				Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			testTime: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.timeRange.IsInRange(tc.testTime)
			if result != tc.expected {
				t.Errorf("Expected IsInRange to return %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultQueryOptions(t *testing.T) {
	// Get default options
	options := DefaultQueryOptions()

	// Test default values
	if options.JQLTemplate != "project = %s AND updatedDate >= %s AND updatedDate < %s" {
		t.Errorf("Expected default JQLTemplate to be 'project = %%s AND updatedDate >= %%s AND updatedDate < %%s', got '%s'", options.JQLTemplate)
	}

	if !options.AssigneeCurrentUser {
		t.Errorf("Expected default AssigneeCurrentUser to be true, got false")
	}

	if options.StatusFilter != "!= Closed" {
		t.Errorf("Expected default StatusFilter to be '!= Closed', got '%s'", options.StatusFilter)
	}

	if !options.InOpenSprints {
		t.Errorf("Expected default InOpenSprints to be true, got false")
	}

	if options.MaxResults != 100 {
		t.Errorf("Expected default MaxResults to be 100, got %d", options.MaxResults)
	}

	expectedFields := []string{"summary", "description", "status", "changelog", "comment"}
	if !reflect.DeepEqual(options.Fields, expectedFields) {
		t.Errorf("Expected default Fields to be %v, got %v", expectedFields, options.Fields)
	}

	if !options.ExpandChangelog {
		t.Errorf("Expected default ExpandChangelog to be true, got false")
	}
} 
