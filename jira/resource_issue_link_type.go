package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

const issueLinkTypeAPIEndpoint = "/rest/api/2/issueLinkType"

func resourceIssueLinkType() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueLinkTypeCreate,
		Read:   resourceIssueLinkTypeRead,
		Update: resourceIssueLinkTypeUpdate,
		Delete: resourceIssueLinkTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"inward": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"outward": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// resourceIssueLinkTypeCreate creates a new jira issue using the jira api
func resourceIssueLinkTypeCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issueLinkType := new(jira.IssueLinkType)

	issueLinkType.Name = d.Get("name").(string)
	issueLinkType.Inward = d.Get("inward").(string)
	issueLinkType.Outward = d.Get("outward").(string)

	req, err := config.jiraClient.NewRequest("POST", issueLinkTypeAPIEndpoint, issueLinkType)

	if err != nil {
		return errors.Wrap(err, "Creating POST Request failed")
	}

	returnedIssueLinkType := new(jira.IssueLinkType)

	_, err = config.jiraClient.Do(req, returnedIssueLinkType)
	if err != nil {
		return errors.Wrap(err, "Creating IssueLinkType Request failed")
	}

	d.SetId(returnedIssueLinkType.ID)

	return resourceIssueLinkTypeRead(d, m)
}

// resourceIssueLinkTypeRead reads issue details using jira api
func resourceIssueLinkTypeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueLinkTypeAPIEndpoint, d.Id())
	req, err := config.jiraClient.NewRequest("GET", urlStr, nil)

	if err != nil {
		return errors.Wrap(err, "Creating IssueLinkType Request failed")
	}

	issueLinkType := new(jira.IssueLinkType)

	_, err = config.jiraClient.Do(req, issueLinkType)
	if err != nil {
		return errors.Wrap(err, "Creating IssueLinkType Request failed")
	}

	d.Set("name", issueLinkType.Name)
	d.Set("inward", issueLinkType.Inward)
	d.Set("outward", issueLinkType.Outward)

	return nil
}

// resourceIssueLinkTypeUpdate updates jira issue using jira api
func resourceIssueLinkTypeUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issueLinkType := new(jira.IssueLinkType)

	issueLinkType.Name = d.Get("name").(string)
	issueLinkType.Inward = d.Get("inward").(string)
	issueLinkType.Outward = d.Get("outward").(string)

	urlStr := fmt.Sprintf("%s/%s", issueLinkTypeAPIEndpoint, d.Id())
	req, err := config.jiraClient.NewRequest("PUT", urlStr, issueLinkType)

	if err != nil {
		return errors.Wrap(err, "Creating PUT Request failed")
	}

	returnedIssueLinkType := new(jira.IssueLinkType)

	_, err = config.jiraClient.Do(req, returnedIssueLinkType)
	if err != nil {
		return errors.Wrap(err, "Creating IssueLinkType Request failed")
	}

	return resourceIssueLinkTypeRead(d, m)
}

// resourceIssueLinkTypeDelete deletes jira issue using the jira api
func resourceIssueLinkTypeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueLinkTypeAPIEndpoint, d.Id())
	req, err := config.jiraClient.NewRequest("DELETE", urlStr, nil)

	if err != nil {
		return errors.Wrap(err, "Creating DELETE Request failed")
	}

	_, err = config.jiraClient.Do(req, nil)
	if err != nil {
		return errors.Wrap(err, "Creating IssueLinkType Request failed")
	}

	return nil
}
