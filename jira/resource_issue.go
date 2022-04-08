package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/trivago/tgo/tcontainer"
)

// resourceIssue is used to define a JIRA issue
func resourceIssue() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssueCreate,
		Read:   resourceIssueRead,
		Update: resourceIssueUpdate,
		Delete: resourceIssueDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIssueImport,
		},

		Schema: map[string]*schema.Schema{
			"assignee": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: caseInsensitiveSuppressFunc,
			},
			"reporter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return caseInsensitiveSuppressFunc(k, old, new, d)
				},
			},
			"fields": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
			"issue_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"labels": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"summary": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parent": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return old == new
				},
			},
			"state_transitions": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
			"delete_transitions": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
			// Computed values
			"issue_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceIssueCreate creates a new jira issue using the jira api
func resourceIssueCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	assignee := d.Get("assignee")
	reporter := d.Get("reporter")
	parent := d.Get("parent")
	fields := d.Get("fields")
	issueType := d.Get("issue_type").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels")
	summary := d.Get("summary").(string)
	projectKey := d.Get("project_key").(string)

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: description,
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: projectKey,
			},
			Summary: summary,
		},
	}

	if assignee != "" {
		i.Fields.Assignee = &jira.User{
			Name: assignee.(string),
		}
	}

	if reporter != "" {
		i.Fields.Reporter = &jira.User{
			Name: reporter.(string),
		}
	}

	if parent != "" {
		i.Fields.Parent = &jira.Parent{
			ID: parent.(string),
		}
	}

	if fields != nil {
		if err := setCustomFields(&i, &fields); err != nil {
			return err
		}
	}

	if labels != nil {
		for _, label := range labels.([]interface{}) {
			i.Fields.Labels = append(i.Fields.Labels, fmt.Sprintf("%v", label))
		}
	}

	issue, res, err := config.jiraClient.Issue.Create(&i)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "creating jira issue failed: %s", body)
	}

	issue, res, err = config.jiraClient.Issue.Get(issue.ID, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	if err := doStateTransition(d, issue, config); err != nil {
		return err
	}

	d.SetId(issue.ID)

	return resourceIssueRead(d, m)
}

// resourceIssueRead reads issue details using jira api
func resourceIssueRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	issue, res, err := config.jiraClient.Issue.Get(d.Id(), nil)
	if err != nil {
		if res.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	if issue.Fields.Assignee != nil {
		d.Set("assignee", issue.Fields.Assignee.Name)
	}

	if issue.Fields.Reporter != nil {
		d.Set("reporter", issue.Fields.Reporter.Name)
	}

	if issue.Fields.Parent != nil {
		d.Set("parent", issue.Fields.Parent.Key)
	}

	// Custom or non-standard fields
	var resourceFieldsRaw, resourceHasFields = d.GetOk("fields")
	if resourceHasFields {
		customFields, fieldsError := getCustomFields(&resourceFieldsRaw, issue)
		if fieldsError != nil {
			return fieldsError
		}
		d.Set("fields", *customFields)
	}

	d.Set("labels", nil)
	if issue.Fields.Labels != nil && len(issue.Fields.Labels) > 0 {
		d.Set("labels", issue.Fields.Labels)
	}

	d.Set("issue_type", issue.Fields.Type.Name)
	if issue.Fields.Description != "" {
		d.Set("description", issue.Fields.Description)
	}
	d.Set("summary", issue.Fields.Summary)
	d.Set("project_key", issue.Fields.Project.Key)
	d.Set("issue_key", issue.Key)
	d.Set("state", issue.Fields.Status.ID)

	return nil
}

// resourceIssueUpdate updates jira issue using jira api
func resourceIssueUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	issueKey := d.Get("issue_key").(string)

	i := jira.Issue{
		Key:    issueKey,
		ID:     d.Id(),
		Fields: &jira.IssueFields{},
	}

	if issueType := d.Get("issue_type").(string); d.HasChange("issue_type") {
		i.Fields.Type = jira.IssueType{
			Name: issueType,
		}
	}

	if description := d.Get("description").(string); d.HasChange("description") {
		i.Fields.Description = description
	}

	if summary := d.Get("summary").(string); d.HasChange("summary") {
		i.Fields.Summary = summary
	}

	if projectKey := d.Get("project_key").(string); d.HasChange("project_key") {
		i.Fields.Project = jira.Project{
			Key: projectKey,
		}
	}

	if assignee := d.Get("assignee"); d.HasChange("assignee") && assignee != "" {
		i.Fields.Assignee = &jira.User{
			Name: assignee.(string),
		}
	}

	if reporter := d.Get("reporter"); d.HasChange("reporter") && reporter != "" {
		i.Fields.Reporter = &jira.User{
			Name: reporter.(string),
		}
	}

	if labels := d.Get("labels"); d.HasChange("labels") && labels != nil {
		for _, label := range labels.([]interface{}) {
			i.Fields.Labels = append(i.Fields.Labels, fmt.Sprintf("%v", label))
		}
	}

	if fields := d.Get("fields"); d.HasChange("fields") && fields != nil && len(fields.(map[string]interface{})) > 0 {
		if i.Fields.Unknowns == nil {
			i.Fields.Unknowns = tcontainer.NewMarshalMap()
		}
		for field, value := range fields.(map[string]interface{}) {
			var decodedValue interface{}
			valueBytes := []byte(value.(string))

			if json.Valid(valueBytes) {
				if err := json.Unmarshal([]byte(value.(string)), &decodedValue); err != nil {
					return err
				}
				i.Fields.Unknowns.Set(field, decodedValue)
			} else {
				i.Fields.Unknowns.Set(field, value.(string))
			}
		}
	}

	issue, res, err := config.jiraClient.Issue.Update(&i)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "updating jira issue failed: %s", body)
	}

	issue, res, err = config.jiraClient.Issue.Get(issue.ID, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	if err := doStateTransition(d, issue, config); err != nil {
		return err
	}

	d.SetId(issue.ID)

	return resourceIssueRead(d, m)
}

// resourceIssueDelete deletes jira issue using the jira api
func resourceIssueDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	id := d.Id()

	if transition, ok := d.GetOk("delete_transitions"); ok {

		issue, res, err := config.jiraClient.Issue.Get(id, nil)
		if err != nil {
			if res.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			body, _ := ioutil.ReadAll(res.Body)
			return errors.Wrapf(err, "getting jira issue failed: %s", body)
		}
		currentState := issue.Fields.Status.ID

		if err := doStateTransitionOnDelete(id, &transition, config, currentState); err != nil {
			return err
		}

	} else {
		res, err := config.jiraClient.Issue.Delete(id)

		if err != nil {
			body, _ := ioutil.ReadAll(res.Body)
			return errors.Wrapf(err, "deleting jira issue failed: %s", body)
		}
	}

	return nil
}

// resourceIssueImport imports jira issue using the jira api
func resourceIssueImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	err := resourceIssueRead(d, m)
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	return []*schema.ResourceData{d}, nil
}

// extractSameKeys pulls the values from extendedInput which match keys is baseInput
func extractSameKeys(baseInput interface{}, extendedInput interface{}) interface{} {
	switch baseInput.(type) {
	case map[string]interface{}:
		if extendedInputMap, ok := extendedInput.(map[string]interface{}); ok {
			output := map[string]interface{}{}
			for key, _ := range baseInput.(map[string]interface{}) {
				output[key] = extendedInputMap[key]
			}
			return output
		}
	case []interface{}:
		if extendedInputSlice, ok := extendedInput.([]interface{}); ok {
			output := []interface{}{}
			for i, element1 := range baseInput.([]interface{}) {
				if len(extendedInputSlice) > i {
					element2 := extendedInputSlice[i]
					output = append(output, extractSameKeys(element1, element2))
				}
			}
			return output
		}
	}

	return extendedInput
}

func setCustomFields(i *jira.Issue, fields *interface{}) error {
	if i.Fields.Unknowns == nil {
		i.Fields.Unknowns = tcontainer.NewMarshalMap()
	}
	for field, value := range (*fields).(map[string]interface{}) {
		var decodedValue interface{}
		valueBytes := []byte(value.(string))

		if json.Valid(valueBytes) {
			if err := json.Unmarshal([]byte(value.(string)), &decodedValue); err != nil {
				return err
			}
			i.Fields.Unknowns.Set(field, decodedValue)
		} else {
			i.Fields.Unknowns.Set(field, value.(string))
		}
	}
	return nil
}

func doStateTransition(d *schema.ResourceData, issue *jira.Issue, config *Config) error {
	key := "state_transitions"
	if desiredState, ok := d.GetOk("state"); ok {
		currentState := issue.Fields.Status.ID
		if currentState != desiredState.(string) {
			if transitionMap, ok := d.GetOk(key); ok {
				transitionFound := false
				for state, transitionValue := range transitionMap.(map[string]interface{}) {
					if state != currentState {
						continue
					}
					transitionFound = true

					var transitions []string
					valueBytes := []byte(transitionValue.(string))
					if json.Valid(valueBytes) {
						if err := json.Unmarshal([]byte(transitionValue.(string)), &transitions); err != nil {
							return err
						}
						for _, transition := range transitions {
							res, err := config.jiraClient.Issue.DoTransition(issue.ID, transition)
							if err != nil {
								body, _ := ioutil.ReadAll(res.Body)
								return errors.Wrapf(err, "transitioning jira issue failed: %s", body)
							}
						}
					} else {
						return errors.Errorf("%s: cannot unmarshal provided string: %s", key, valueBytes)
					}
					break
				}
				if !transitionFound {
					return errors.Errorf("%s: cannot find transition sequence for state: %s", key, currentState)
				}
			}
		}
	}
	return nil
}

func doStateTransitionOnDelete(id string, transitionMap *interface{}, config *Config, currentState string) error {
	transitionFound := false
	for state, transitionValue := range (*transitionMap).(map[string]interface{}) {
		if state != currentState {
			continue
		}
		transitionFound = true

		var transitions []string
		valueBytes := []byte(transitionValue.(string))
		if json.Valid(valueBytes) {
			if err := json.Unmarshal([]byte(transitionValue.(string)), &transitions); err != nil {
				return err
			}
			for _, transition := range transitions {
				res, err := config.jiraClient.Issue.DoTransition(id, transition)
				if err != nil {
					body, _ := ioutil.ReadAll(res.Body)
					return errors.Wrapf(err, "transitioning on delete for jira issue failed: %s", body)
				}
			}
		} else {
			return errors.Errorf("delete_transitions: cannot unmarshal provided string: %s", valueBytes)
		}
		break
	}
	if !transitionFound {
		return errors.Errorf("delete_transitions: cannot find transition sequence for state: %s", currentState)
	}
	return nil
}

func getCustomFields(resourceFieldsRaw *interface{}, issue *jira.Issue) (*map[string]string, error) {
	incomingFields := make(map[string]string)
	resourceFields := (*resourceFieldsRaw).(map[string]interface{})
	for field := range issue.Fields.Unknowns {
		if existingField, fieldExists := resourceFields[field]; fieldExists {
			if value, valueExists := issue.Fields.Unknowns.Value(field); valueExists {
				existingFieldBytes := []byte(existingField.(string))

				if json.Valid(existingFieldBytes) {
					var decodedExistingValue interface{}
					if err := json.Unmarshal([]byte(existingField.(string)), &decodedExistingValue); err != nil {
						return nil, err
					}

					marshalledValue, _ := json.Marshal(extractSameKeys(decodedExistingValue, value))
					incomingFields[field] = string(marshalledValue)
				} else {
					switch value.(type) {
					case string:
						incomingFields[field] = value.(string)
					case bool:
						incomingFields[field] = fmt.Sprintf("%t", value.(bool))
					case int:
						incomingFields[field] = fmt.Sprintf("%d", value.(int))
					case float32:
						incomingFields[field] = fmt.Sprintf("%f", value.(float32))
					case float64:
						incomingFields[field] = fmt.Sprintf("%f", value.(float64))
					case uint:
						incomingFields[field] = fmt.Sprintf("%d", value.(uint))
					}
				}
			}
		}
	}
	return &incomingFields, nil
}
