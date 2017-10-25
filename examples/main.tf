resource "jira_issue" "example" {
  assignee    = "anubhavmishra"
  reporter    = "anubhavmishra"

  issue_type  = "Task"

  // description is optional  
  description = "This is a test issue"
  summary     = "Created using Terraform"

  project_key = "PROJ"
}
