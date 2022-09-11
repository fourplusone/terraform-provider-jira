resource "jira_webhook" "demo_hook" {
  name = "Terraform Hook"
  url = "https://demohook"
  jql = "project = PROJ"
  
  // See https://developer.atlassian.com/server/jira/platform/webhooks/ for supported events
  events = ["jira:issue_created"]
}
