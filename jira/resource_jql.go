package jira

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// resourceComment is used to define a JIRA comment
func resourceJQL() *schema.Resource {
	return &schema.Resource{
		Read: resourceJQLRead,

		Schema: map[string]*schema.Schema{
			"jql": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_keys": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceJQLRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	jql := d.Get("jql").(string)

	var issueKeys []string

	handler := func(i jira.Issue) error {
		issueKeys = append(issueKeys, i.Key)
		return nil
	}

	err := config.jiraClient.Issue.SearchPages(jql, nil, handler)
	if err != nil {
		return errors.Wrapf(err, "searching jira issue failed")
	}

	d.SetId(jql)
	d.Set("issue_keys", issueKeys)

	return nil
}
