// Create a group named "Terraform Managed"
resource "jira_group" "tf_group" {
  name = "Terraform Managed"
}

// User "bot" will be a Member of "Terraform Managed"
resource "jira_group_membership" "gm_1" {
  username = "bot"
  group = "${jira_group.tf_group.name}"
}
