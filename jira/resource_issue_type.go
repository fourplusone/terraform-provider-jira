package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

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

	returnedIssueType := new(jira.IssueType)
	err := request(config.jiraClient, "POST", issueTypeAPIEndpoint, issueType, returnedIssueType)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(returnedIssueType.ID)

	resourceIssueTypeUpdate(d, m)

	return resourceIssueTypeRead(d, m)
}

// resourceIssueTypeRead reads issue details using jira api
func resourceIssueTypeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueTypeAPIEndpoint, d.Id())

	issueType := new(jira.IssueType)
	err := request(config.jiraClient, "GET", urlStr, nil, issueType)

	if err != nil {
		return errors.Wrap(err, "Request failed")
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
	returnedIssueType := new(jira.IssueType)

	err := request(config.jiraClient, "PUT", urlStr, issueType, returnedIssueType)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceIssueTypeRead(d, m)
}

// resourceIssueTypeDelete deletes jira issue using the jira api
func resourceIssueTypeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueTypeAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
