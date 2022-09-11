package jira

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentCreate,
		Read:   resourceComponentRead,
		Update: resourceComponentUpdate,
		Delete: resourceComponentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"assignee_type": {
				Type: schema.TypeString,
				ValidateFunc: func(val interface{}, s string) ([]string, []error) {
					v := strings.ToLower(val.(string))
					if !(v == "project_default" ||
						v == "component_lead" ||
						v == "project_lead" ||
						v == "unassigned") {
						return nil, []error{fmt.Errorf("assigneeType needs to be one of project_default, component_lead, project_lead or unassigned")}
					}
					return nil, nil
				},
				DiffSuppressFunc: caseInsensitiveSuppressFunc,
				Optional:         true,
				Default:          "project_default",
			},

			"lead": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"project_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}

}

func resourceComponentCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.jiraClient

	componentOptions := &jira.CreateComponentOptions{}

	componentOptions.Name = d.Get("name").(string)
	componentOptions.Description = d.Get("description").(string)
	componentOptions.Project = d.Get("project_key").(string)
	componentOptions.AssigneeType = strings.ToUpper(d.Get("assignee_type").(string))
	componentOptions.LeadUserName = d.Get("lead").(string)
	component, _, err := client.Component.Create(componentOptions)
	if err != nil {
		return err
	}
	d.SetId(component.ID)
	return resourceComponentRead(d, m)
}

func resourceComponentRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	id := d.Id()
	urlStr := fmt.Sprintf("%s/%s", componentAPIEndpoint, id)

	component := &jira.ProjectComponent{}

	err := request(config.jiraClient, "GET", urlStr, nil, component)
	if err != nil {
		if errors.Is(err, ResourceNotFoundError) {
			d.SetId("")
			return nil
		}
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(component.ID)
	d.Set("name", component.Name)
	d.Set("description", component.Description)
	d.Set("project_key", component.Project)
	d.Set("assignee_type", component.AssigneeType)
	d.Set("lead", component.Lead.Name)

	return nil
}

func resourceComponentUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	id := d.Id()
	urlStr := fmt.Sprintf("%s/%s", componentAPIEndpoint, id)

	componentOptions := &jira.CreateComponentOptions{}

	componentOptions.Name = d.Get("name").(string)
	componentOptions.Description = d.Get("description").(string)
	componentOptions.Project = d.Get("project_key").(string)
	componentOptions.AssigneeType = strings.ToUpper(d.Get("assignee_type").(string))
	componentOptions.LeadUserName = d.Get("lead").(string)

	err := request(config.jiraClient, "PUT", urlStr, componentOptions, nil)

	if err != nil {
		return err
	}

	return resourceComponentRead(d, m)
}

func resourceComponentDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	id := d.Id()
	urlStr := fmt.Sprintf("%s/%s", componentAPIEndpoint, id)

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return err
	}

	return nil
}
