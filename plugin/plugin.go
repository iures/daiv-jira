package plugin

import (
	// Import contexts package
	"daiv-jira/plugin/jira"
	"fmt"

	extJira "github.com/andygrunwald/go-jira"
	plug "github.com/iures/daivplug"
)

type JiraPlugin struct {
	client    *jira.JiraClient
	config    *jira.JiraConfig
	user      *extJira.User
	formatter jira.ReportFormatter
}

// New creates a new instance of the plugin
func New() *JiraPlugin {
	return &JiraPlugin{}
}

// Name returns the unique identifier for this plugin
func (p *JiraPlugin) Name() string {
	return "daiv-jira"
}

// Manifest returns the plugin manifest
func (p *JiraPlugin) Manifest() *plug.PluginManifest {
	return &plug.PluginManifest{
		ConfigKeys: []plug.ConfigKey{
			{
				Type:        plug.ConfigTypeString,
				Key:         "jira.username",
				Name:        "Jira Username",
				Description: "The username for the Jira user",
				Required:    true,
				Secret:      false,
			},
			{
				Type:        plug.ConfigTypeString,
				Key:         "jira.token",
				Name:        "Jira API Token",
				Description: "The API token for the Jira user",
				Required:    true,
				EnvVar:      "JIRA_API_TOKEN",
			},
			{
				Type:        plug.ConfigTypeString,
				Key:         "jira.url",
				Name:        "Jira URL",
				Description: "The URL for the Jira instance",
				Required:    true,
				Secret:      false,
			},
			{
				Type:        plug.ConfigTypeString,
				Key:         "jira.project",
				Name:        "Jira Project",
				Description: "The project to generate the report for",
				Required:    true,
				Secret:      false,
			},
			{
				Type:        plug.ConfigTypeString,
				Key:         "jira.format",
				Name:        "Report Format",
				Description: "The format for the activity report (xml, json, or markdown)",
				Required:    false,
				Secret:      false,
			},
		},
	}
}

// Initialize sets up the plugin with its configuration
func (p *JiraPlugin) Initialize(settings map[string]interface{}) error {
	config := &jira.JiraConfig{
		Username: settings["jira.username"].(string),
		Token:    settings["jira.token"].(string),
		URL:      settings["jira.url"].(string),
		Project:  settings["jira.project"].(string),
	}

	client, err := jira.NewJiraClient(config)
	if err != nil {
		return fmt.Errorf("failed to create Jira client: %w", err)
	}

	p.client = client
	p.config = config
	p.user, err = p.client.GetSelf()

	if err != nil {
		return fmt.Errorf("failed to get Jira user: %w", err)
	}

	// Set the formatter based on configuration
	format, ok := settings["jira.format"].(string)
	if !ok || format == "" {
		format = "JSON" // Default to JSON if not specified
	}

	switch format {
	case "json":
		p.formatter = jira.NewJSONFormatter()
	case "markdown":
		p.formatter = jira.NewMarkdownFormatter()
	case "xml":
		p.formatter = jira.NewXMLFormatter()
	default:
		p.formatter = jira.NewJSONFormatter()
	}

	return nil
}

// Shutdown performs cleanup when the plugin is being disabled/removed
func (p *JiraPlugin) Shutdown() error {
	// No resources to clean up
	return nil
}

// GetStandupContext implements the StandupPlugin interface
func (p *JiraPlugin) GetStandupContext(timeRange plug.TimeRange) (plug.StandupContext, error) {
	// Create service
	service := jira.NewActivityService(p.client)
	
	// Get activity report from service
	report, err := service.GetActivityReport(timeRange, p.user)
	if err != nil {
		return plug.StandupContext{}, fmt.Errorf("failed to get activity report: %w", err)
	}
	
	// Format the report using the configured formatter
	content, err := p.formatter.Format(report)
	if err != nil {
		return plug.StandupContext{}, fmt.Errorf("failed to format activity report: %w", err)
	}

	return plug.StandupContext{
		PluginName: p.Name(),
		Content:    content,
	}, nil
}
