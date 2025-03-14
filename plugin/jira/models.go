package jira

import (
	"time"
)

// Domain models for Jira activity data
// These models are independent of both the external Jira API and presentation formats

// ActivityReport represents processed activity data for a specific time range
type ActivityReport struct {
	TimeRange TimeRange
	User      User
	Issues    []Issue
}

// TimeRange represents a time period for the report
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// IsInRange checks if a given time is within the time range
func (tr TimeRange) IsInRange(t time.Time) bool {
	return (t.Equal(tr.Start) || t.After(tr.Start)) && t.Before(tr.End)
}

// User represents a Jira user
type User struct {
	AccountID   string
	DisplayName string
	Email       string
}

// Issue represents a Jira issue with relevant activity data
type Issue struct {
	Key     string
	Summary string
	Status  string
	Comments []Comment
	Changes  []Change
}

// Comment represents a comment on a Jira issue
type Comment struct {
	Timestamp time.Time
	Author    string
	Content   string
}

// Change represents a change to a Jira issue
type Change struct {
	Timestamp time.Time
	Author    string
	Field     string
	FromValue string
	ToValue   string
}

// QueryOptions represents configurable options for Jira queries
type QueryOptions struct {
	// JQL template with placeholders for dynamic values
	JQLTemplate string
	
	// Whether to include only issues assigned to the current user
	AssigneeCurrentUser bool
	
	// Project key to filter issues by
	Project string
	
	// Status filter (e.g., "!= Closed" to exclude closed issues)
	StatusFilter string
	
	// Whether to include only issues in open sprints
	InOpenSprints bool
	
	// Maximum number of results to return
	MaxResults int
	
	// Fields to include in the response
	Fields []string
	
	// Whether to expand changelog in the response
	ExpandChangelog bool
}

// DefaultQueryOptions returns the default query options
func DefaultQueryOptions() QueryOptions {
	return QueryOptions{
		JQLTemplate:       "project = %s AND updatedDate >= %s AND updatedDate < %s",
		AssigneeCurrentUser: true,
		StatusFilter:      "!= Closed",
		InOpenSprints:     true,
		MaxResults:        100,
		Fields:            []string{"summary", "description", "status", "changelog", "comment"},
		ExpandChangelog:   true,
	}
} 
