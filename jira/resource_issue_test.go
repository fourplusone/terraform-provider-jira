package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccJiraIssue_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "jira_issue.example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraIssueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraIssueConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraIssueExists(resourceName),
					testAccCheckJiraIssueHasLabels(resourceName),
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

func testAccCheckJiraIssueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_issue" {
			continue
		}
		id := rs.Primary.ID

		_, resp, _ := client.Issue.Get(id, nil)

		if resp.StatusCode != 404 {
			return fmt.Errorf("Issue %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func testAccCheckJiraIssueExists(n string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient
		_, resp, _ := client.Issue.Get(rs.Primary.ID, nil)

		if resp.StatusCode != 200 {
			return fmt.Errorf("Issue %q does not exists", rs.Primary.ID)
		}
		return nil
	}

}

func testAccCheckJiraIssueHasLabels(n string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient
		issue, resp, _ := client.Issue.Get(rs.Primary.ID, nil)

		if resp.StatusCode != 200 {
			return fmt.Errorf("Issue %q does not exists", rs.Primary.ID)
		}

		labels := issue.Fields.Labels

		if labels == nil || len(labels) == 0 {
			return fmt.Errorf("Issue %q does not have any labels", rs.Primary.ID)
		}

		return nil
	}

}

func testAccJiraIssueConfig(rInt int) string {
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

resource "jira_issue" "example" {
	issue_type    = "Task"
	project_key   = "${jira_project.foo.key}"
	summary       = "Created using Terraform"
	labels        = ["label1", "label2", "label3", "label4"]
}
`, rInt, rInt, rInt%100000)
}
