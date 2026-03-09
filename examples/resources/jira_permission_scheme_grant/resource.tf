resource "jira_permission_scheme" "editors" {
  name        = "Project Editors"
  description = "Permission scheme for project editors"
}

resource "jira_group" "developers" {
  name = "developers"
}

// Grant BROWSE_PROJECTS to a group
resource "jira_permission_scheme_grant" "browse" {
  scheme_id        = jira_permission_scheme.editors.id
  permission       = "BROWSE_PROJECTS"
  holder_type      = "group"
  holder_parameter = jira_group.developers.name
}

// Grant ADMINISTER_PROJECTS to the project lead
resource "jira_permission_scheme_grant" "admin" {
  scheme_id   = jira_permission_scheme.editors.id
  permission  = "ADMINISTER_PROJECTS"
  holder_type = "projectLead"
}

// Grant CREATE_ISSUES to a specific role
resource "jira_role" "developers_role" {
  name        = "Developers"
  description = "Development team role"
}

resource "jira_permission_scheme_grant" "create_issues" {
  scheme_id        = jira_permission_scheme.editors.id
  permission       = "CREATE_ISSUES"
  holder_type      = "projectRole"
  holder_parameter = jira_role.developers_role.id
}
