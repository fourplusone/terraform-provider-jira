resource "jira_issue_link_type" "blocks" {
  name = "Blocks"
  inward = "is blocked by"
  outward = "blocks"
}
