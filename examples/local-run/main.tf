terraform {
  required_providers {
    jira = {
      version = "0.1"
      source  = "idealo.de/pt/jira"
    }
  }
}

provider "jira" {
  url = "http://localhost:8080"
}

resource "jira_user" "foo" {
  name = "project-user-1"
  email = "example@example.org"
}

resource "jira_project" "foo" {
  name = "foo-name"
  key = "PX1000"
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
  lead = jira_user.foo.name
}

resource "jira_issue" "example" {
  issue_type = "Task"

  // description is optional
  description = "This is a test issue that's part of an epic"
  summary     = "Created using Terraform"
  labels      = ["label1", "label2"]
  state = 10001
  state_transitions = {
    10000 = jsonencode(["21"])
  }
  delete_transitions = {
    10001 = jsonencode(["51"])
  }

  project_key   = jira_project.foo.key
}