# terraform-provider-jira
Terraform Provider for creating, updating and deleting JIRA issues.
This can be used to interlink infrastructure management with JIRA issues closely.

## Install

* Download `terraform-provider-jira` binary from [Github](https://github.com/anubhavmishra/terraform-provider-jira/releases)
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
export JIRA_USERNAME=username
export JIRA_PASSWORD=password
```

Create terraform config file

```hcl
resource "jira_issue" "example" {
  assignee    = "anubhavmishra"
  reporter    = "anubhavmishra"

  issue_type  = "Task"

  // description is optional  
  description = "This is a test issue"
  summary     = "Created using Terraform"

  project_key = "PROJ"
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


