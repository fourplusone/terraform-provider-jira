package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccJiraGroup_basic(t *testing.T) {
	// var group gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraGroupExists("jira_group.foo"),
				),
			},
		},
	})
}

func testAccCheckJiraGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_group" {
			continue
		}
		id := rs.Primary.ID

		_, resp, _ := client.Group.Get(id)

		if resp.StatusCode != 404 {
			return fmt.Errorf("Group %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func testAccCheckJiraGroupExists(n string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No group ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient
		_, resp, _ := client.Group.Get(rs.Primary.ID)

		if resp.StatusCode != 200 {
			return fmt.Errorf("Group %q does not exists", rs.Primary.ID)
		}
		return nil
	}

}

func testAccJiraGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "jira_group" "foo" {
  name = "foo-name-%d"
}`, rInt)
}
