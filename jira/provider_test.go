package jira

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"jira": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("JIRA_URL"); v == "" {
		t.Fatal("JIRA_URL must be set for acceptance tests")
	}

	if v := os.Getenv("JIRA_USER"); v == "" {
		t.Fatal("JIRA_USER must be set for acceptance tests")
	}

	if v := os.Getenv("JIRA_PASSWORD"); v == "" {
		t.Fatal("JIRA_PASSWORD must be set for acceptance tests")
	}
}
