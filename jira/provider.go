package jira

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

// Provider returns a terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"jira_issue": resourceIssue(),
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
