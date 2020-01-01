package jira

import (
	"io/ioutil"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// resourceComment is used to define a JIRA comment
func resourceComment() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommentCreate,
		Read:   resourceCommentRead,
		Update: resourceCommentUpdate,
		Delete: resourceCommentDelete,

		Schema: map[string]*schema.Schema{
			"body": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceCommentCreate creates a new jira comment using the jira api
func resourceCommentCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	body := d.Get("body").(string)
	issueKey := d.Get("issue_key").(string)

	c := jira.Comment{Body: body}

	comment, res, err := config.jiraClient.Issue.AddComment(issueKey, &c)

	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "creating jira issue failed: %s", body)
	}

	d.SetId(comment.ID)

	return resourceCommentRead(d, m)
}

// resourceCommentRead reads comment details using jira api
func resourceCommentRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issue, res, err := config.jiraClient.Issue.Get(d.Get("issue_key").(string), nil)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	var comment *jira.Comment
	for _, c := range issue.Fields.Comments.Comments {
		if c.ID == d.Id() {
			comment = c
			break
		}
	}

	if comment == nil {
		d.SetId("")
		return nil
	}

	d.Set("body", comment.Body)

	return nil
}

// resourceCommentUpdate updates jira comment using jira api
func resourceCommentUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	body := d.Get("body").(string)
	issueKey := d.Get("issue_key").(string)

	i := jira.Comment{
		ID:   d.Id(),
		Body: body,
	}

	comment, res, err := config.jiraClient.Issue.UpdateComment(issueKey, &i)

	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "updating jira comment failed: %s", body)
	}

	d.SetId(comment.ID)

	return resourceCommentRead(d, m)
}

// resourceCommentDelete deletes jira comment using the jira api
func resourceCommentDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	issueKey := d.Get("issue_key").(string)
	id := d.Id()

	err := config.jiraClient.Issue.DeleteComment(issueKey, id)
	if err != nil {
		return errors.Wrapf(err, "deleting jira comment failed")
	}
	return nil
}
