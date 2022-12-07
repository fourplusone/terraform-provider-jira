package jira

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// Provider returns a terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_URL", nil),
				Description: "Base url of the JIRA instance.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_USER", nil),
				Description: "User to be used",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_PASSWORD", nil),
				Description: "Password/API Key of the user",
			},
			"pat_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_PAT_TOKEN", nil),
				Description: "PAT token of the user",
			},
			"custom_auth_header_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_CUSTOM_AUTH_HEADER_KEY", nil),
				Description: "Set a custom header key",
			},
			"custom_auth_header_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_CUSTOM_AUTH_HEADER_VALUE", nil),
				Description: "Set the value for the custom header",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"jira_comment":            resourceComment(),
			"jira_filter":             resourceFilter(),
			"jira_group":              resourceGroup(),
			"jira_group_membership":   resourceGroupMembership(),
			"jira_issue":              resourceIssue(),
			"jira_issue_link":         resourceIssueLink(),
			"jira_issue_type":         resourceIssueType(),
			"jira_issue_link_type":    resourceIssueLinkType(),
			"jira_project":            resourceProject(),
			"jira_project_category":   resourceProjectCategory(),
			"jira_project_membership": resourceProjectMembership(),
			"jira_webhook":            resourceWebhook(),
			"jira_role":               resourceRole(),
			"jira_user":               resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jira_field": resourceField(),
			"jira_jql":   resourceJQL(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure configures the provider by creating and authenticating JIRA client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var c Config
	if err := c.createAndAuthenticateClient(d); err != nil {
		return nil, errors.Wrap(err, "creating config failed")
	}
	return &c, nil
}
