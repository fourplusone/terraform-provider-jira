package jira

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// PermissionScheme represents a JIRA Permission Scheme
type PermissionScheme struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

const permissionSchemeAPIEndpoint = "/rest/api/2/permissionscheme"

func resourcePermissionScheme() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionSchemeCreate,
		Read:   resourcePermissionSchemeRead,
		Update: resourcePermissionSchemeUpdate,
		Delete: resourcePermissionSchemeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the permission scheme",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the permission scheme",
			},
		},
	}
}

func setPermissionScheme(w *PermissionScheme, d *schema.ResourceData) {
	w.Name = d.Get("name").(string)
	w.Description = d.Get("description").(string)
}

func setPermissionSchemeResource(w *PermissionScheme, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(w.ID))
	d.Set("name", w.Name)
	d.Set("description", w.Description)
}

func resourcePermissionSchemeCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	scheme := new(PermissionScheme)
	returnedScheme := new(PermissionScheme)

	setPermissionScheme(scheme, d)

	err := request(config.jiraClient, "POST", permissionSchemeAPIEndpoint, scheme, returnedScheme)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setPermissionSchemeResource(returnedScheme, d)

	return resourcePermissionSchemeRead(d, m)
}

func resourcePermissionSchemeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", permissionSchemeAPIEndpoint, d.Id())

	scheme := new(PermissionScheme)
	err := request(config.jiraClient, "GET", urlStr, nil, scheme)

	if err != nil {
		if errors.Is(err, ResourceNotFoundError) {
			d.SetId("")
			return nil
		}
		return errors.Wrap(err, "Request failed")
	}

	setPermissionSchemeResource(scheme, d)

	return nil
}

func resourcePermissionSchemeUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	scheme := new(PermissionScheme)

	setPermissionScheme(scheme, d)

	urlStr := fmt.Sprintf("%s/%s", permissionSchemeAPIEndpoint, d.Id())
	returnedScheme := new(PermissionScheme)

	err := request(config.jiraClient, "PUT", urlStr, scheme, returnedScheme)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourcePermissionSchemeRead(d, m)
}

func resourcePermissionSchemeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", permissionSchemeAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
