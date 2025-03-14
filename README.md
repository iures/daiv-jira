# Daiv Jira

A Jira integration plugin for the daiv CLI tool. This plugin allows you to generate activity reports from Jira issues for use in standup meetings and other contexts.

## Features

- Retrieves Jira issues based on configurable query parameters
- Filters issues by time range, status, assignee, and more
- Supports multiple output formats (XML, JSON, Markdown, HTML)
- Fully configurable JQL queries
- Customizable field selection

## Project Structure

- **main.go**: Plugin entry point that exports the Plugin interface
- **plugin/plugin.go**: Core plugin implementation (configuration, lifecycle, etc.)
- **plugin/jira/**: Directory containing Jira integration components
  - **plugin/jira/client.go**: Jira API client implementation
  - **plugin/jira/models.go**: Domain models for Jira data
  - **plugin/jira/repository.go**: Data access layer for Jira
  - **plugin/jira/service.go**: Business logic for processing Jira data
  - **plugin/jira/formatters.go**: Output formatters (XML, JSON, Markdown)
- **Makefile**: Build automation for the plugin

## Installation

### From GitHub

```
daiv plugin install YOUR_GITHUB_USERNAME/daiv-jira
```

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/YOUR_GITHUB_USERNAME/daiv-jira.git
   cd daiv-jira
   ```

2. Build the plugin:
   ```
   make install
   ```
   
   Or manually:
   ```
   go build -o out/daiv-jira.so -buildmode=plugin
   daiv plugin install ./out/daiv-jira.so
   ```

## Configuration

This plugin requires the following configuration:

### Required Settings

- **jira.username**: Your Jira username
- **jira.token**: Your Jira API token
- **jira.url**: The URL of your Jira instance
- **jira.project**: The Jira project key to query

### Optional Settings

- **jira.format**: Output format (xml, json, markdown, or html)
- **jira.query.jql_template**: Custom JQL template with placeholders for project, start date, and end date
- **jira.query.assignee_current_user**: Whether to include only issues assigned to the current user (true/false)
- **jira.query.status_filter**: Filter issues by status using JQL syntax (e.g., '!= Closed' to exclude closed issues)
- **jira.query.in_open_sprints**: Whether to include only issues in open sprints (true/false)
- **jira.query.max_results**: Maximum number of results to return
- **jira.query.fields**: Comma-separated list of fields to include in the response

You can configure these settings when you first run daiv after installing the plugin, or by using the `daiv config set` command.

## Usage

After installation and configuration, the plugin will be automatically loaded when you start daiv.

### Generating a Standup Report

```
daiv standup
```

This will generate a report of your Jira activity for the default time range (usually the last 24 hours).

### Customizing the Time Range

```
daiv standup --from "2023-03-01" --to "2023-03-14"
```

### Changing the Output Format

You can change the default output format in the configuration, or specify it for a single command:

```
daiv config set jira.format markdown
```

## Development

This plugin includes a Makefile with the following commands:

- `make build`: Build the plugin
- `make install`: Build and install the plugin
- `make clean`: Clean build artifacts
- `make tidy`: Run go mod tidy

## Architecture

The plugin follows a clean architecture approach with clear separation of concerns:

1. **Domain Models**: Independent data structures representing Jira entities
2. **Repository Layer**: Handles data access to the Jira API
3. **Service Layer**: Contains business logic for processing Jira data
4. **Formatters**: Transform domain models into different output formats
5. **Plugin Layer**: Integrates with the daiv CLI tool

This architecture makes the plugin flexible, maintainable, and testable.

