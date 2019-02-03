package jira

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

// Provider returns a terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"jira_comment":          resourceComment(),
			"jira_group":            resourceGroup(),
			"jira_group_membership": resourceGroupMembership(),
			"jira_issue":            resourceIssue(),
			"jira_issue_link":       resourceIssueLink(),
			"jira_issue_type":       resourceIssueType(),
			"jira_issue_link_type":  resourceIssueLinkType(),
			"jira_project":          resourceProject(),
			"jira_webhook":          resourceWebhook(),
			"jira_role":             resourceRole(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jira_jql": resourceJQL(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure configures the provider by creating and authenticating JIRA client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var c Config
	if err := c.createAndAuthenticateClient(); err != nil {
		return nil, errors.Wrap(err, "creating config failed")
	}
	return &c, nil
}
