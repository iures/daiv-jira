# Jira Query Options

This document provides detailed information about the query options available in the Daiv Jira plugin.

## Overview

The Daiv Jira plugin allows you to customize how it queries Jira for issues. These options control what issues are included in your reports, how many results are returned, and what fields are included.

## Output Formats

The plugin supports multiple output formats for the activity report:

- **JSON**: A structured JSON format suitable for programmatic processing
- **XML**: An XML format for integration with XML-based systems
- **Markdown**: A human-readable format suitable for display in text editors and chat systems
- **HTML**: A rich HTML format with styling for viewing in web browsers

You can set the output format using the `jira.format` configuration option.

## Configuration Options

### JQL Template (`jira.query.jql_template`)

The JQL template is the base query used to find issues in Jira. It should include placeholders (`%s`) for:
1. Project key
2. Start date
3. End date

**Default**: `"project = %s AND updatedDate >= %s AND updatedDate < %s"`

**Example Custom Value**: `"project = %s AND created >= %s AND created < %s"`

This would change the query to look for issues created in the date range rather than updated.

### Assignee Filter (`jira.query.assignee_current_user`)

Controls whether to include only issues assigned to the current user.

**Default**: `true`

**Possible Values**:
- `true`: Only include issues assigned to the current user
- `false`: Include issues regardless of assignee

### Status Filter (`jira.query.status_filter`)

Filters issues by their status using JQL syntax.

**Default**: `"!= Closed"`

**Example Values**:
- `"= In Progress"`: Only include issues with status "In Progress"
- `"!= Done"`: Exclude issues with status "Done"
- `"IN (Open, 'In Progress')"`: Include issues with status "Open" or "In Progress"

**Note**: Make sure to use valid JQL syntax. The operator must be one of: `=`, `!=`, `<`, `>`, `<=`, `>=`, `~`, `!~`, `IN`, `NOT IN`, `IS`, or `IS NOT`.

### Sprint Filter (`jira.query.in_open_sprints`)

Controls whether to include only issues in open sprints.

**Default**: `true`

**Possible Values**:
- `true`: Only include issues in open sprints
- `false`: Include issues regardless of sprint status

### Maximum Results (`jira.query.max_results`)

The maximum number of issues to return from Jira.

**Default**: `100`

**Example Values**: Any positive integer (e.g., `50`, `200`)

### Fields (`jira.query.fields`)

A comma-separated list of fields to include in the Jira API response.

**Default**: `"summary,description,status,changelog,comment"`

**Example Values**:
- `"summary,status"`: Only include summary and status fields
- `"summary,description,status,priority,assignee"`: Include additional fields

## Advanced JQL Examples

Here are some examples of advanced JQL templates you might want to use:

### Issues Created or Updated in Date Range

```
(created >= %s AND created < %s) OR (updatedDate >= %s AND updatedDate < %s)
```

Note: This would require modifying the code to support additional placeholders.

### Issues with Specific Priority

```
project = %s AND updatedDate >= %s AND updatedDate < %s AND priority IN (High, Highest)
```

### Issues with Specific Labels

```
project = %s AND updatedDate >= %s AND updatedDate < %s AND labels IN (important, urgent)
```

## Troubleshooting

If you encounter errors related to JQL syntax, check that:

1. Your status filter uses valid JQL operators (`=`, `!=`, etc.)
2. Field names are correct for your Jira instance
3. Values with spaces are properly quoted in JQL (e.g., `"In Progress"`)
4. Your JQL template has the correct number of `%s` placeholders

For more information on JQL syntax, refer to the [Atlassian JQL documentation](https://support.atlassian.com/jira-software-cloud/docs/advanced-search-reference-jql-fields/). 
