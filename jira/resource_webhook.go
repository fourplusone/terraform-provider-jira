package jira

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// WebhookFilter represents the JQL Filter for Webhook Events
type WebhookFilter struct {
	JQL string `json:"issue-related-events-section,omitempty"`
}

// Webhook represents a JIRA Webhook
type Webhook struct {
	Self        string        `json:"self,omitempty" structs:"self,omitempty"`
	Name        string        `json:"name,omitempty" structs:"name,omitempty"`
	URL         string        `json:"url,omitempty" structs:"url,omitempty"`
	Events      []string      `json:"events,omitempty" structs:"events,omitempty"`
	Filters     WebhookFilter `json:"filters"`
	ExcludeBody bool          `json:"excludeBody,omitempty" structs:"excludeBody,omitempty"`
}

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Update: resourceWebhookUpdate,
		Delete: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"jql": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"events": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_body": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func setWebhook(w *Webhook, d *schema.ResourceData) {
	w.Name = d.Get("name").(string)
	w.URL = d.Get("url").(string)

	resourceEvents := d.Get("events").([]interface{})
	events := make([]string, 0, len(resourceEvents))

	for _, event := range resourceEvents {
		events = append(events, event.(string))
	}

	w.Events = events
	w.Filters.JQL = d.Get("jql").(string)
	w.ExcludeBody = d.Get("exclude_body").(bool)
}

func setWebhookResource(w *Webhook, d *schema.ResourceData) {
	components := strings.Split(w.Self, "/")
	ID := components[len(components)-1]

	d.SetId(ID)
	d.Set("name", w.Name)
	d.Set("url", w.URL)
	d.Set("events", w.Events)
	d.Set("exclude_body", w.ExcludeBody)
	d.Set("jql", w.Filters.JQL)
}

// resourceWebhookCreate creates a new jira issue using the jira api
func resourceWebhookCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	webhook := new(Webhook)
	returnedWebhook := new(Webhook)

	setWebhook(webhook, d)

	err := request(config.jiraClient, "POST", webhookAPIEndpoint, webhook, returnedWebhook)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setWebhookResource(returnedWebhook, d)

	return resourceWebhookRead(d, m)
}

// resourceWebhookRead reads issue details using jira api
func resourceWebhookRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", webhookAPIEndpoint, d.Id())

	webhook := new(Webhook)
	err := request(config.jiraClient, "GET", urlStr, nil, webhook)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setWebhookResource(webhook, d)

	return nil
}

// resourceWebhookUpdate updates jira issue using jira api
func resourceWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	webhook := new(Webhook)

	setWebhook(webhook, d)

	urlStr := fmt.Sprintf("%s/%s", webhookAPIEndpoint, d.Id())
	returnedWebhook := new(Webhook)

	err := request(config.jiraClient, "PUT", urlStr, webhook, returnedWebhook)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceWebhookRead(d, m)
}

// resourceWebhookDelete deletes jira issue using the jira api
func resourceWebhookDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", webhookAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
