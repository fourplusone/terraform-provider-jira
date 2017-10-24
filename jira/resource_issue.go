package jira

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceIssue() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueCreate,
		Read:   resourceIssueRead,
		Update: resourceIssueUpdate,
		Delete: resourceIssueDelete,

		Schema: map[string]*schema.Schema{
			"assignee": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"reporter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"summary": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed values
			"issue_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIssueCreate(d *schema.ResourceData, m interface{}) error {
	assignee := d.Get("assignee").(string)
	reporter := d.Get("reporter").(string)
	issueType := d.Get("issue_type").(string)
	description := d.Get("description").(string)
	summary := d.Get("summary").(string)
	projectKey := d.Get("project_key").(string)

	jiraClient, err := jira.NewClient(nil, "http://localhost:8080")
	if err != nil {
		return errors.Wrap(err, "creating jira client failed")
	}
	jiraClient.Authentication.SetBasicAuth("username", "password")

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: assignee,
			},
			Reporter: &jira.User{
				Name: reporter,
			},
			Description: description,
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: projectKey,
			},
			Summary: summary,
		},
	}
	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		return errors.Wrap(err, "creating jira issue failed")
	}

	d.SetId(issue.ID)

	return resourceIssueRead(d, m)
}

func resourceIssueRead(d *schema.ResourceData, m interface{}) error {
	jiraClient, err := jira.NewClient(nil, "http://localhost:8080")
	if err != nil {
		return errors.Wrap(err, "creating jira client failed")
	}
	jiraClient.Authentication.SetBasicAuth("username", "password")

	issue, _, err := jiraClient.Issue.Get(d.Id(), nil)
	if err != nil {
		return errors.Wrap(err, "getting jira issue failed")
	}

	d.Set("assignee", issue.Fields.Assignee.Name)
	d.Set("reporter", issue.Fields.Reporter.Name)
	d.Set("issue_type", issue.Fields.Type.Name)
	if issue.Fields.Description != "" {
		d.Set("description", issue.Fields.Description)
	}
	d.Set("summary", issue.Fields.Summary)
	d.Set("project_key", issue.Fields.Project.Key)

	return nil
}

func resourceIssueUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIssueDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
