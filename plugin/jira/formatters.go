package jira

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// FormattedContent represents formatted content with its content type
type FormattedContent struct {
	ContentType string // MIME type of the content
	Content     string // The formatted content
}

// ReportFormatter is an interface for formatting activity reports
type ReportFormatter interface {
	Format(report *ActivityReport) (*FormattedContent, error)
	Name() string // Returns the name of the formatter
}

// XMLFormatter formats activity reports as XML
type XMLFormatter struct{}

// NewXMLFormatter creates a new XML formatter
func NewXMLFormatter() *XMLFormatter {
	return &XMLFormatter{}
}

// Name returns the name of the formatter
func (f *XMLFormatter) Name() string {
	return "xml"
}

// Format formats an activity report as XML
func (f *XMLFormatter) Format(report *ActivityReport) (*FormattedContent, error) {
	if len(report.Issues) == 0 {
		return &FormattedContent{
			ContentType: "application/xml",
			Content:     "<jira_report></jira_report>",
		}, nil
	}

	xmlReport := jiraXMLReport{
		Issues: make([]xmlIssue, 0, len(report.Issues)),
	}

	for _, issue := range report.Issues {
		xmlIssue := xmlIssue{
			Key:     issue.Key,
			Status:  issue.Status,
			Summary: issue.Summary,
		}

		// Process comments
		comments := make([]xmlComment, 0, len(issue.Comments))
		for _, comment := range issue.Comments {
			comments = append(comments, xmlComment{
				Timestamp: comment.Timestamp.Format("2006-01-02 15:04:05"),
				Author:    comment.Author,
				Content:   comment.Content,
			})
		}
		xmlIssue.Comments = xmlComments{Comments: comments}

		// Process changes
		changes := make([]xmlChange, 0, len(issue.Changes))
		for _, change := range issue.Changes {
			changes = append(changes, xmlChange{
				Timestamp: change.Timestamp.Format("2006-01-02 15:04:05"),
				Author:    change.Author,
				Field:     change.Field,
				From:      change.FromValue,
				To:        change.ToValue,
			})
		}
		xmlIssue.Changelog = xmlChangelog{Changes: changes}

		xmlReport.Issues = append(xmlReport.Issues, xmlIssue)
	}

	// Marshal to XML with proper indentation
	output, err := xml.MarshalIndent(xmlReport, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header and return
	return &FormattedContent{
		ContentType: "application/xml",
		Content:     xml.Header + string(output),
	}, nil
}

// JSONFormatter formats activity reports as JSON
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Name returns the name of the formatter
func (f *JSONFormatter) Name() string {
	return "json"
}

// Format formats an activity report as JSON
func (f *JSONFormatter) Format(report *ActivityReport) (*FormattedContent, error) {
	if len(report.Issues) == 0 {
		return &FormattedContent{
			ContentType: "application/json",
			Content:     "{}",
		}, nil
	}

	// Create a JSON-friendly structure
	type jsonComment struct {
		Timestamp string `json:"timestamp"`
		Author    string `json:"author"`
		Content   string `json:"content"`
	}

	type jsonChange struct {
		Timestamp string `json:"timestamp"`
		Author    string `json:"author"`
		Field     string `json:"field"`
		From      string `json:"from"`
		To        string `json:"to"`
	}

	type jsonIssue struct {
		Key      string        `json:"key"`
		Status   string        `json:"status"`
		Summary  string        `json:"summary"`
		Comments []jsonComment `json:"comments"`
		Changes  []jsonChange  `json:"changes"`
	}

	type jsonReport struct {
		TimeRange struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"timeRange"`
		User struct {
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
		} `json:"user"`
		Issues []jsonIssue `json:"issues"`
	}

	// Convert domain model to JSON structure
	jReport := jsonReport{}
	jReport.TimeRange.Start = report.TimeRange.Start.Format(time.RFC3339)
	jReport.TimeRange.End = report.TimeRange.End.Format(time.RFC3339)
	jReport.User.DisplayName = report.User.DisplayName
	jReport.User.Email = report.User.Email
	
	for _, issue := range report.Issues {
		jIssue := jsonIssue{
			Key:      issue.Key,
			Status:   issue.Status,
			Summary:  issue.Summary,
			Comments: make([]jsonComment, 0, len(issue.Comments)),
			Changes:  make([]jsonChange, 0, len(issue.Changes)),
		}

		for _, comment := range issue.Comments {
			jIssue.Comments = append(jIssue.Comments, jsonComment{
				Timestamp: comment.Timestamp.Format(time.RFC3339),
				Author:    comment.Author,
				Content:   comment.Content,
			})
		}

		for _, change := range issue.Changes {
			jIssue.Changes = append(jIssue.Changes, jsonChange{
				Timestamp: change.Timestamp.Format(time.RFC3339),
				Author:    change.Author,
				Field:     change.Field,
				From:      change.FromValue,
				To:        change.ToValue,
			})
		}

		jReport.Issues = append(jReport.Issues, jIssue)
	}

	// Marshal to JSON with proper indentation
	output, err := json.MarshalIndent(jReport, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return &FormattedContent{
		ContentType: "application/json",
		Content:     string(output),
	}, nil
}

// MarkdownFormatter formats activity reports as Markdown
type MarkdownFormatter struct{}

// NewMarkdownFormatter creates a new Markdown formatter
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// Name returns the name of the formatter
func (f *MarkdownFormatter) Name() string {
	return "markdown"
}

// Format formats an activity report as Markdown
func (f *MarkdownFormatter) Format(report *ActivityReport) (*FormattedContent, error) {
	if len(report.Issues) == 0 {
		return &FormattedContent{
			ContentType: "text/markdown",
			Content:     "No activity found for the specified time range.",
		}, nil
	}

	var sb strings.Builder

	// Add report header
	sb.WriteString(fmt.Sprintf("# Jira Activity Report\n\n"))
	sb.WriteString(fmt.Sprintf("**Time Range:** %s to %s\n\n", 
		report.TimeRange.Start.Format("2006-01-02"),
		report.TimeRange.End.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**User:** %s (%s)\n\n", 
		report.User.DisplayName, 
		report.User.Email))
	
	// Group issues by status
	statusGroups := make(map[string][]Issue)
	for _, issue := range report.Issues {
		statusGroups[issue.Status] = append(statusGroups[issue.Status], issue)
	}

	// Add issues by status
	for status, issues := range statusGroups {
		sb.WriteString(fmt.Sprintf("## %s Issues\n\n", status))
		
		for _, issue := range issues {
			sb.WriteString(fmt.Sprintf("### [%s] %s\n\n", issue.Key, issue.Summary))
			
			// Add changes section if there are any
			if len(issue.Changes) > 0 {
				sb.WriteString("#### Changes\n\n")
				sb.WriteString("| Time | Field | From | To |\n")
				sb.WriteString("|------|-------|------|----|\n")
				
				for _, change := range issue.Changes {
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
						change.Timestamp.Format("2006-01-02 15:04"),
						change.Field,
						change.FromValue,
						change.ToValue))
				}
				sb.WriteString("\n")
			}
			
			// Add comments section if there are any
			if len(issue.Comments) > 0 {
				sb.WriteString("#### Comments\n\n")
				
				for _, comment := range issue.Comments {
					sb.WriteString(fmt.Sprintf("**%s** - %s\n\n", 
						comment.Author,
						comment.Timestamp.Format("2006-01-02 15:04")))
					sb.WriteString(fmt.Sprintf("%s\n\n", comment.Content))
				}
			}
			
			sb.WriteString("---\n\n")
		}
	}

	return &FormattedContent{
		ContentType: "text/markdown",
		Content:     sb.String(),
	}, nil
}

// HTMLFormatter formats activity reports as HTML
type HTMLFormatter struct{}

// NewHTMLFormatter creates a new HTML formatter
func NewHTMLFormatter() *HTMLFormatter {
	return &HTMLFormatter{}
}

// Name returns the name of the formatter
func (f *HTMLFormatter) Name() string {
	return "html"
}

// Format formats an activity report as HTML
func (f *HTMLFormatter) Format(report *ActivityReport) (*FormattedContent, error) {
	if len(report.Issues) == 0 {
		return &FormattedContent{
			ContentType: "text/html",
			Content:     "<html><body><h1>Jira Activity Report</h1><p>No activity found for the specified time range.</p></body></html>",
		}, nil
	}

	var sb strings.Builder

	// Start HTML document
	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	sb.WriteString("<title>Jira Activity Report</title>\n")
	sb.WriteString("<style>\n")
	sb.WriteString("body { font-family: Arial, sans-serif; margin: 20px; }\n")
	sb.WriteString("h1 { color: #0052CC; }\n") // Jira blue
	sb.WriteString("h2 { color: #172B4D; border-bottom: 1px solid #DFE1E6; padding-bottom: 8px; }\n")
	sb.WriteString("h3 { margin-top: 20px; }\n")
	sb.WriteString(".issue { background-color: #F4F5F7; border-radius: 3px; padding: 15px; margin-bottom: 15px; }\n")
	sb.WriteString(".issue-key { color: #0052CC; font-weight: bold; }\n")
	sb.WriteString(".issue-summary { font-size: 16px; margin-bottom: 10px; }\n")
	sb.WriteString(".metadata { color: #6B778C; font-size: 14px; margin-bottom: 15px; }\n")
	sb.WriteString(".changes, .comments { margin-top: 10px; }\n")
	sb.WriteString(".change, .comment { background-color: white; border: 1px solid #DFE1E6; padding: 10px; margin-bottom: 8px; }\n")
	sb.WriteString(".author { color: #0052CC; font-weight: bold; }\n")
	sb.WriteString(".timestamp { color: #6B778C; font-size: 12px; }\n")
	sb.WriteString("</style>\n")
	sb.WriteString("</head>\n<body>\n")

	// Add report header
	sb.WriteString("<h1>Jira Activity Report</h1>\n")
	sb.WriteString("<div class=\"metadata\">\n")
	sb.WriteString(fmt.Sprintf("<p><strong>Time Range:</strong> %s to %s</p>\n", 
		report.TimeRange.Start.Format("2006-01-02"),
		report.TimeRange.End.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("<p><strong>User:</strong> %s (%s)</p>\n", 
		report.User.DisplayName, 
		report.User.Email))
	sb.WriteString("</div>\n")
	
	// Group issues by status
	statusGroups := make(map[string][]Issue)
	for _, issue := range report.Issues {
		statusGroups[issue.Status] = append(statusGroups[issue.Status], issue)
	}

	// Add issues by status
	for status, issues := range statusGroups {
		sb.WriteString(fmt.Sprintf("<h2>%s Issues</h2>\n", status))
		
		for _, issue := range issues {
			sb.WriteString("<div class=\"issue\">\n")
			sb.WriteString(fmt.Sprintf("<h3><span class=\"issue-key\">[%s]</span> <span class=\"issue-summary\">%s</span></h3>\n", 
				issue.Key, issue.Summary))
			
			// Add changes section if there are any
			if len(issue.Changes) > 0 {
				sb.WriteString("<div class=\"changes\">\n")
				sb.WriteString("<h4>Changes</h4>\n")
				for _, change := range issue.Changes {
					sb.WriteString("<div class=\"change\">\n")
					sb.WriteString(fmt.Sprintf("<p><span class=\"author\">%s</span> changed <strong>%s</strong> from \"%s\" to \"%s\"</p>\n", 
						change.Author, change.Field, change.FromValue, change.ToValue))
					sb.WriteString(fmt.Sprintf("<p class=\"timestamp\">%s</p>\n", 
						change.Timestamp.Format("2006-01-02 15:04:05")))
					sb.WriteString("</div>\n")
				}
				sb.WriteString("</div>\n")
			}
			
			// Add comments section if there are any
			if len(issue.Comments) > 0 {
				sb.WriteString("<div class=\"comments\">\n")
				sb.WriteString("<h4>Comments</h4>\n")
				for _, comment := range issue.Comments {
					sb.WriteString("<div class=\"comment\">\n")
					sb.WriteString(fmt.Sprintf("<p><span class=\"author\">%s</span></p>\n", comment.Author))
					sb.WriteString(fmt.Sprintf("<p>%s</p>\n", comment.Content))
					sb.WriteString(fmt.Sprintf("<p class=\"timestamp\">%s</p>\n", 
						comment.Timestamp.Format("2006-01-02 15:04:05")))
					sb.WriteString("</div>\n")
				}
				sb.WriteString("</div>\n")
			}
			
			sb.WriteString("</div>\n")
		}
	}
	
	// Close HTML document
	sb.WriteString("</body>\n</html>")

	return &FormattedContent{
		ContentType: "text/html",
		Content:     sb.String(),
	}, nil
}

// XML structures for proper marshaling
type jiraXMLReport struct {
	XMLName xml.Name   `xml:"jira_report"`
	Issues  []xmlIssue `xml:"issue"`
}

type xmlIssue struct {
	Key      string      `xml:"key"`
	Status   string      `xml:"status"`
	Summary  string      `xml:"summary"`
	Comments xmlComments `xml:"comments"`
	Changelog xmlChangelog `xml:"changelog"`
}

type xmlComments struct {
	Comments []xmlComment `xml:"comment"`
}

type xmlComment struct {
	Timestamp string `xml:"timestamp"`
	Author    string `xml:"author"`
	Content   string `xml:"content"`
}

type xmlChangelog struct {
	Changes []xmlChange `xml:"change"`
}

type xmlChange struct {
	Timestamp string `xml:"timestamp"`
	Author    string `xml:"author"`
	Field     string `xml:"field"`
	From      string `xml:"from"`
	To        string `xml:"to"`
} 
