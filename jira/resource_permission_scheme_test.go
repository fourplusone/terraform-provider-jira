package jira

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJiraPermissionScheme_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "jira_permission_scheme.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraPermissionSchemeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraPermissionSchemeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraPermissionSchemeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-scheme-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "description", "Test permission scheme"),
				),
			},
			{
				Config: testAccJiraPermissionSchemeConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraPermissionSchemeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test-scheme-updated-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated test permission scheme"),
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

func testAccCheckJiraPermissionSchemeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_permission_scheme" {
			continue
		}

		urlStr := fmt.Sprintf("%s/%s", permissionSchemeAPIEndpoint, rs.Primary.ID)

		scheme := new(PermissionScheme)
		err := request(client, "GET", urlStr, nil, scheme)

		if err == nil {
			return fmt.Errorf("Permission scheme %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckJiraPermissionSchemeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No permission scheme ID is set")
		}

		client := testAccProvider.Meta().(*Config).jiraClient

		urlStr := fmt.Sprintf("%s/%s", permissionSchemeAPIEndpoint, rs.Primary.ID)

		scheme := new(PermissionScheme)
		err := request(client, "GET", urlStr, nil, scheme)

		if err != nil {
			return fmt.Errorf("Permission scheme %q does not exist: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccJiraPermissionSchemeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "jira_permission_scheme" "foo" {
  name        = "test-scheme-%d"
  description = "Test permission scheme"
}`, rInt)
}

func testAccJiraPermissionSchemeConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "jira_permission_scheme" "foo" {
  name        = "test-scheme-updated-%d"
  description = "Updated test permission scheme"
}`, rInt)
}
