package jira

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// PermissionGrant represents a permission grant within a JIRA Permission Scheme
type PermissionGrant struct {
	ID         int                   `json:"id,omitempty"`
	Holder     PermissionGrantHolder `json:"holder"`
	Permission string                `json:"permission"`
}

// PermissionGrantHolder represents who holds the permission
type PermissionGrantHolder struct {
	Type      string `json:"type"`
	Parameter string `json:"parameter,omitempty"`
}

func permissionSchemeGrantAPIEndpoint(schemeID string) string {
	return fmt.Sprintf("%s/%s/permission", permissionSchemeAPIEndpoint, schemeID)
}

func resourcePermissionSchemeGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionSchemeGrantCreate,
		Read:   resourcePermissionSchemeGrantRead,
		Delete: resourcePermissionSchemeGrantDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePermissionSchemeGrantImport,
		},

		Schema: map[string]*schema.Schema{
			"scheme_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the permission scheme",
			},
			"permission": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The permission to grant (e.g. BROWSE_PROJECTS, CREATE_ISSUES, ADMINISTER_PROJECTS)",
			},
			"holder_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the permission holder (e.g. group, user, projectRole, projectLead, assignee, reporter, applicationRole)",
			},
			"holder_parameter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The parameter for the holder type (e.g. group name, user key, role ID). Not required for types like projectLead, assignee, reporter.",
			},
		},
	}
}

func resourcePermissionSchemeGrantCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	schemeID := d.Get("scheme_id").(string)

	grant := &PermissionGrant{
		Permission: d.Get("permission").(string),
		Holder: PermissionGrantHolder{
			Type:      d.Get("holder_type").(string),
			Parameter: d.Get("holder_parameter").(string),
		},
	}

	returnedGrant := new(PermissionGrant)

	err := request(config.jiraClient, "POST", permissionSchemeGrantAPIEndpoint(schemeID), grant, returnedGrant)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(strconv.Itoa(returnedGrant.ID))

	return resourcePermissionSchemeGrantRead(d, m)
}

func resourcePermissionSchemeGrantRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	schemeID := d.Get("scheme_id").(string)

	urlStr := fmt.Sprintf("%s/%s", permissionSchemeGrantAPIEndpoint(schemeID), d.Id())

	grant := new(PermissionGrant)
	err := request(config.jiraClient, "GET", urlStr, nil, grant)

	if err != nil {
		if errors.Is(err, ResourceNotFoundError) {
			d.SetId("")
			return nil
		}
		return errors.Wrap(err, "Request failed")
	}

	d.Set("permission", grant.Permission)
	d.Set("holder_type", grant.Holder.Type)
	d.Set("holder_parameter", grant.Holder.Parameter)

	return nil
}

func resourcePermissionSchemeGrantDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	schemeID := d.Get("scheme_id").(string)

	urlStr := fmt.Sprintf("%s/%s", permissionSchemeGrantAPIEndpoint(schemeID), d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}

// resourcePermissionSchemeGrantImport supports importing with the format "scheme_id/grant_id"
func resourcePermissionSchemeGrantImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	schemeID, grantID, err := parseImportID(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid import ID %q, expected format: scheme_id/grant_id", d.Id())
	}

	d.Set("scheme_id", schemeID)
	d.SetId(grantID)

	return []*schema.ResourceData{d}, nil
}
