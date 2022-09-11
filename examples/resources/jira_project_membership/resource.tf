
resource "jira_role" "role" {
  name = "Project Manager"
  description = "The Project Managers"
}

resource "jira_group" "tf_group" {
  name = "Terraform Managed"
}

// Grant Project Access to user "bot" 
resource "jira_project_membership" "member" {
  project_key = "TRF"
  role_id = "${jira_role.role.id}"
  username = "bot"
}

// Grant Project Access to group "bot" 
resource "jira_project_membership" "group_member" {
  project_key = "TRF"
  role_id = "${jira_role.role.id}"
  group = "${jira_group.tf_group.name}"
}
