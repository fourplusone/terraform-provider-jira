# terraform-provider-jira

[![Build & Tests](https://github.com/fourplusone/terraform-provider-jira/actions/workflows/go.yml/badge.svg)](https://github.com/fourplusone/terraform-provider-jira/actions/workflows/go.yml)

Terraform Provider for managing JIRA. (__[View on registry.terraform.io](https://registry.terraform.io/providers/fourplusone/jira/latest)__)

## Data Sources

- Issue Keys from JQL
- Custom Fields

## Resources

- Comments
- Components
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
- Users
- Webhooks

This can be used to interlink infrastructure management with JIRA issues closely.

![terraform-provider-jira demo](./images/terraform-provider-jira.gif)

## Compatibility 


| Version | Status       |  
|---------|--------------|
| Latest  | Tested in CI |
| 7.x     | Untested     |
| Cloud   | Untested     |   

## Install

## Terraform v0.13 or newer

Copy this code into yout terraform configuration file (for example `main.tf`)

```hcl
terraform {
  required_providers {
    jira = {
      source = "fourplusone/jira"
      version = "0.1.16"
    }
  }
}
```

Run `terraform init`

## Example Usage

Set JIRA URL, Username and Password using environment variables

```bash
export JIRA_URL=http://localhost:8080
export JIRA_USER=username
export JIRA_PASSWORD=password
```

It's also possible to use an API-Token from JIRA cloud. In this case, set 

```bash
export JIRA_URL=https://yourinstance.atlassian.net
export JIRA_USER=username@example.org
export JIRA_PASSWORD=<API-Key>
```

Create terraform config file

```hcl

provider "jira" {
  url = "https://myjira.atlassian.net" # Can also be set using the JIRA_URL environment variable
  # user = "xxxx"                      # Can also be set using the JIRA_USER environment variable
  # password = "xxxx"                  # Can also be set using the JIRA_PASSWORD environment variable
}

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

resource "jira_component" "example_component" {
  
  name = "Sample Component"
  project_key = "${jira_project.example.key}"
  
  // (optional) Description of the component
  description = "Sample Description" 
  
  // (optional) Component lead
  lead = "${jira_user.demo_user.name}"

  // (optional) assignee type. Can be one of project_default, component_lead, project_lead or unassigned.
	assignee_type = "component_lead"
}

resource "jira_issue" "another_example" {
  issue_type  = "${jira_issue_type.task.name}"
  summary     = "Also Created using Terraform"
  labels      = ["label1", "label2"]
  project_key = "PROJ"
}

data "jira_field" "epic_link" {
  name = "Epic Link"
}

resource "jira_issue" "custom_fields_example" {
  issue_type  = "${jira_issue_type.task.name}"
  summary     = "Also Created using Terraform"
  fields      = {
    (jira_field.epic_link.id) = jira_issue.example_epic.issue_key
  }
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

  // Any authenticated user
  permissions {
    type = "authenticated"
  }
}

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

resource "jira_user" "demo_user" {
  name = "bot"
  email = "bot@example.org"
  display_name = "The Bot"
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


resource "jira_custom_field" "L5" {
  name = "My custom field"
  type = "com.atlassian.jira.plugin.system.customfieldtypes:textfield"
  searcher_key = "com.atlassian.jira.plugin.system.customfieldtypes:textsearcher"
}


resource "jira_issue_type_scheme" "scheme" {
  name = "JIRA Task Force Scheme"
  description = "The issue types for projects engaged in SFS program"
  default_issue_type_id = jira_issue_type.task.id
  issue_type_ids = [jira_issue_type.task.id]
  project_ids = [jira_project.project_a.id]
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

## Building 

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:fourplusone/terraform-provider-jira
...
```

Enter the provider directory and run `make build` to build the provider

```sh
$ make build
```

## Testing
Testing requires a JIRA instance. To set up a temporary, local JIRA instance you can use the [Docker Atlas](https://github.com/fourplusone/docker-atlas) container:

```sh
$ git clone https://github.com/fourplusone/docker-atlas
$ cd docker-atlas
$ docker build -f jira.dockerfile . -t jira-docker
$ docker run --rm -t -p 2990:2990 jira-docker
```

Building and starting the JIRA instance will take several minutes.

To run the tests, run `make test` inside the provider folder

```sh
$ make test
```

## Rationale
Working in Operations engineering organizations infrastructure is often driven by tickets. Why not track infrastructure using tickets but this time we will use code.
This just showcases that you can pretty much Terraform anything!


## Credits

- Anubhav Mishra (anubhavmishra)
- Matthias Bartelmess (fourplusone)
