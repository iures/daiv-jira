package main

import (
	"daiv-jira/plugin" // Import the plugin package

	plug "github.com/iures/daivplug"
)

// Plugin is exported as a symbol for the daiv plugin system to find.
// This Jira plugin provides integration with Jira for activity reporting.
// It allows users to query Jira issues based on configurable parameters
// and generate reports in various formats (XML, JSON, Markdown).
var Plugin plug.Plugin = plugin.New()
