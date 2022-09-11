data "jira_jql" "issues" {
  jql = "project = TRF ORDER BY key ASC"
}
