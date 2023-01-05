resource "jira_issue" "example" {
  issue_type  = "Task"
  project_key = "PROJ"
  summary     = "Created using Terraform"
}

resource "jira_issue" "another_example" {
  issue_type  = "Task"
  project_key = "PROJ"
  summary     = "Created using Terraform"
}

resource "jira_issue_link_type" "blocks" {
  name = "Blocks"
  inward = "is blocked by"
  outward = "blocks"
}

resource "jira_issue_link" "linked" {
  inward_key = "${jira_issue.example.issue_key}"
  outward_key = "${jira_issue.another_example.issue_key}"
  link_type = "${jira_issue_link_type.blocks.id}"
}
