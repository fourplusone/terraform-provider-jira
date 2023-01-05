data "jira_field" "epic_link" {
  name = "Epic Link"
}

resource "jira_issue" "custom_fields_example" {
  issue_type  = "Task"
  summary     = "Also Created using Terraform"
  fields      = {
    (jira_field.epic_link.id) = jira_issue.example_epic.issue_key
  }
  project_key = "PROJ"
}
