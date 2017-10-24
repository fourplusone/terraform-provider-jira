package jira

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"jira_issue": resourceIssue(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var c Config
	if err := c.createAndAuthenticateClient(); err != nil {
		return nil, errors.Wrap(err, "creating config failed")
	}
	return &c, nil
}
