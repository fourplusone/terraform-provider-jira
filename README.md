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

This should initialize the plugin.



