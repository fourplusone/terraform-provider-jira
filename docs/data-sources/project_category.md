# jira_project_category Data Source

Get a project category from Jira

## Example Usage

```hcl
data "jira_project_category" "test" {
  project_id = "10041"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The Id of the project category.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Id of the project category.
* `description` - A block of text describing the category.
* `name` - The name of the project category.
