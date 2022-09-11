package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var (
	fieldsCache []jira.Field
)

// JIRA field
func resourceField() *schema.Resource {
	return &schema.Resource{
		Read: resourceFieldRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"clause_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"navigable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"searchable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceFieldRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	name := d.Get("name").(string)

	if fieldsCache == nil || len(fieldsCache) == 0 {
		fields, _, err := config.jiraClient.Field.GetList()
		if err != nil {
			return errors.Wrapf(err, "fetching jira fields failed")
		}
		fieldsCache = fields
	}

	field := findFieldByName(fieldsCache, name)
	if field == nil {
		return errors.New(fmt.Sprintf("field with name '%s' not found", name))
	}

	d.SetId(field.ID)
	d.Set("clause_names", field.ClauseNames)
	d.Set("custom", field.Custom)
	d.Set("id", field.ID)
	d.Set("key", field.Key)
	d.Set("navigable", field.Navigable)
	d.Set("searchable", field.Searchable)

	return nil
}

func findFieldByName(fields []jira.Field, name string) *jira.Field {
	for _, field := range fields {
		if field.Name == name {
			return &field
		}
	}
	return nil
}
