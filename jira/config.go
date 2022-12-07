package jira

import (
	"log"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

type Config struct {
	jiraClient *jira.Client
}

type AddHeaderTransport struct {
	customAuthHeaderKey   string
	customAuthHeaderValue string
	Transport             http.RoundTripper
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(adt.customAuthHeaderKey, adt.customAuthHeaderValue)
	return adt.Transport.RoundTrip(req)
}

func (c *Config) createAndAuthenticateClient(d *schema.ResourceData) error {
	log.Printf("[INFO] creating jira client using environment variables")

	transport := AddHeaderTransport{}

	if d.Get("custom_auth_header_key") != nil && d.Get("custom_auth_header_value") != nil {
		transport = AddHeaderTransport{
			customAuthHeaderKey:   d.Get("custom_auth_header_key").(string),
			customAuthHeaderValue: d.Get("custom_auth_header_value").(string),
			Transport:             http.DefaultTransport,
		}
	}

	if d.Get("pat_token") != nil {
		tp := jira.BearerAuthTransport{
			Token:     d.Get("pat_token").(string),
			Transport: &transport,
		}

		jiraClient, err := jira.NewClient(tp.Client(), d.Get("url").(string))
		if err != nil {
			return errors.Wrap(err, "creating jira client failed")
		}

		c.jiraClient = jiraClient

	} else if d.Get("username") != nil && d.Get("password") != nil {
		tp := jira.BasicAuthTransport{
			Username:  d.Get("username").(string),
			Password:  d.Get("password").(string),
			Transport: &transport,
		}

		jiraClient, err := jira.NewClient(tp.Client(), d.Get("url").(string))
		if err != nil {
			return errors.Wrap(err, "creating jira client failed")
		}

		c.jiraClient = jiraClient
	}

	return nil
}
