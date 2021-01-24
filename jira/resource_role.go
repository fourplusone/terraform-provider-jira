package jira

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// Role represents a JIRA Role
type Role struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
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
			},
		},
	}
}

func setRole(w *Role, d *schema.ResourceData) {
	w.Name = d.Get("name").(string)
	w.Description = d.Get("description").(string)
}

func setRoleResource(w *Role, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(w.ID))
	d.Set("name", w.Name)
	d.Set("description", w.Description)
}

// resourceRoleCreate creates a new jira issue using the jira api
func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	role := new(Role)
	returnedRole := new(Role)

	setRole(role, d)

	err := request(config.jiraClient, "POST", roleAPIEndpoint, role, returnedRole)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setRoleResource(returnedRole, d)

	return resourceRoleRead(d, m)
}

// resourceRoleRead reads issue details using jira api
func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", roleAPIEndpoint, d.Id())

	role := new(Role)
	err := request(config.jiraClient, "GET", urlStr, nil, role)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setRoleResource(role, d)

	return nil
}

// resourceRoleUpdate updates jira issue using jira api
func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	role := new(Role)

	setRole(role, d)

	urlStr := fmt.Sprintf("%s/%s", roleAPIEndpoint, d.Id())
	returnedRole := new(Role)

	err := request(config.jiraClient, "PUT", urlStr, role, returnedRole)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceRoleRead(d, m)
}

// resourceRoleDelete deletes jira issue using the jira api
func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", roleAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
