resource "jira_filter" "filter" {
  name = "Simple Filter"
  jql = "project = PROJ"

  // Optional Fields
  description = "All Issues in PROJ"
  favourite = false

  // All Members of project with ID 13102
  permissions {
    type = "project"
    project_id = "13102"
  }

  // All Members of Group "Team A"
  permissions {
    type = "group"
    group_name = "Team A"
  }

  // Any authenticated user
  permissions {
    type = "authenticated"
  }
}
