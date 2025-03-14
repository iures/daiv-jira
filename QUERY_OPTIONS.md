# Jira Query Options

This document provides detailed information about the query options available in the Daiv Jira plugin.

## Overview

The Daiv Jira plugin allows you to customize how it queries Jira for issues. These options control what issues are included in your reports, how many results are returned, and what fields are included.

The plugin also intelligently filters out issues that don't have any relevant activity (comments or changes) within the specified time range, ensuring that your reports only include issues with meaningful updates.

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

## JQL Date Format

When constructing JQL queries, the plugin uses the date format `YYYY-MM-DD` (e.g., `2023-01-15`) without the time component. This is the format expected by Jira's JQL parser.

For example, a JQL query might look like:
```
project = PROJECT AND updatedDate >= 2023-01-01 AND updatedDate < 2023-01-02
```

If you're customizing the JQL template, make sure to use this date format to avoid parsing errors.

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

## Advanced Features

### Smart Filtering

In addition to the JQL-based filtering, the plugin applies a smart filtering mechanism that examines each issue's comments and changelog entries. Issues that don't have any comments or changes within the specified time range are automatically filtered out, even if they match the JQL query.

This ensures that your activity reports only include issues with meaningful activity during the time period you're interested in, reducing noise and making your reports more relevant.

For example, if an issue was updated during your specified time range but the update was just an automated field change or a comment outside your time range, it won't be included in your report.

## Troubleshooting

If you encounter errors related to JQL syntax, check that:

1. Your status filter uses valid JQL operators (`=`, `!=`, etc.)
2. Field names are correct for your Jira instance
3. Values with spaces are properly quoted in JQL (e.g., `"In Progress"`)
4. Your JQL template has the correct number of `%s` placeholders

If you're not seeing issues you expect in your report, remember that the plugin filters out issues without relevant activity (comments or changes) in the specified time range. This is by design to ensure your reports only include meaningful updates.

For more information on JQL syntax, refer to the [Atlassian JQL documentation](https://support.atlassian.com/jira-software-cloud/docs/advanced-search-reference-jql-fields/). 
