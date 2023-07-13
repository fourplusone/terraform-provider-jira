package jira

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "URL for your Jira instance. Can be specified with the JIRA_URL environment variable.",
			},
			"user": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"token"},
				RequiredWith: []string{"password"},
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("JIRA_USER", nil),
				Description:  "Username for your user. Can be specified with the JIRA_USER environment variable.",
			},
			"password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_PASSWORD", nil),
				Description: "Password for the user, can also be an API Token. Can be specified with the JIRA_PASSWORD environment variable.",
			},
			"token": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JIRA_TOKEN", nil),
				Description: "Personal access token of a user. Can be specified with the JIRA_TOKEN environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"jira_comment":            resourceComment(),
			"jira_component":          resourceComponent(),
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
			"jira_custom_field":       resourceCustomField(),
			"jira_issue_type_scheme":  resourceIssueTypeScheme(),
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
