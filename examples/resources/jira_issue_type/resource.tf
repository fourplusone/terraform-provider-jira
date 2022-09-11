// The types will be globally available in JIRA
resource "jira_issue_type" "task" {
  description = "A Task."
  name = "Task"
  avatar_id = 10318
}
