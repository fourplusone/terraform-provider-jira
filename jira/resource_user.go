package jira

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func usernameFallbackSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return old == d.Get("name")
	}
	return old == new
}

// resourceUser is used to define a JIRA issue
func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: usernameFallbackSuppressFunc,
			},
		},
	}
}

// resourceUserCreate creates a new jira user using the jira api
func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	user := new(jira.User)
	user.Name = d.Get("name").(string)
	user.EmailAddress = d.Get("email").(string)

	dn, ok := d.GetOkExists("display_name")
	user.DisplayName = dn.(string)

	if !ok {
		user.DisplayName = user.Name
	}

	createdUser, _, err := config.jiraClient.User.Create(user)

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(createdUser.Key)

	return resourceUserRead(d, m)
}

// resourceUserRead reads issue details using jira api
func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	user, _, err := config.jiraClient.User.Get(d.Id())
	if err != nil {
		return errors.Wrap(err, "getting jira user failed")
	}

	d.Set("name", user.Name)
	d.Set("display_name", user.DisplayName)
	d.Set("email", user.EmailAddress)
	return nil
}

// resourceUserDelete deletes jira issue using the jira api
func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	_, err := config.jiraClient.User.Delete(d.Id())

	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
