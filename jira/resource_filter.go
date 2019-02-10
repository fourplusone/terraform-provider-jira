package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

// FilterRequest represents a Filter in Jira
type FilterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Jql         string `json:"jql"`
	Favourite   bool   `json:"favourite"`
}

// resourceFilter is used to define a JIRA Filter
func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceFilterCreate,
		Read:   resourceFilterRead,
		Update: resourceFilterUpdate,
		Delete: resourceFilterDelete,
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
			"jql": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"favourite": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func setFilter(w *FilterRequest, d *schema.ResourceData) {
	w.Name = d.Get("name").(string)
	w.Description = d.Get("description").(string)
	w.Jql = d.Get("jql").(string)
	w.Favourite = d.Get("favourite").(bool)
}

func setFilterResource(w *jira.Filter, d *schema.ResourceData) {
	d.SetId(w.ID)
	d.Set("name", w.Name)
	d.Set("description", w.Description)
	d.Set("jql", w.Jql)
	d.Set("favourite", w.Favourite)
}

// resourceFilterCreate creates a new jira filter using the jira api
func resourceFilterCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	filter := new(FilterRequest)
	returnedFilter := new(jira.Filter)
	setFilter(filter, d)

	err := request(config.jiraClient, "POST", filterAPIEndpoint, filter, returnedFilter)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setFilterResource(returnedFilter, d)

	return nil
}

// resourceFilterRead reads filter details using jira api
func resourceFilterRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", filterAPIEndpoint, d.Id())

	filter := new(jira.Filter)
	err := request(config.jiraClient, "GET", urlStr, nil, filter)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setFilterResource(filter, d)

	return nil
}

// resourceFilterUpdate updates jira filter using jira api
func resourceFilterUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	filter := new(FilterRequest)

	setFilter(filter, d)

	urlStr := fmt.Sprintf("%s/%s", filterAPIEndpoint, d.Id())
	returnedFilter := new(jira.Filter)

	err := request(config.jiraClient, "PUT", urlStr, filter, returnedFilter)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceFilterRead(d, m)
}

// resourceFilterDelete deletes jira filter using the jira api
func resourceFilterDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", filterAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
