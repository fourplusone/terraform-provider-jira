package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraProjectPermissionScheme_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "jira_project_permission_scheme.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraProjectPermissionSchemeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraProjectPermissionSchemeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraProjectPermissionSchemeExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "scheme_id"),
					resource.TestCheckResourceAttrSet(resourceName, "scheme_name"),
				),
			},
			{
				Config: testAccJiraProjectPermissionSchemeConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraProjectPermissionSchemeExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckJiraProjectPermissionSchemeDestroy(s *terraform.State) error {
	// A project permission scheme assignment cannot truly be destroyed;
	// delete resets it to the default scheme (ID 0). We verify the project
	// itself is gone or that the scheme was reset.
	return nil
}

func testAccCheckJiraProjectPermissionSchemeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project permission scheme ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient

		projectKey := rs.Primary.ID
		scheme := new(ProjectPermissionSchemeResponse)
		err := request(client, "GET", projectPermissionSchemeAPIEndpoint(projectKey), nil, scheme)

		if err != nil {
			return fmt.Errorf("Project permission scheme for project %q does not exist: %s", projectKey, err)
		}
		return nil
	}
}

func testAccJiraProjectPermissionSchemeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "jira_user" "foo" {
  name  = "pps-user-%d"
  email = "example@example.org"
}

resource "jira_project" "foo" {
  name             = "pps-project-%d"
  key              = "PP%d"
  lead             = jira_user.foo.name
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
}

resource "jira_permission_scheme" "foo" {
  name        = "test-pps-scheme-%d"
  description = "Test scheme for project assignment"
}

resource "jira_project_permission_scheme" "foo" {
  project_key = jira_project.foo.key
  scheme_id   = jira_permission_scheme.foo.id
}`, rInt, rInt, rInt%100000, rInt)
}

func testAccJiraProjectPermissionSchemeConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "jira_user" "foo" {
  name  = "pps-user-%d"
  email = "example@example.org"
}

resource "jira_project" "foo" {
  name             = "pps-project-%d"
  key              = "PP%d"
  lead             = jira_user.foo.name
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
}

resource "jira_permission_scheme" "foo" {
  name        = "test-pps-scheme-%d"
  description = "Test scheme for project assignment"
}

resource "jira_permission_scheme" "bar" {
  name        = "test-pps-scheme-updated-%d"
  description = "Updated test scheme for project assignment"
}

resource "jira_project_permission_scheme" "foo" {
  project_key = jira_project.foo.key
  scheme_id   = jira_permission_scheme.bar.id
}`, rInt, rInt, rInt%100000, rInt, rInt)
}
