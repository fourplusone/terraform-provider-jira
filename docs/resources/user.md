---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jira_user Resource - terraform-provider-jira"
subcategory: ""
description: |-
  
---

# jira_user (Resource)



## Example Usage

```terraform
resource "jira_user" "demo_user" {
  name = "bot"
  email = "bot@example.org"
  display_name = "The Bot"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String)
- `name` (String)

### Optional

- `display_name` (String)

### Read-Only

- `id` (String) The ID of this resource.


