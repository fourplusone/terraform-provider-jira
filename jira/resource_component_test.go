package jira

import (
	"errors"
	"fmt"
	"testing"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccJiraComponent_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "jira_component.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJiraComponentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraComponentConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJiraComponentExists(resourceName),
				),
			},
			{
				Config: testAccJiraComponentConfigB(rInt),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckJiraComponentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).jiraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jira_component" {
			continue
		}
		id := rs.Primary.ID

		urlStr := fmt.Sprintf("%s/%s", componentAPIEndpoint, id)
		component := &jira.ProjectComponent{}
		err := request(client, "GET", urlStr, nil, component)

		if !errors.Is(err, ResourceNotFoundError) {
			return fmt.Errorf("Component %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func testAccCheckJiraComponentExists(n string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No component ID is set")
		}

		id := rs.Primary.ID

		client := testAccProvider.Meta().(*Config).jiraClient

		urlStr := fmt.Sprintf("%s/%s", componentAPIEndpoint, id)
		component := &jira.ProjectComponent{}
		err := request(client, "GET", urlStr, nil, component)

		if errors.Is(err, ResourceNotFoundError) {
			return fmt.Errorf("Component %q does not exists", id)
		}
		return err
	}

}

func testEnvConfig(rInt int) string {
	return fmt.Sprintf(`

	resource "jira_user" "u1" {
		name = "u1-name-%d"
		email = "example@example.org"
	}

	resource "jira_user" "u2" {
		name = "u2-name-%d"
		email = "example@example.org"
	}

	resource "jira_project" "p1" {
		name = "p1-name-%d"
		lead = jira_user.u1.name
		permission_scheme = 10000
		key = "PX%d"
		project_type_key = "software"
		project_template_key = "com.pyxis.greenhopper.jira:basic-software-development-template"
	}
	
	resource "jira_project" "p2" {
		name = "p2-name-%d"
		lead = jira_user.u2.name
		permission_scheme = 10000
		key = "PZ%d"
		project_type_key = "software"
		project_template_key = "com.pyxis.greenhopper.jira:basic-software-development-template"
	  }

	`, rInt, rInt, rInt, rInt%100000, rInt, rInt%100000)
}

func testAccJiraComponentConfig(rInt int) string {
	return fmt.Sprintf(`
	%s

	resource "jira_component" "foo" {
		name = "foo-component-%d"
		project_key = "${jira_project.p1.key}"
	}

	resource "jira_component" "foo-bar" {
		name = "foo-component-%d-2"
		description = "Sample Description"
		lead = "${jira_user.u1.name}"
		assignee_type = "component_lead"
		project_key = "${jira_project.p1.key}"
	}
	`, testEnvConfig(rInt), rInt, rInt)
}

func testAccJiraComponentConfigB(rInt int) string {
	return fmt.Sprintf(`
	%s
	
	resource "jira_component" "foo" {
		name = "My Component %d"
		project_key = "${jira_project.p1.key}"
	}

	resource "jira_component" "foo-bar" {
		name = "My Component %d 1"
		description = "Hello World"
		lead = "${jira_user.u2.name}"
		assignee_type = "component_lead"
		project_key = "${jira_project.p1.key}"
	}
	`, testEnvConfig(rInt), rInt, rInt)
}
