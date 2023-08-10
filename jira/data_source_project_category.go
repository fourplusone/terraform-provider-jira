package jira

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func datasourceProjectCategory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectCategoryRead,
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectCategoryRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", projectCategoryAPIEndpoint, d.Get("project_id"))

	projectCategory := new(ProjectCategory)
	err := request(config.jiraClient, "GET", urlStr, nil, projectCategory)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setProjectCategoryResource(projectCategory, d)

	return nil
}
