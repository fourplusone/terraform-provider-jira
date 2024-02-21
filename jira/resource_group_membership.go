package jira

import (
	"fmt"
	"net/url"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// GroupMembership The struct sent to the JIRA instance to create a new GroupMembership
type GroupMembership struct {
	AccountId string `json:"name,omitempty" structs:"name,omitempty"`
}

// Group
type Group struct {
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

// Groups List of groups the user belongs to
type Groups struct {
	Items []Group `json:"items,omitempty" structs:"items,omitempty"`
}

// UserGroups Wrapper for the groups of a user
type UserGroups struct {
	Groups Groups `json:"groups,omitempty" structs:"groups,omitempty"`
}

func getGroups(jiraClient *jira.Client, accountId string) (*UserGroups, *jira.Response, error) {

	relativeURL, _ := url.Parse("/rest/api/2/user")
	query := relativeURL.Query()
	query.Set("accountId", accountId)
	query.Set("expand", "groups")

	relativeURL.RawQuery = query.Encode()

	req, err := jiraClient.NewRequest("GET", relativeURL.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(UserGroups)
	resp, err := jiraClient.Do(req, user)
	if err != nil {
		return nil, resp, jira.NewJiraError(resp, err)
	}
	return user, resp, nil
}

// resourceGroupMembership is used to define a JIRA issue
func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Delete: resourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceGroupMembershipCreate creates a new jira issue using the jira api
func resourceGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	account_id := d.Get("account_id").(string)
	group := d.Get("group").(string)

	groupMembership := new(GroupMembership)
	groupMembership.AccountId = account_id

	relativeURL, _ := url.Parse(groupUserAPIEndpoint)
	query := relativeURL.Query()
	query.Set("groupname", group)
	relativeURL.RawQuery = query.Encode()

	err := request(config.jiraClient, "POST", relativeURL.String(), groupMembership, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	d.SetId(fmt.Sprintf("%s:%s", account_id, group))

	return resourceGroupMembershipRead(d, m)
}

// resourceGroupMembershipRead reads issue details using jira api
func resourceGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	components := strings.SplitN(d.Id(), ":", 2)
	accountId := components[0]
	groupname := components[1]

	groups, _, err := getGroups(config.jiraClient, accountId)
	if err != nil {
		return errors.Wrap(err, "getting jira group failed")
	}

	d.Set("account_id", accountId)
	d.Set("group", groupname)

	for _, group := range groups.Groups.Items {
		if group.Name == groupname {
			return nil
		}
	}

	return errors.Errorf("Cannot find group %s", groupname)
}

// resourceGroupMembershipDelete deletes jira issue using the jira api
func resourceGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	relativeURL, _ := url.Parse(groupUserAPIEndpoint)

	query := relativeURL.Query()
	query.Set("accountId", d.Get("account_id").(string))
	query.Set("groupname", d.Get("group").(string))

	relativeURL.RawQuery = query.Encode()

	err := request(config.jiraClient, "DELETE", relativeURL.String(), nil, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
