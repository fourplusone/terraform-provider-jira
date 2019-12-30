package jira

import (
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// GroupRequest The struct sent to the JIRA instance to create a new Group
type GroupRequest struct {
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

// resourceGroup is used to define a JIRA issue
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Delete: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceGroupCreate creates a new jira issue using the jira api
func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	group := new(GroupRequest)
	group.Name = d.Get("name").(string)

	err := request(config.jiraClient, "POST", groupAPIEndpoint, group, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(group.Name)

	return resourceGroupRead(d, m)
}

// resourceGroupRead reads issue details using jira api
func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	_, _, err := config.jiraClient.Group.Get(d.Id())
	if err != nil {
		return errors.Wrap(err, "getting jira group failed")
	}

	d.Set("name", d.Id())

	return nil
}

// resourceGroupDelete deletes jira issue using the jira api
func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	relativeURL, _ := url.Parse(groupAPIEndpoint)

	query := relativeURL.Query()
	query.Set("groupname", d.Get("name").(string))

	relativeURL.RawQuery = query.Encode()

	err := request(config.jiraClient, "DELETE", relativeURL.String(), nil, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
