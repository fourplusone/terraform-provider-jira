package jira

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// resourceProjectCreate creates a new jira issue using the jira api
func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	// config := m.(*Config)
	// name := d.Get("name").(string)
	// key := d.Get("key").(string)

	// config.jiraClient.Project.Get()
	// d.SetId(issue.ID)

	return resourceProjectRead(d, m)
}

// resourceProjectRead reads issue details using jira api
func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	project, _, err := config.jiraClient.Project.Get(d.Id())
	if err != nil {
		return errors.Wrap(err, "getting jira project failed")
	}

	d.Set("name", project.Name)
	d.Set("key", project.Key)

	// issue, _, err := config.jiraClient.Group.
	// if err != nil {
	// 	return errors.Wrap(err, "getting jira issue failed")
	// }

	// d.Set("assignee", issue.Fields.Assignee.Name)
	// d.Set("reporter", issue.Fields.Reporter.Name)
	// d.Set("issue_type", issue.Fields.Type.Name)
	// if issue.Fields.Description != "" {
	// 	d.Set("description", issue.Fields.Description)
	// }
	// d.Set("summary", issue.Fields.Summary)
	// d.Set("project_key", issue.Fields.Project.Key)
	// d.Set("issue_key", issue.Key)

	return nil
}

// resourceProjectUpdate updates jira issue using jira api
func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectRead(d, m)
}

// resourceProjectDelete deletes jira issue using the jira api
func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
