package jira

import (
	"log"
	"net/http"
	"sync"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

type Config struct {
	jiraClient *jira.Client
	jiraLock   sync.Mutex
}

func (c *Config) createAndAuthenticateClient(d *schema.ResourceData) error {
	log.Printf("[INFO] creating jira client using environment variables")

	var httpClient *http.Client

	token, ok := d.GetOk("token")
	if ok {
		transport := jira.BearerAuthTransport{Token: token.(string)}
		httpClient = transport.Client()
	} else {
		transport := &jira.BasicAuthTransport{
			Username: d.Get("user").(string),
			Password: d.Get("password").(string),
		}
		httpClient = transport.Client()
	}

	jiraClient, err := jira.NewClient(httpClient, d.Get("url").(string))
	if err != nil {
		return errors.Wrap(err, "creating jira client failed")
	}

	c.jiraClient = jiraClient

	return nil
}
