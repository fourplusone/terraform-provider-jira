package jira

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// ProjectCategory represents a JIRA ProjectCategory
type ProjectCategory struct {
	Self        string `json:"self,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ID          string `json:"id,omitempty"`
}

func resourceProjectCategory() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCategoryCreate,
		Read:   resourceProjectCategoryRead,
		Update: resourceProjectCategoryUpdate,
		Delete: resourceProjectCategoryDelete,
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

func setProjectCategory(w *ProjectCategory, d *schema.ResourceData) {
	w.Name = d.Get("name").(string)
	w.Description = d.Get("description").(string)
}

func setProjectCategoryResource(w *ProjectCategory, d *schema.ResourceData) {
	d.SetId(w.ID)
	d.Set("name", w.Name)
	d.Set("description", w.Description)
}

// resourceProjectCategoryCreate creates a new jira issue using the jira api
func resourceProjectCategoryCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	projectCategory := new(ProjectCategory)
	returnedProjectCategory := new(ProjectCategory)

	setProjectCategory(projectCategory, d)

	err := request(config.jiraClient, "POST", projectCategoryAPIEndpoint, projectCategory, returnedProjectCategory)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setProjectCategoryResource(returnedProjectCategory, d)

	return resourceProjectCategoryRead(d, m)
}

// resourceProjectCategoryRead reads issue details using jira api
func resourceProjectCategoryRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", projectCategoryAPIEndpoint, d.Id())

	projectCategory := new(ProjectCategory)
	err := request(config.jiraClient, "GET", urlStr, nil, projectCategory)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setProjectCategoryResource(projectCategory, d)

	return nil
}

// resourceProjectCategoryUpdate updates jira issue using jira api
func resourceProjectCategoryUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	projectCategory := new(ProjectCategory)

	setProjectCategory(projectCategory, d)

	urlStr := fmt.Sprintf("%s/%s", projectCategoryAPIEndpoint, d.Id())
	returnedProjectCategory := new(ProjectCategory)

	err := request(config.jiraClient, "PUT", urlStr, projectCategory, returnedProjectCategory)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceProjectCategoryRead(d, m)
}

// resourceProjectCategoryDelete deletes jira issue using the jira api
func resourceProjectCategoryDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", projectCategoryAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
