package jira

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

// API Endpoints
const filterAPIEndpoint = "/rest/api/2/filter"
const groupAPIEndpoint = "/rest/api/2/group"
const groupUserAPIEndpoint = "/rest/api/2/group/user"

const issueLinkAPIEndpoint = "/rest/api/2/issueLink"
const issueLinkTypeAPIEndpoint = "/rest/api/2/issueLinkType"
const issueTypeAPIEndpoint = "/rest/api/2/issuetype"

const projectAPIEndpoint = "/rest/api/2/project"
const projectCategoryAPIEndpoint = "/rest/api/2/projectCategory"
const roleAPIEndpoint = "/rest/api/2/role"
const userAPIEndpoint = "/rest/api/2/user"
const webhookAPIEndpoint = "/rest/webhooks/1.0/webhook"

func projectWithSharedConfigurationAPIEndpoint(projectID int) string {
	return fmt.Sprintf("/rest/project-templates/1.0/createshared/%d", projectID)
}

func projectRoleAPIEndpoint(projectKey string) string {
	return fmt.Sprintf("/rest/api/2/project/%s/role", projectKey)
}

func filterPermissionEndpoint(filterID string) string {
	return fmt.Sprintf("%s/%s/permission", filterAPIEndpoint, filterID)

}

func request(client *jira.Client, method string, endpoint string, in interface{}, out interface{}) error {

	req, err := client.NewRequest(method, endpoint, in)

	if err != nil {
		return errors.Wrapf(err, "Creating %s Request failed", method)
	}

	res, err := client.Do(req, out)
	if err != nil {

		if in != nil && res != nil {
			typeName := reflect.TypeOf(in).Name()
			body, readErr := ioutil.ReadAll(res.Response.Body)
			if readErr != nil {
				return errors.Wrapf(readErr, "Creating %s Request failed", typeName)
			}
			return errors.Wrapf(err, "Creating %s Request failed: %s", typeName, body)

		}
		return errors.Wrapf(err, "Creating Request failed")
	}

	return nil
}

func caseInsensitiveSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
