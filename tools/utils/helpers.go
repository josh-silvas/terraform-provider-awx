package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

// DiagFetch : Return the message for the fetch method
func DiagFetch(method string, id interface{}, err error) diag.Diagnostics {
	if err == nil {
		err = fmt.Errorf("fetch failed")
	}
	return Diagf(
		fmt.Sprintf("Unable to fetch %s", method),
		"Unable to fetch %s with id %v, got %s",
		method, id, err,
	)
}

// DiagCreate : Return the message for the create method
func DiagCreate(method string, err error) diag.Diagnostics {
	if err == nil {
		err = fmt.Errorf("create failed")
	}
	return Diagf(
		fmt.Sprintf("Unable to create %s", method),
		"Unable to create %s got %s",
		method, err,
	)
}

// DiagUpdate : Return the message for the update method
func DiagUpdate(method string, id interface{}, err error) diag.Diagnostics {
	if err == nil {
		err = fmt.Errorf("update failed")
	}
	return Diagf(
		fmt.Sprintf("Unable to update %s", method),
		"Unable to update %s with id %v: got %s",
		method, id, err,
	)
}

// DiagNotFound : Return the message for the not found method
func DiagNotFound(method string, id interface{}, err error) diag.Diagnostics {
	if err == nil {
		err = fmt.Errorf("not found")
	}
	return Diagf(
		fmt.Sprintf("Unable to fetch %s", method),
		"Unable to load %s with id %v: got %s",
		method, id, err,
	)
}

// DiagDelete : Return the message for the delete method
func DiagDelete(method string, id interface{}, err error) diag.Diagnostics {
	if err == nil {
		err = fmt.Errorf("delete failed")
	}
	return Diagf(
		fmt.Sprintf("Fail to delete %s", method),
		"Fail to delete %s with id %v: got %s",
		method, id, err,
	)
}

// StateIDToInt : Convert the ID from the state to an integer
func StateIDToInt(tfElement string, d *schema.ResourceData) (int, diag.Diagnostics) {
	var diags diag.Diagnostics
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return id, Diagf(
			fmt.Sprintf("%s, State ID Not Converted", tfElement),
			"Value in State %s is not numeric, %s", d.Id(), err,
		)
	}
	return id, diags
}

// Diagf : Return the message for the diag method
func Diagf(diagSummary, diagDetails string, detailsVars ...interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  diagSummary,
		Detail:   fmt.Sprintf(diagDetails, detailsVars...),
	})
	return diags
}

// Normalize : normalize the input data of type interface{} to string of either JSON or YAML
func Normalize(s interface{}) string {
	result := ""
	if j, ok := NormalizeJSONOk(s); ok {
		result = j
	} else if y, ok := NormalizeYamlOk(s); ok {
		result = y
	} else {
		result = s.(string)
	}
	return result
}

// NormalizeJSONOk : normalize the input data of type interface{} to string of JSON
func NormalizeJSONOk(s interface{}) (string, bool) {
	if s == nil || s == "" {
		return "", true
	}
	var j interface{}
	err := json.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %s", err), false
	}
	b, _ := json.Marshal(j)
	return string(b[:]), true
}

// NormalizeYamlOk : normalize the input data of type interface{} to string of YAML
func NormalizeYamlOk(s interface{}) (string, bool) {
	if s == nil || s == "" {
		return "", true
	}
	var j interface{}
	if err := yaml.Unmarshal([]byte(s.(string)), &j); err != nil {
		return fmt.Sprintf("Error parsing YAML: %s", err), false
	}
	b, err := yaml.Marshal(j)
	if err != nil {
		return fmt.Sprintf("Error parsing YAML: %s", err), false
	}
	return string(b[:]), true
}

// UnmarshalYAML : Unmarshal the YAML string to map[string]interface{}
func UnmarshalYAML(str string) map[string]interface{} {
	asMap := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(str), &asMap)
	if err != nil {
		asMap = nil
	}
	return asMap
}

// MarshalYAML : Marshal the interface{} to YAML string
func MarshalYAML(v interface{}) string {
	extraDataBytes, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(extraDataBytes)
}

// UnmarshalJSON : Unmarshal the JSON string to map[string]interface{}
func UnmarshalJSON(str string) map[string]interface{} {
	asMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &asMap)
	if err != nil {
		asMap = nil
	}
	return asMap
}

// MarshalJSON : Marshal the interface{} to JSON string
func MarshalJSON(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
