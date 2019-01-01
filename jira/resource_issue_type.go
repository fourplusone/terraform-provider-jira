package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

const issueTypeAPIEndpoint = "/rest/api/2/issuetype"

// IssueTypeRequest The struct sent to the JIRA instance to create a new Issue Type
// JIRA API docs: https://docs.atlassian.com/software/jira/docs/api/REST/7.6.1/#api/2/issuetype-createIssueType
type IssueTypeRequest struct {
	Description string `json:"description,omitempty" structs:"description,omitempty"`
	Name        string `json:"name,omitempty" structs:"name,omitempty"`
	Type        string `json:"type,omitempty" structs:"type,omitempty"`
	AvatarID    int    `json:"avatarId,omitempty" structs:"avatarId,omitempty"`
}

func resourceIssueType() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueTypeCreate,
		Read:   resourceIssueTypeRead,
		Update: resourceIssueTypeUpdate,
		Delete: resourceIssueTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"is_subtask": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"avatar_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// resourceIssueTypeCreate creates a new jira issue using the jira api
func resourceIssueTypeCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issueType := new(IssueTypeRequest)

	issueType.Name = d.Get("name").(string)
	issueType.Description = d.Get("description").(string)
	issueType.Type = "standard"

	if d.Get("is_subtask").(bool) {
		issueType.Type = "subtask"
	}

	req, err := config.jiraClient.NewRequest("POST", issueTypeAPIEndpoint, issueType)

	if err != nil {
		return errors.Wrap(err, "Creating POST Request failed")
	}

	returnedIssueType := new(jira.IssueType)

	_, err = config.jiraClient.Do(req, returnedIssueType)
	if err != nil {
		return errors.Wrap(err, "Creating IssueType Request failed")
	}

	d.SetId(returnedIssueType.ID)

	resourceIssueTypeUpdate(d, m)

	return resourceIssueTypeRead(d, m)
}

// resourceIssueTypeRead reads issue details using jira api
func resourceIssueTypeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueTypeAPIEndpoint, d.Id())
	req, err := config.jiraClient.NewRequest("GET", urlStr, nil)

	if err != nil {
		return errors.Wrap(err, "Creating IssueType Request failed")
	}

	issueType := new(jira.IssueType)

	_, err = config.jiraClient.Do(req, issueType)
	if err != nil {
		return errors.Wrap(err, "Creating IssueType Request failed")
	}

	d.Set("name", issueType.Name)
	d.Set("description", issueType.Description)
	d.Set("is_subtask", issueType.Subtask)
	d.Set("avatar_id", issueType.AvatarID)

	return nil
}

// resourceIssueTypeUpdate updates jira issue using jira api
func resourceIssueTypeUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issueType := new(IssueTypeRequest)

	issueType.Name = d.Get("name").(string)
	issueType.Description = d.Get("description").(string)
	issueType.AvatarID = d.Get("avatar_id").(int)

	urlStr := fmt.Sprintf("%s/%s", issueTypeAPIEndpoint, d.Id())
	req, err := config.jiraClient.NewRequest("PUT", urlStr, issueType)

	if err != nil {
		return errors.Wrap(err, "Creating PUT Request failed")
	}

	returnedIssueType := new(jira.IssueType)

	_, err = config.jiraClient.Do(req, returnedIssueType)
	if err != nil {
		return errors.Wrap(err, "Creating IssueType Request failed")
	}

	return resourceIssueTypeRead(d, m)
}

// resourceIssueTypeDelete deletes jira issue using the jira api
func resourceIssueTypeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueTypeAPIEndpoint, d.Id())
	req, err := config.jiraClient.NewRequest("DELETE", urlStr, nil)

	if err != nil {
		return errors.Wrap(err, "Creating DELETE Request failed")
	}

	_, err = config.jiraClient.Do(req, nil)
	if err != nil {
		return errors.Wrap(err, "Creating IssueType Request failed")
	}

	return nil
}
