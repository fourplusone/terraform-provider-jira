package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

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

	returnedIssueLinkType := new(jira.IssueLinkType)

	err := request(config.jiraClient, "POST", issueLinkTypeAPIEndpoint, issueLinkType, returnedIssueLinkType)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(returnedIssueLinkType.ID)

	return resourceIssueLinkTypeRead(d, m)
}

// resourceIssueLinkTypeRead reads issue details using jira api
func resourceIssueLinkTypeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	urlStr := fmt.Sprintf("%s/%s", issueLinkTypeAPIEndpoint, d.Id())
	issueLinkType := new(jira.IssueLinkType)

	err := request(config.jiraClient, "GET", urlStr, nil, issueLinkType)
	if err != nil {
		return errors.Wrap(err, "Request failed")
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
	returnedIssueLinkType := new(jira.IssueLinkType)

	err := request(config.jiraClient, "PUT", urlStr, issueLinkType, returnedIssueLinkType)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceIssueLinkTypeRead(d, m)
}

// resourceIssueLinkTypeDelete deletes jira issue using the jira api
func resourceIssueLinkTypeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueLinkTypeAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
