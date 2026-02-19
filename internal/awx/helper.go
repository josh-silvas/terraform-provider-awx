package awx

import (
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

// The AWX API returns '$encrypted$' in place of the password/ssh_key_data. We do not want to write that placeholder to the
// Terraform state file as it would break diffing and cause the SCM credential to be recreated on every apply.
func setSanitizedEncryptedValue(d *schema.ResourceData, fieldName string, value interface{}) error {
	if value == "$encrypted$" {
		stateValue := d.Get(fieldName).(string)
		if stateValue == "$encrypted$" {
			stateValue = "UPDATE_ME"
		}
		if err := d.Set(fieldName, stateValue); err != nil {
			return err
		}
	} else {
		if err := d.Set(fieldName, value); err != nil {
			return err
		}
	}

	return nil
}

func setSanitizedEncryptedCredential(d *schema.ResourceData, fieldName string, cred *awx.Credential) error {
	return setSanitizedEncryptedValue(d, fieldName, cred.Inputs[fieldName])
}

// suppressValueDiff suppresses diffs for setting values when AWX returns encrypted placeholders
func suppressValueDiff(k, oldVal, newVal string, d *schema.ResourceData) bool {
	// If the oldVal value contains $encrypted$, AWX has masked the secret
	// Suppress the diff since we can't compare the actual values
	if strings.Contains(oldVal, "$encrypted$") {
		return true
	}

	// For non-secret values, use the standard JSON comparison
	return SuppressEquivalentJSONDiffs(k, oldVal, newVal, d)
}

// suppressYAMLDiff suppresses diffs when YAML strings are semantically equivalent
func suppressYAMLDiff(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	// If both are empty, they're equivalent
	if oldVal == "" && newVal == "" {
		return true
	}

	// Parse both YAML strings and compare the resulting structures
	oldData := utils.UnmarshalYAML(oldVal)
	newData := utils.UnmarshalYAML(newVal)

	// If both failed to parse (both nil), compare as strings
	if oldData == nil && newData == nil {
		return oldVal == newVal
	}

	// If one parsed and one didn't, they're different
	if (oldData == nil) != (newData == nil) {
		return false
	}

	// Compare the parsed structures using reflection
	return reflect.DeepEqual(oldData, newData)
}
