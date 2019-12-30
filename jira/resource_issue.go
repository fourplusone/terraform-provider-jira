package jira

import (
	"io/ioutil"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// resourceIssue is used to define a JIRA issue
func resourceIssue() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueCreate,
		Read:   resourceIssueRead,
		Update: resourceIssueUpdate,
		Delete: resourceIssueDelete,

		Schema: map[string]*schema.Schema{
			"assignee": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: caseInsensitiveSuppressFunc,
			},
			"reporter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return caseInsensitiveSuppressFunc(k, old, new, d)
				},
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
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return old == new
				},
			},
			"state_transition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_transition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"issue_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceIssueCreate creates a new jira issue using the jira api
func resourceIssueCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	assignee := d.Get("assignee")
	reporter := d.Get("reporter")
	issueType := d.Get("issue_type").(string)
	description := d.Get("description").(string)
	summary := d.Get("summary").(string)
	projectKey := d.Get("project_key").(string)

	i := jira.Issue{
		Fields: &jira.IssueFields{
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

	if assignee != "" {
		i.Fields.Assignee = &jira.User{
			Name: assignee.(string),
		}
	}

	if reporter != "" {
		i.Fields.Reporter = &jira.User{
			Name: reporter.(string),
		}
	}

	issue, res, err := config.jiraClient.Issue.Create(&i)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "creating jira issue failed: %s", body)
	}

	issue, res, err = config.jiraClient.Issue.Get(issue.ID, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	if state, ok := d.GetOk("state"); ok {
		if issue.Fields.Status.ID != state.(string) {
			if transition, ok := d.GetOk("state_transition"); ok {
				res, err := config.jiraClient.Issue.DoTransition(issue.ID, transition.(string))
				if err != nil {
					body, _ := ioutil.ReadAll(res.Body)
					return errors.Wrapf(err, "transitioning jira issue failed: %s", body)
				}
			}
		}
	}

	d.SetId(issue.ID)

	return resourceIssueRead(d, m)
}

// resourceIssueRead reads issue details using jira api
func resourceIssueRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issue, res, err := config.jiraClient.Issue.Get(d.Id(), nil)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	if issue.Fields.Assignee != nil {
		d.Set("assignee", issue.Fields.Assignee.Name)
	}

	if issue.Fields.Reporter != nil {
		d.Set("reporter", issue.Fields.Reporter.Name)
	}

	d.Set("issue_type", issue.Fields.Type.Name)
	if issue.Fields.Description != "" {
		d.Set("description", issue.Fields.Description)
	}
	d.Set("summary", issue.Fields.Summary)
	d.Set("project_key", issue.Fields.Project.Key)
	d.Set("issue_key", issue.Key)
	d.Set("state", issue.Fields.Status.ID)

	return nil
}

// resourceIssueUpdate updates jira issue using jira api
func resourceIssueUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	assignee := d.Get("assignee")
	reporter := d.Get("reporter")
	issueType := d.Get("issue_type").(string)
	description := d.Get("description").(string)
	summary := d.Get("summary").(string)
	projectKey := d.Get("project_key").(string)
	issueKey := d.Get("issue_key").(string)

	i := jira.Issue{
		Key: issueKey,
		ID:  d.Id(),
		Fields: &jira.IssueFields{
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

	if assignee != "" {
		i.Fields.Assignee = &jira.User{
			Name: assignee.(string),
		}
	}

	if reporter != "" {
		i.Fields.Reporter = &jira.User{
			Name: reporter.(string),
		}
	}

	issue, res, err := config.jiraClient.Issue.Update(&i)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "updating jira issue failed: %s", body)
	}

	issue, res, err = config.jiraClient.Issue.Get(issue.ID, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	if state, ok := d.GetOk("state"); ok {
		if issue.Fields.Status.ID != state.(string) {
			if transition, ok := d.GetOk("state_transition"); ok {
				res, err := config.jiraClient.Issue.DoTransition(issue.ID, transition.(string))
				if err != nil {
					body, _ := ioutil.ReadAll(res.Body)
					return errors.Wrapf(err, "transitioning jira issue failed: %s", body)
				}
			}
		}
	}

	d.SetId(issue.ID)

	return resourceIssueRead(d, m)
}

// resourceIssueDelete deletes jira issue using the jira api
func resourceIssueDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	id := d.Id()

	if transition, ok := d.GetOk("delete_transition"); ok {
		res, err := config.jiraClient.Issue.DoTransition(id, transition.(string))
		if err != nil {
			body, _ := ioutil.ReadAll(res.Body)
			return errors.Wrapf(err, "deleting jira issue failed: %s", body)
		}

	} else {
		res, err := config.jiraClient.Issue.Delete(id)

		if err != nil {
			body, _ := ioutil.ReadAll(res.Body)
			return errors.Wrapf(err, "deleting jira issue failed: %s", body)
		}
	}

	return nil
}
