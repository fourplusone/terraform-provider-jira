package jira

import (
	"log"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// FieldRequest The struct sent to the JIRA instance to create a new Project
type FieldRequest struct {
	Name                string `json:"name,omitempty" structs:"name,omitempty"`
	Description         string `json:"description,omitempty" structs:"description,omitempty"`
	Type                string `json:"type,omitempty" structs:"type,omitempty"`
	SearcherKey         string `json:"searcherKey,omitempty" structs:"searcherKey,omitempty"`
}

// resourceCustomField is used to define a JIRA custom field
func resourceCustomField() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomFieldCreate,
		Read:   resourceCustomFieldRead,
		Delete: resourceCustomFieldDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"searcher_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getCustomFieldById(client *jira.Client, id string) (*jira.Field, *jira.Response, error) {
	fields, resp, err := client.Field.GetList()
	if fields == nil {
		return nil, nil, err
	}

	for _, field := range fields {
		if field.ID == id {
			return &field, resp, nil
		}
	}

	return nil, resp, errors.Errorf("Custom Field not found")
}

// resourceCustomFieldRead reads custom field details using jira api
func resourceCustomFieldRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	field, _, err := getCustomFieldById(config.jiraClient, d.Id())
	if err != nil {
		return errors.Wrap(err, "getting jira field failed")
	}

	log.Printf("Read custom field (id=%s)", field.ID)

	d.Set("name", field.Name)
	return nil
}

// resourceCustomFieldCreate creates a new jira custom field using the jira api
func resourceCustomFieldCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	field := &FieldRequest{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Type:                d.Get("type").(string),
		SearcherKey:         d.Get("searcher_key").(string),
	}

	returnedField := new(jira.Field)

	err := request(config.jiraClient, "POST", fieldAPIEndpoint, field, returnedField)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	log.Printf("Created new Custom Field: %s", returnedField.ID)

	d.SetId(returnedField.ID)

	for err != nil {
		time.Sleep(200 * time.Millisecond)
		err = resourceCustomFieldRead(d, m)
	}

	return err
}

func resourceCustomFieldDelete(d *schema.ResourceData, m interface{}) error {
	return errors.Errorf("There is no way to delete a custom field via REST API")
}
