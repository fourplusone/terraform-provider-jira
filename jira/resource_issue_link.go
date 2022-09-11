package jira

import (
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceIssueLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueLinkCreate,
		Read:   resourceIssueLinkRead,
		Delete: resourceIssueLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"inward_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"outward_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"link_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceIssueLinkCreate creates a new jira issue using the jira api
func resourceIssueLinkCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issueLink := new(jira.IssueLink)

	issueLink.InwardIssue = &jira.Issue{Key: d.Get("inward_key").(string)}
	issueLink.OutwardIssue = &jira.Issue{Key: d.Get("outward_key").(string)}
	issueLink.Type = jira.IssueLinkType{ID: d.Get("link_type").(string)}

	resp, err := config.jiraClient.Issue.AddLink(issueLink)

	if err != nil {
		return errors.Wrap(err, "Creating Issue Link failed")
	}

	location, err := resp.Location()

	if err != nil {
		return errors.Wrap(err, "Creating Issue Link failed")
	}

	components := strings.Split(location.Path, "/")
	ID := components[len(components)-1]

	d.SetId(ID)

	return resourceIssueLinkRead(d, m)
}

// resourceIssueLinkRead reads issue details using jira api
func resourceIssueLinkRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueLinkAPIEndpoint, d.Id())
	issueLink := new(jira.IssueLink)

	err := request(config.jiraClient, "GET", urlStr, nil, issueLink)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.Set("inward_key", issueLink.InwardIssue.Key)
	d.Set("outward_key", issueLink.OutwardIssue.Key)
	d.Set("link_type", issueLink.Type.ID)

	return nil
}

// resourceIssueLinkDelete deletes jira issue using the jira api
func resourceIssueLinkDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueLinkAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
