package jira

import (
	"fmt"
	"strconv"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// ProjectRequest The struct sent to the JIRA instance to create a new Project
type ProjectRequest struct {
	Key                 string `json:"key,omitempty" structs:"key,omitempty"`
	Name                string `json:"name,omitempty" structs:"name,omitempty"`
	ProjectTypeKey      string `json:"projectTypeKey,omitempty" structs:"projectTypeKey,omitempty"`
	ProjectTemplateKey  string `json:"projectTemplateKey,omitempty" structs:"projectTemplateKey,omitempty"`
	Description         string `json:"description,omitempty" structs:"description,omitempty"`
	Lead                string `json:"lead,omitempty" structs:"lead,omitempty"`
	LeadAccountID       string `json:"leadAccountId,omitempty" structs:"leadAccountId,omitempty"`
	URL                 string `json:"url,omitempty" structs:"url,omitempty"`
	AssigneeType        string `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	AvatarID            int    `json:"avatar_id,omitempty" structs:"avatar_id,omitempty"`
	IssueSecurityScheme int    `json:"issueSecurityScheme,omitempty" structs:"issueSecurityScheme,omitempty"`
	PermissionScheme    int    `json:"permissionScheme,omitempty" structs:"permissionScheme,omitempty"`
	NotificationScheme  int    `json:"notificationScheme,omitempty" structs:"notificationScheme,omitempty"`
	CategoryID          int    `json:"categoryId,omitempty" structs:"categoryId,omitempty"`
}

type SharedConfigurationProjectResponse struct {
	ProjectID int `json:"projectId,omitempty"`
}

// IDResponse The struct sent from the JIRA instance after creating a new Project
type IDResponse struct {
	ID int `json:"id,omitempty" structs:"id,omitempty"`
}

// GetJiraResourceID Fetches the ID of a JIRA resource
func GetJiraResourceID(client *jira.Client, urlStr string) (*int, error) {
	req, err := client.NewRequest("GET", urlStr, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Creating Request failed")
	}

	response := new(IDResponse)

	resp, err := client.Do(req, response)
	if err != nil {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Creating Project Request failed")
	}

	return &response.ID, nil
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"shared_configuration_project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"project_type_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_template_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"lead": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"lead_account_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignee_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UNASSIGNED",
			},
			"avatar_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"issue_security_scheme": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"permission_scheme": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"notification_scheme": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"category_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// resourceProjectCreate creates a new jira issue using the jira api
func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	sharedProjectID, useSharedConfiguration := d.GetOk("shared_configuration_project_id")
	if useSharedConfiguration {
		project := &ProjectRequest{
			Key:           d.Get("key").(string),
			Name:          d.Get("name").(string),
			Lead:          d.Get("lead").(string),
			LeadAccountID: d.Get("lead_account_id").(string),
		}

		returnedProject := new(SharedConfigurationProjectResponse)

		endpoint := projectWithSharedConfigurationAPIEndpoint(sharedProjectID.(int))

		err := request(config.jiraClient, "POST", endpoint, project, returnedProject)

		if err != nil {
			return errors.Wrap(err, "Request failed")
		}

		d.SetId(strconv.Itoa(returnedProject.ProjectID))

		err = resourceProjectUpdate(d, m)

		if err != nil {
			return errors.Wrap(err, "Request failed")
		}

	} else {
		project := &ProjectRequest{
			Key:                 d.Get("key").(string),
			Name:                d.Get("name").(string),
			ProjectTypeKey:      d.Get("project_type_key").(string),
			ProjectTemplateKey:  d.Get("project_template_key").(string),
			Description:         d.Get("description").(string),
			Lead:                d.Get("lead").(string),
			LeadAccountID:       d.Get("lead_account_id").(string),
			URL:                 d.Get("url").(string),
			AssigneeType:        d.Get("assignee_type").(string),
			AvatarID:            d.Get("avatar_id").(int),
			IssueSecurityScheme: d.Get("issue_security_scheme").(int),
			PermissionScheme:    d.Get("permission_scheme").(int),
			NotificationScheme:  d.Get("notification_scheme").(int),
			CategoryID:          d.Get("category_id").(int),
		}

		returnedProject := new(IDResponse)

		err := request(config.jiraClient, "POST", projectAPIEndpoint, project, returnedProject)
		if err != nil {
			return errors.Wrap(err, "Request failed")
		}

		d.SetId(strconv.Itoa(returnedProject.ID))
	}

	return resourceProjectRead(d, m)

}

// resourceProjectRead reads issue details using jira api
func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	project, _, err := config.jiraClient.Project.Get(d.Id())
	if err != nil {
		return errors.Wrap(err, "getting jira project failed")
	}

	id, _ := strconv.Atoi(d.Id())
	d.Set("project_id", id)
	d.Set("key", project.Key)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("lead", project.Lead)
	d.Set("lead_account_id", project.Lead.AccountID)
	d.Set("url", project.URL)
	d.Set("assignee_type", project.AssigneeType)
	d.Set("category_id", project.ProjectCategory.ID)

	issuesecuritylevelscheme, err := GetJiraResourceID(config.jiraClient, fmt.Sprintf("%s/%s/issuesecuritylevelscheme", projectAPIEndpoint, d.Id()))
	if err != nil {
		return errors.Wrap(err, "getting issuesecuritylevelscheme failed")
	}
	d.Set("issue_security_scheme", issuesecuritylevelscheme)

	notificationscheme, err := GetJiraResourceID(config.jiraClient, fmt.Sprintf("%s/%s/notificationscheme", projectAPIEndpoint, d.Id()))
	if err != nil {
		return errors.Wrap(err, "getting notificationscheme failed")
	}
	d.Set("notification_scheme", notificationscheme)

	permissionscheme, err := GetJiraResourceID(config.jiraClient, fmt.Sprintf("%s/%s/permissionscheme", projectAPIEndpoint, d.Id()))
	if err != nil {
		return errors.Wrap(err, "getting permissionscheme failed")
	}
	d.Set("permission_scheme", permissionscheme)

	return nil
}

// resourceProjectUpdate updates jira issue using jira api
func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	project := &ProjectRequest{
		Key:                 d.Get("key").(string),
		Name:                d.Get("name").(string),
		ProjectTypeKey:      d.Get("project_type_key").(string),
		ProjectTemplateKey:  d.Get("project_template_key").(string),
		Description:         d.Get("description").(string),
		Lead:                d.Get("lead").(string),
		LeadAccountID:       d.Get("lead_account_id").(string),
		URL:                 d.Get("url").(string),
		AssigneeType:        d.Get("assignee_type").(string),
		AvatarID:            d.Get("avatar_id").(int),
		IssueSecurityScheme: d.Get("issue_security_scheme").(int),
		PermissionScheme:    d.Get("permission_scheme").(int),
		NotificationScheme:  d.Get("notification_scheme").(int),
		CategoryID:          d.Get("category_id").(int),
	}
	urlStr := fmt.Sprintf("%s/%s", projectAPIEndpoint, d.Id())

	returnedProject := new(jira.Project)

	err := request(config.jiraClient, "PUT", urlStr, project, returnedProject)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return resourceProjectRead(d, m)
}

// resourceProjectDelete deletes jira issue using the jira api
func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	urlStr := fmt.Sprintf("%s/%s", projectAPIEndpoint, d.Id())

	err := request(config.jiraClient, "DELETE", urlStr, nil, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	return nil
}
