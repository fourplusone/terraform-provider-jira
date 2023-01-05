---
layout: "jira"
page_title: "Provider: Jira"
description: |-
  The Jira provider provides resources to interact with and manage Jira.
---

# Jira Provider

The Jira provider provides resources to interact with and manage Jira projects.

## Example Usage

Terraform 0.13 and later:

```terraform
terraform {
  required_providers {
    jira = {
      source = "fourplusone/jira"
      version = "0.1.15"
    }
  }
}

# Configure the Jira Provider
provider "jira" {
  url      = "https://myjira.atlassian.net" # Can also be set using the JIRA_URL environment variable
  user     = "xxxx"                         # Can also be set using the JIRA_USER environment variable
  password = "xxxx"                         # Can also be set using the JIRA_PASSWORD environment variable
  token = "xxxx"                            # Can also be set using the JIRA_TOKEN environment variable
}
```

## Schema

- **url** (String, Required) URL for your Jira instance. Can be specified with the JIRA_URL environment variable.

- **user** (String, Optional) Username for your user. Can be specified with the JIRA_USER environment variable.

- **password** (String, Optional) Password for the user, can also be an API Token. Can be specified with the JIRA_PASSWORD environment variable.

- **token** (String, Optional) Password for the user, can also be an API Token. Can be specified with the JIRA_PASSWORD environment variable.
