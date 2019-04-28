# terraform-provider-jira

[![CircleCI](https://circleci.com/gh/fourplusone/terraform-provider-jira.svg?style=svg)](https://circleci.com/gh/fourplusone/terraform-provider-jira)

Terraform Provider for managing JIRA. 

## Data Sources

- Issue Keys from JQL

## Resources

- Comments
- Filters & Filter Permissions
- Groups
- Group Memberships
- Issues
- Issue Links
- Issue Types
- Issue Link Types
- Projects
- Project Categories
- Project Roles
- Roles
- Webhooks

This can be used to interlink infrastructure management with JIRA issues closely.

![terraform-provider-jira demo](./images/terraform-provider-jira.gif)

## Install

* Download `terraform-provider-jira` binary from [Github](https://github.com/fourplusone/terraform-provider-jira/releases)
* Unzip the zip file
* Then move `terraform-provider-jira` binary to `$HOME/.terraform.d/plugins` directory

```bash
mkdir -p $HOME/.terraform.d/plugins
mv terraform-provider-jira $HOME/.terraform.d/plugins/terraform-provider-jira

```

* Run `terraform init` in your terraform project

```bash
terraform init
```

Used to initialize the plugin

## Example Usage

Set JIRA URL, Username and Password using environment variables

```bash
export JIRA_URL=http://localhost:8080
export JIRA_USER=username
export JIRA_PASSWORD=password
```

Create terraform config file

```hcl
// The types will be globally available in JIRA
resource "jira_issue_type" "task" {
  description = "A Task.",
  name = "Task",
  avatar_id = 10318
}

resource "jira_issue_link_type" "blocks" {
  name = "Blocks"
  inward = "is blocked by"
  outward = "blocks"
}

resource "jira_issue" "example" {
  issue_type  = "${jira_issue_type.task.name}"
  project_key = "PROJ"
  summary     = "Created using Terraform"

  // description is optional  
  description = "This is a test issue" 

  // (optional) Instead of deleting the issue, perform this transition 
  delete_transition = 21

  // (optional) Make sure, the issue is in the desired state
  // using state_transition
  state = 10000
  state_transition = 31 
}

resource "jira_comment" "example_comment" {
  body = "Commented using terraform"
  issue_key = "${jira_issue.example.issue_key}"
}

resource "jira_issue" "another_example" {
  issue_type  = "${jira_issue_type.task.name}"
  summary     = "Also Created using Terraform"
  project_key = "PROJ"
}

resource "jira_issue_link" "linked" {
  inward_key = "${jira_issue.example.issue_key}"
  outward_key = "${jira_issue.another_example.issue_key}"
  link_type = "${jira_issue_link_type.blocks.id}"
}

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
}

resource "jira_project_category" "categroy" {
  name = "Managed"
  description = "Managed Projects"
}

resource "jira_project" "project_a" {
  key = "TRF"
  name = "Terraform"
  project_type_key = "business"
  project_template_key = "com.atlassian.jira-core-project-templates:jira-core-project-management"
  lead = "bot"
  permission_scheme = 10400
  notification_scheme = 10300
  category_id = "${jira_project_category.category.id}"
}


// Create a group named "Terraform Managed"
resource "jira_group" "tf_group" {
  name = "Terraform Managed"
}

// User "bot" will be a Member of "Terraform Managed"

resource "jira_group_membership" "gm_1" {
  username = "bot"
  group = "${jira_group.tf_group.name}"
}

resource "jira_role" "role" {
  name = "Project Manager"
  description = "The Project Managers"
}

// Grant Project Access to user "bot" 
resource "jira_project_membership" "member" {
  project_key = "${jira_project.project_a.key}"
  role_id = "${jira_role.role.id}"
  username = "bot"
}

// Grant Project Access to group "bot" 
resource "jira_project_membership" "group_member" {
  project_key = "${jira_project.project_a.key}"
  role_id = "${jira_role.role.id}"
  group = "${jira_group.tf_group.name}"
}


resource "jira_webhook" "demo_hook" {
  name = "Terraform Hook"
  url = "https://demohook"
  jql = "project = PROJ"
  
  // See https://developer.atlassian.com/server/jira/platform/webhooks/ for supported events
  events = ["jira:issue_created"]
}


data "jira_jql" "issues" {
  jql = "project = ${jira_project.project_a.key} ORDER BY key ASC"
}

```

Run `terraform init`

```bash
terraform init
```

```bash

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

Run `terraform plan`

```bash
terraform plan
```

Check if the terraform plan looks good

Run `terraform apply`

```bash
terraform apply
```

## Rationale
Working in Operations engineering organizations infrastructure is often driven by tickets. Why not track infrastructure using tickets but this time we will use code.
This just showcases that you can pretty much Terraform anything!


## Credits

- Anubhav Mishra (anubhavmishra)
- Matthias Bartelmess (fourplusone)