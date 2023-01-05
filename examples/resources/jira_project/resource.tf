resource "jira_project_category" "category" {
  name = "Managed"
  description = "Managed Projects"
}

resource "jira_project" "project_a" {
  key = "TRF"
  name = "Terraform"
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
  lead = "bot"
  // For JIRA Cloud use lead_account_id instead
  lead_account_id = "xxxxxx:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  permission_scheme = 10400
  notification_scheme = 10300
  category_id = "${jira_project_category.category.id}"
}


// Create a Project with a shared configuration
resource "jira_project" "project_shared" {
	key = "SHARED"
	name = "Project (with shared config)"
	lead = "bot"
	shared_configuration_project_id = "${jira_project.project_a.project_id}"
}
