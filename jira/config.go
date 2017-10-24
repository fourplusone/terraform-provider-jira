package jira

import (
	"log"
	"os"

	jira "github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

type Config struct {
	jiraClient *jira.Client
}

func (c *Config) createAndAuthenticateClient() error {
	log.Printf("[INFO] creating jira client using environment variables")
	jiraClient, err := jira.NewClient(nil, os.Getenv("JIRA_URL"))
	if err != nil {
		return errors.Wrap(err, "creating jira client failed")
	}
	jiraClient.Authentication.SetBasicAuth(os.Getenv("JIRA_USER"), os.Getenv("JIRA_PASSWORD"))

	c.jiraClient = jiraClient

	return nil
}
