---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jira_issue_link Resource - terraform-provider-jira"
subcategory: ""
description: |-
  
---

# jira_issue_link (Resource)



## Example Usage

```terraform
resource "jira_issue_link_type" "blocks" {
  name = "Blocks"
  inward = "is blocked by"
  outward = "blocks"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `inward_key` (String)
- `link_type` (String)
- `outward_key` (String)

### Read-Only

- `id` (String) The ID of this resource.


