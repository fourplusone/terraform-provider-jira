package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccJiraProject_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraProjectExists("jira_project.foo"),
				),
			},
		},
	})
}

func TestAccJiraProject_shared(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraSharedProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraProjectExists("jira_project.foo_shared"),
				),
			},
		},
	})
}

func testAccCheckJiraProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_project" {
			continue
		}
		id := rs.Primary.ID

		_, resp, _ := client.Project.Get(id)

		if resp.StatusCode != 404 {
			return fmt.Errorf("Project %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func testAccCheckJiraProjectExists(n string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient
		_, resp, _ := client.Project.Get(rs.Primary.ID)

		if resp.StatusCode != 200 {
			return fmt.Errorf("Project %q does not exists", rs.Primary.ID)
		}
		return nil
	}

}

func testAccJiraProjectConfig(rInt int) string {
	return fmt.Sprintf(`

resource "jira_user" "foo" {
	name = "project-user-%d"
	email = "example@example.org"
}

resource "jira_project" "foo" {
  name = "foo-name-%d"
  key = "PX%d"
  lead = "${jira_user.foo.name}"
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
}`, rInt, rInt, rInt%100000)
}

func testAccJiraSharedProjectConfig(rInt int) string {
	return fmt.Sprintf(`

resource "jira_user" "foo" {
	name = "project-user-%d"
	email = "example@example.org"
}

resource "jira_project" "foo" {
  name = "foo-name-%d"
  key = "PX%d"
  lead = "${jira_user.foo.name}"
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
}

resource "jira_project" "foo_shared" {
	name = "foo-shared-%d"
	key = "PC%d"
	lead = "${jira_user.foo.name}"
	shared_configuration_project_id = "${jira_project.foo.project_id}"
}`, rInt, rInt, rInt%100000, rInt, rInt%100000)
}
