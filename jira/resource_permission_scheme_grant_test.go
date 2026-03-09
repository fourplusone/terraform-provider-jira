package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraPermissionSchemeGrant_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "jira_permission_scheme_grant.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraPermissionSchemeGrantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraPermissionSchemeGrantConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraPermissionSchemeGrantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "permission", "BROWSE_PROJECTS"),
					resource.TestCheckResourceAttr(resourceName, "holder_type", "group"),
					resource.TestCheckResourceAttrSet(resourceName, "holder_parameter"),
				),
			},
		},
	})
}

func TestAccJiraPermissionSchemeGrant_projectLead(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "jira_permission_scheme_grant.lead"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraPermissionSchemeGrantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraPermissionSchemeGrantConfigProjectLead(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraPermissionSchemeGrantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "permission", "ADMINISTER_PROJECTS"),
					resource.TestCheckResourceAttr(resourceName, "holder_type", "projectLead"),
				),
			},
		},
	})
}

func testAccCheckJiraPermissionSchemeGrantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_permission_scheme_grant" {
			continue
		}

		schemeID := rs.Primary.Attributes["scheme_id"]
		urlStr := fmt.Sprintf("%s/%s", permissionSchemeGrantAPIEndpoint(schemeID), rs.Primary.ID)

		grant := new(PermissionGrant)
		err := request(client, "GET", urlStr, nil, grant)

		if err == nil {
			return fmt.Errorf("Permission scheme grant %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckJiraPermissionSchemeGrantExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No permission scheme grant ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient

		schemeID := rs.Primary.Attributes["scheme_id"]
		urlStr := fmt.Sprintf("%s/%s", permissionSchemeGrantAPIEndpoint(schemeID), rs.Primary.ID)

		grant := new(PermissionGrant)
		err := request(client, "GET", urlStr, nil, grant)

		if err != nil {
			return fmt.Errorf("Permission scheme grant %q does not exist: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccJiraPermissionSchemeGrantConfig(rInt int) string {
	return fmt.Sprintf(`
resource "jira_permission_scheme" "foo" {
  name        = "test-grant-scheme-%d"
  description = "Test permission scheme for grants"
}

resource "jira_group" "foo" {
  name = "test-grant-group-%d"
}

resource "jira_permission_scheme_grant" "foo" {
  scheme_id        = jira_permission_scheme.foo.id
  permission       = "BROWSE_PROJECTS"
  holder_type      = "group"
  holder_parameter = jira_group.foo.name
}`, rInt, rInt)
}

func testAccJiraPermissionSchemeGrantConfigProjectLead(rInt int) string {
	return fmt.Sprintf(`
resource "jira_permission_scheme" "lead" {
  name        = "test-lead-scheme-%d"
  description = "Test permission scheme for project lead grant"
}

resource "jira_permission_scheme_grant" "lead" {
  scheme_id   = jira_permission_scheme.lead.id
  permission  = "ADMINISTER_PROJECTS"
  holder_type = "projectLead"
}`, rInt)
}
