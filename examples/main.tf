data "jira_field" "epic_name" {
  name = "Epic Name"
}

resource "jira_issue" "example" {
  assignee = "anubhavmishra"
  reporter = "anubhavmishra"

  issue_type = "Task"

  // description is optional
  description = "This is a test issue"
  summary     = "Created using Terraform"

  project_key = "PROJ"
}

resource "jira_issue" "example_epic" {
  assignee = "anubhavmishra"
  reporter = "anubhavmishra"

  issue_type = "Epic"

  // description is optional
  description = "This is an epic description"
  summary     = "This is an epic summary"

  labels = ["one", "two", "buckle-my-shoe"]

  // System and custom fields are optional; see the field data source to reference internal JIRA field IDs by name
  fields = {
    (data.jira_field.epic_name.id) = "Example epic name"
  }

  project_key = "PROJ"
}
