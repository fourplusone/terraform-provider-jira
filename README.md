# terraform-provider-jira

[![Build Status](https://travis-ci.com/fourplusone/terraform-provider-jira.svg?branch=master)](https://travis-ci.com/fourplusone/terraform-provider-jira)

Terraform Provider for managing JIRA. Supports

- Issues
- Issue Links
- Issue Types
- Issue Link Types

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
  summary     = "Created using Terraform"

  // description is optional  
  description = "This is a test issue"
  

  project_key = "PROJ"
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
