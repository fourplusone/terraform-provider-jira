resource "jira_permission_scheme" "editors" {
  name        = "Project Editors"
  description = "Permission scheme for project editors"
}

// Assign a permission scheme to a project
resource "jira_project_permission_scheme" "my_project" {
  project_key = "MYPROJ"
  scheme_id   = jira_permission_scheme.editors.id
}
