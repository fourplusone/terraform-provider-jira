package jira

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// ProjectPermissionSchemeResponse represents the response from JIRA when reading
// a project's permission scheme
type ProjectPermissionSchemeResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectPermissionSchemeRequest represents the request body for assigning
// a permission scheme to a project
type ProjectPermissionSchemeRequest struct {
	ID int `json:"id"`
}

func projectPermissionSchemeAPIEndpoint(projectKey string) string {
	return fmt.Sprintf("/rest/api/2/project/%s/permissionscheme", projectKey)
}

func resourceProjectPermissionScheme() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectPermissionSchemeCreate,
		Read:   resourceProjectPermissionSchemeRead,
		Update: resourceProjectPermissionSchemeUpdate,
		Delete: resourceProjectPermissionSchemeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectPermissionSchemeImport,
		},

		Schema: map[string]*schema.Schema{
			"project_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key of the project to assign the permission scheme to",
			},
			"scheme_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the permission scheme to assign to the project",
			},
			"scheme_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the assigned permission scheme",
			},
		},
	}
}

func resourceProjectPermissionSchemeCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	projectKey := d.Get("project_key").(string)
	schemeID := d.Get("scheme_id").(int)

	req := &ProjectPermissionSchemeRequest{
		ID: schemeID,
	}

	returnedScheme := new(ProjectPermissionSchemeResponse)

	err := request(config.jiraClient, "PUT", projectPermissionSchemeAPIEndpoint(projectKey), req, returnedScheme)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(projectKey)

	return resourceProjectPermissionSchemeRead(d, m)
}

func resourceProjectPermissionSchemeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	projectKey := d.Id()

	returnedScheme := new(ProjectPermissionSchemeResponse)

	err := request(config.jiraClient, "GET", projectPermissionSchemeAPIEndpoint(projectKey), nil, returnedScheme)

	if err != nil {
		if errors.Is(err, ResourceNotFoundError) {
			d.SetId("")
			return nil
		}
		return errors.Wrap(err, "Request failed")
	}

	d.Set("project_key", projectKey)
	d.Set("scheme_id", returnedScheme.ID)
	d.Set("scheme_name", returnedScheme.Name)

	return nil
}

func resourceProjectPermissionSchemeUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	projectKey := d.Get("project_key").(string)
	schemeID := d.Get("scheme_id").(int)

	req := &ProjectPermissionSchemeRequest{
		ID: schemeID,
	}

	returnedScheme := new(ProjectPermissionSchemeResponse)

	err := request(config.jiraClient, "PUT", projectPermissionSchemeAPIEndpoint(projectKey), req, returnedScheme)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceProjectPermissionSchemeRead(d, m)
}

// resourceProjectPermissionSchemeDelete resets the project to the default
// permission scheme (ID 0) since every project must have a permission scheme.
func resourceProjectPermissionSchemeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	projectKey := d.Get("project_key").(string)

	req := &ProjectPermissionSchemeRequest{
		ID: 0,
	}

	err := request(config.jiraClient, "PUT", projectPermissionSchemeAPIEndpoint(projectKey), req, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}

func resourceProjectPermissionSchemeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("project_key", d.Id())
	return []*schema.ResourceData{d}, nil
}
