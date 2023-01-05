resource "jira_issue" "example" {
  issue_type  = "${jira_issue_type.task.name}"
  project_key = "PROJ"
  summary     = "Created using Terraform"
}

resource "jira_comment" "example_comment" {
  body = "Commented using terraform"
  issue_key = "${jira_issue.example.issue_key}"
}
