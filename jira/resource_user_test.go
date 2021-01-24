package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraUser_basic(t *testing.T) {
	// var group gitlab.User
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraUserConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraUserExists("jira_user.foo"),
				),
			},
		},
	})
}

func testAccCheckJiraUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_user" {
			continue
		}
		id := rs.Primary.ID

		_, resp, _ := getUserByKey(client, id)

		if resp.StatusCode != 404 {
			return fmt.Errorf("User %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func testAccCheckJiraUserExists(n string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient
		_, resp, _ := getUserByKey(client, rs.Primary.ID)

		if resp.StatusCode != 200 {
			return fmt.Errorf("User %q does not exists", rs.Primary.ID)
		}
		return nil
	}

}

func testAccJiraUserConfig(rInt int) string {
	return fmt.Sprintf(`
resource "jira_user" "foo" {
  name = "foo-name-%d"
  email = "example@example.org"
}`, rInt)
}
