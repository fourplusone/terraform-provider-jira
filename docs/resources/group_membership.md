---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jira_group_membership Resource - terraform-provider-jira"
subcategory: ""
description: |-
  
---

# jira_group_membership (Resource)



## Example Usage

```terraform
// Create a group named "Terraform Managed"
resource "jira_group" "tf_group" {
  name = "Terraform Managed"
}

// User "bot" will be a Member of "Terraform Managed"
resource "jira_group_membership" "gm_1" {
  username = "bot"
  group = "${jira_group.tf_group.name}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group` (String)
- `username` (String)

### Read-Only

- `id` (String) The ID of this resource.


