package jira

import (
	"fmt"
	"net/url"
	"strconv"
	"log"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func resourceIssueTypeScheme() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueTypeSchemeCreate,
		Read:   resourceIssueTypeSchemeRead,
		Delete: resourceIssueTypeSchemeDelete,
		Update: resourceIssueTypeSchemeUpdate,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"issue_type_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{ Type: schema.TypeString, },
				Required: true,
				ForceNew: true,
			},
			"project_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{ Type: schema.TypeString, },
				Required: true,
				ForceNew: true,
			},
			"default_issue_type_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

type IssueTypeSchemeItem struct {
	IssueTypeSchemeId  string  `json:"issueTypeSchemeId"`
	IssueTypeId        string  `json:"issueTypeId"`
}

type GetIssueTypeSchemeItemsResult struct {
	StartAt    int                    `json:"startAt"`
	MaxResults int                    `json:"maxResults"`
	Total      int                    `json:"total"`
	IsLast     bool                   `json:"isLast"`
	Values     []IssueTypeSchemeItem  `json:"values"`
}

func listIssueTypeSchemeItems(client *jira.Client, issueTypeSchemeId string, items *schema.Set) error {
	relativeURL, _ := url.Parse(issueTypeSchemeMappingAPIEndpoint)
	query := relativeURL.Query()

	resp := new(GetIssueTypeSchemeItemsResult)
	resp.IsLast = false
	resp.StartAt = 0
	resp.MaxResults = 200

	for !resp.IsLast {
		query.Set("startAt", strconv.Itoa(resp.StartAt + len(resp.Values)))
		query.Set("maxResults", strconv.Itoa(resp.MaxResults))
		query.Set("issueTypeSchemeId", issueTypeSchemeId)
		relativeURL.RawQuery = query.Encode()

		err := request(client, "GET", relativeURL.String(), nil, resp)
		if err != nil {
			return err
		}

		for _, v := range resp.Values {
			items.Add(v.IssueTypeId)
		}
	}

	return nil
}

type AddProjectAssociationRequest struct {
	IdsOrKeys []string `json:"idsOrKeys" structs:"idsOrKeys"`
}

// https://docs.atlassian.com/software/jira/docs/api/REST/8.10.1/#api/2/issuetypescheme-setProjectAssociationsForScheme
func setProjectAssociationToScheme(client *jira.Client, schemeId string, items *schema.Set) error {
	ep := fmt.Sprintf("%s/%s/associations", issueTypeSchemeAPIEndpoint, schemeId)
	req := new(AddProjectAssociationRequest)
	req.IdsOrKeys = make([]string, items.Len())
	for i, idOrKey := range items.List() {
		req.IdsOrKeys[i] = idOrKey.(string)
	}

	return request(client, "PUT", ep, req, nil)
}

func getProjectAssociationFromScheme(client *jira.Client, schemeId string, items *schema.Set) error {
	ep := fmt.Sprintf("%s/%s/associations", issueTypeSchemeAPIEndpoint, schemeId)
	resp := new(jira.ProjectList)
	err := request(client, "GET", ep, nil, resp)
	if err != nil {
		return err
	}

	for _, v := range *resp {
		items.Add(v.ID)
	}

	return nil
}

type IssueTypeScheme struct {
	Name                string `json:"name,omitempty" structs:"name,omitempty"`
	Description         string `json:"description,omitempty" structs:"description,omitempty"`
	IssueTypeIds        []string `json:"issueTypeIds,omitempty" structs:"issueTypeIds,omitempty"`
	DefaultIssueTypeId  string `json:"defaultIssueTypeId,omitempty" structs:"defaultIssueTypeId,omitempty"`
}

type GetIssueTypeSchemeResult struct {
	StartAt    int                       `json:"startAt"`
	MaxResults int                       `json:"maxResults"`
	Total      int                       `json:"total"`
	Values     []IssueTypeScheme         `json:"values"`
}

func resourceIssueTypeSchemeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	relativeURL, _ := url.Parse(issueTypeSchemeAPIEndpoint)
	query := relativeURL.Query()
	query.Set("id", d.Id())
	relativeURL.RawQuery = query.Encode()

	resp := new(GetIssueTypeSchemeResult)

	err := request(config.jiraClient, "GET", relativeURL.String(), nil, resp)
	if err != nil || resp.Total != 1 {
		return errors.Wrap(err, "getting JIRA Issue Type Scheme failed")
	}

	scheme := resp.Values[0]

	d.Set("name", scheme.Name)
	d.Set("description", scheme.Description)
	d.Set("default_issue_type_id", scheme.DefaultIssueTypeId)

	items := schema.NewSet(schema.HashString, make([]interface{}, 0))
	err = listIssueTypeSchemeItems(config.jiraClient, d.Id(), items)
	if err != nil {
		return errors.Wrap(err, "getting JIRA Issue Type Scheme items failed")
	}

	log.Printf("Read Issue Type Scheme (%s) with %d items", scheme.Name, items.Len())

	// TypeSet members needs to be checked
	if err := d.Set("issue_type_ids", items); err != nil {
		return errors.Wrap(err, "Error setting issue type ids")
	}

	items = schema.NewSet(schema.HashString, make([]interface{}, 0))
	err = getProjectAssociationFromScheme(config.jiraClient, d.Id(), items)
	if err != nil {
		return errors.Wrap(err, "getting JIRA Issue Type Scheme associated projects failed")
	}

	log.Printf("Read Issue Type Scheme (%s) with %d projects", scheme.Name, items.Len())

	// TypeSet members needs to be checked
	if err := d.Set("project_ids", items); err != nil {
		return errors.Wrap(err, "Error setting project ids")
	}

	return nil
}

type CreateIssueTypeSchemeResponse struct {
	ID string `json:"id" structs:"id"`
}

func resourceIssueTypeSchemeCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	scheme := new(IssueTypeScheme)
	scheme.Name = d.Get("name").(string)
	scheme.Description = d.Get("description").(string)
	scheme.DefaultIssueTypeId = d.Get("default_issue_type_id").(string)

	items := d.Get("issue_type_ids").(*schema.Set).List()
	scheme.IssueTypeIds = make([]string, len(items))
	for i, v := range items {
		scheme.IssueTypeIds[i] = v.(string)
	}

	resp := new(CreateIssueTypeSchemeResponse)

	err := request(config.jiraClient, "POST", issueTypeSchemeAPIEndpoint, scheme, resp)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	if resp.ID == "" {
		return errors.Errorf("New Issue Type Scheme ID was NOT returned %+v", resp)
	}

	log.Printf("Created new Issue Type Scheme: %s", resp.ID)

	d.SetId(resp.ID)

	err = setProjectAssociationToScheme(config.jiraClient, resp.ID, d.Get("project_ids").(*schema.Set))
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceIssueTypeSchemeRead(d, m)
}

func resourceIssueTypeSchemeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueTypeSchemeAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}

type UpdateIssueTypeScheme struct {
	Name                string `json:"name,omitempty" structs:"name,omitempty"`
	Description         string `json:"description,omitempty" structs:"description,omitempty"`
	DefaultIssueTypeId  string `json:"defaultIssueTypeId,omitempty" structs:"defaultIssueTypeId,omitempty"`
}

func resourceIssueTypeSchemeUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", issueTypeSchemeAPIEndpoint, d.Id())

	scheme := new(UpdateIssueTypeScheme)
	scheme.Name = d.Get("name").(string)
	scheme.Description = d.Get("description").(string)
	scheme.DefaultIssueTypeId = d.Get("default_issue_type_id").(string)

	err := request(config.jiraClient, "PUT", urlStr, scheme, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	err = setProjectAssociationToScheme(config.jiraClient, d.Id(), d.Get("project_ids").(*schema.Set))
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceIssueTypeSchemeRead(d, m)
}
