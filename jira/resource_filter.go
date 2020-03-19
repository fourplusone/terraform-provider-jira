package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	jira "github.com/andygrunwald/go-jira"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// FilterRequest represents a Filter in Jira
type FilterRequest struct {
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	Jql              string                   `json:"jql"`
	Favourite        bool                     `json:"favourite"`
	SharePermissions []FilterPermissionResult `json:"sharePermissions"`
}

type FilterPermissionRequest struct {
	Type          string `json:"type"`
	ProjectID     string `json:"projectId"`
	Group         string `json:"groupname"`
	ProjectRoleID string `json:"projectRoleId"`
}

type ProjectPermission struct {
	ID string `json:"id"`
}

type GroupPermission struct {
	Name string `json:"name"`
}

type RolePermission struct {
	ID int `json:"id"`
}

type FilterPermissionResult struct {
	Type          string            `json:"type"`
	ID            int               `json:"id"`
	Project       ProjectPermission `json:"project"`
	Group         GroupPermission   `json:"group"`
	ProjectRoleID RolePermission    `json:"role"`
}

// resourceFilter is used to define a JIRA Filter
func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceFilterCreate,
		Read:   resourceFilterRead,
		Update: resourceFilterUpdate,
		Delete: resourceFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"jql": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"favourite": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceFilterPermissionsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, s string) ([]string, []error) {
								if !(v.(string) == "global" ||
									v.(string) == "group" ||
									v.(string) == "project" ||
									v.(string) == "project_role" ||
									v.(string) == "authenticated") {
									return nil, []error{fmt.Errorf("type needs to be one of global, group, project, project_role or authenticated")}
								}
								return nil, nil
							},
						},
						"project_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"project_role_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func setFilter(w *FilterRequest, d *schema.ResourceData) {
	w.Name = d.Get("name").(string)
	w.Description = d.Get("description").(string)
	w.Jql = d.Get("jql").(string)
	w.Favourite = d.Get("favourite").(bool)
}

func setFilterResource(w *jira.Filter, d *schema.ResourceData) {
	d.SetId(w.ID)
	d.Set("name", w.Name)
	d.Set("description", w.Description)
	d.Set("jql", w.Jql)
	d.Set("favourite", w.Favourite)

	permissions := &schema.Set{
		F: resourceFilterPermissionsHash,
	}

	for _, f := range w.SharePermissions {
		permissionResultMap := f.(map[string]interface{})
		permissionResult := FilterPermissionResult{}
		marshalledJSON, _ := json.Marshal(permissionResultMap)

		json.Unmarshal(marshalledJSON, &permissionResult)

		projectRoleID := strconv.Itoa(permissionResult.ProjectRoleID.ID)
		if projectRoleID == "0" {
			projectRoleID = ""
		}

		permissionID := strconv.Itoa(permissionResult.ID)
		if permissionID == "0" {
			permissionID = ""
		}

		permissionType := permissionResult.Type
		if permissionType == "loggedin" {
			permissionType = "authenticated"
		}

		m := map[string]interface{}{
			"group_name":      permissionResult.Group.Name,
			"project_id":      permissionResult.Project.ID,
			"project_role_id": projectRoleID,
			"type":            permissionType,
			"id":              permissionID,
		}
		permissions.Add(m)
	}

	d.Set("permissions", permissions)

}

// resourceFilterCreate creates a new jira filter using the jira api
func resourceFilterCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	filter := new(FilterRequest)
	permissions := d.Get("permissions").(*schema.Set)
	returnedFilter := new(jira.Filter)
	setFilter(filter, d)

	err := request(config.jiraClient, "POST", filterAPIEndpoint, filter, returnedFilter)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setFilterResource(returnedFilter, d)

	err = filterAddPermissions(permissions.List(), returnedFilter.ID, config)
	if err != nil {
		return err
	}

	setFilterResource(returnedFilter, d)

	return nil
}

// resourceFilterRead reads filter details using jira api
func resourceFilterRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", filterAPIEndpoint, d.Id())

	filter := new(jira.Filter)
	err := request(config.jiraClient, "GET", urlStr, nil, filter)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	setFilterResource(filter, d)

	return nil
}

// resourceFilterUpdate updates jira filter using jira api
func resourceFilterUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	if d.HasChange("permissions") {
		o, n := d.GetChange("permissions")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		err := filterRevokePermissions(os.Difference(ns).List(), d.Id(), config)
		if err != nil {
			return err
		}

		err = filterAddPermissions(ns.Difference(os).List(), d.Id(), config)
		if err != nil {
			return err
		}
	}

	// First issueing a PUT seems to cause a racing condition in JIRA.
	// Therefore the following order must be obeyed:
	// * Removing old permissions
	// * Set new permissions
	// * Update the filter itself

	filter := new(FilterRequest)
	setFilter(filter, d)

	urlStr := fmt.Sprintf("%s/%s", filterAPIEndpoint, d.Id())
	returnedFilter := new(jira.Filter)

	err := request(config.jiraClient, "PUT", urlStr, filter, returnedFilter)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceFilterRead(d, m)
}

// resourceFilterDelete deletes jira filter using the jira api
func resourceFilterDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", filterAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}

func resourceFilterPermissionsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["type"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["project_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["project_role_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["group_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func filterRevokePermissions(configured []interface{}, filterID string, config *Config) error {
	for _, data := range configured {
		d := data.(map[string]interface{})
		url := fmt.Sprintf("%s/%s", filterPermissionEndpoint(filterID), d["id"].(string))

		err := request(
			config.jiraClient,
			"DELETE",
			url,
			nil, nil)

		if err != nil {
			return errors.Wrap(err, "Request failed")
		}
	}
	return nil
}

func filterAddPermissions(configured []interface{}, filterID string, config *Config) error {

	for _, data := range configured {
		d := data.(map[string]interface{})
		permission := FilterPermissionRequest{
			Type:          d["type"].(string),
			ProjectID:     d["project_id"].(string),
			Group:         d["group_name"].(string),
			ProjectRoleID: d["project_role_id"].(string),
		}
		err := request(
			config.jiraClient,
			"POST",
			filterPermissionEndpoint(filterID),
			permission, nil)

		if err != nil {
			return errors.Wrap(err, "Request failed")
		}
	}
	return nil
}
