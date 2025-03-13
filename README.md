# Daiv Jira

A plugin for the daiv CLI tool.

## Project Structure

- **main.go**: Plugin entry point that exports the Plugin interface
- **plugin/plugin.go**: Core plugin implementation (configuration, lifecycle, etc.)
- **plugin/contexts/**: Directory containing context providers for LLM integration
  - **plugin/contexts/standup.go**: Implementation of the standup context provider
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

- daiv-jira.apikey: API key for the service

You can configure these settings when you first run daiv after installing the plugin.

## Usage

After installation, the plugin will be automatically loaded when you start daiv.

## Development

This plugin includes a Makefile with the following commands:

- `make build`: Build the plugin
- `make install`: Build and install the plugin
- `make clean`: Clean build artifacts
- `make tidy`: Run go mod tidy

### Adding New Context Providers

As the daiv Plugin interface grows, you can add new context providers in the
`plugin/contexts/` directory, following the pattern established by `standup.go`.
These contexts provide information to the LLM when running daiv commands.

