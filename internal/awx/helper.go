package awx

import (
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
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

// writeOnlySuffix and writeOnlyVersionSuffix name the write-only variant of a secret field and the
// version field used to trigger an update when the write-only value changes.
const (
	writeOnlySuffix        = "_wo"
	writeOnlyVersionSuffix = "_wo_version"
)

// addWriteOnlySecret adds "<field>_wo" and "<field>_wo_version" to the schema. The write-only value is never
// persisted to state (requires Terraform >= 1.11), so bumping the version is what makes Terraform send it again.
// A required field becomes optional, with exactly one of "<field>" and "<field>_wo" required instead.
func addWriteOnlySecret(s map[string]*schema.Schema, field, description string) {
	wo := field + writeOnlySuffix
	version := field + writeOnlyVersionSuffix

	s[wo] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		WriteOnly:   true,
		Sensitive:   true,
		Description: description + " Write-only: never stored in state. Requires Terraform >= 1.11. Use `" + field + "` for older versions.",
	}
	if s[field].Required {
		s[field].Required = false
		s[field].Optional = true
		s[field].ExactlyOneOf = []string{field, wo}
		s[wo].ExactlyOneOf = []string{field, wo}
	} else {
		s[field].ConflictsWith = append(s[field].ConflictsWith, wo)
		s[wo].ConflictsWith = []string{field}
	}
	s[version] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		RequiredWith: []string{wo},
		Description:  "Version of `" + wo + "`. Change it to send a new value of `" + wo + "` to AWX.",
	}
}

// writeOnlyVersionKeys returns the "<field>_wo_version" keys for the given secrets, for use with d.HasChanges.
func writeOnlyVersionKeys(secrets map[string]string) []string {
	keys := make([]string, 0, len(secrets))
	for field := range secrets {
		keys = append(keys, field+writeOnlyVersionSuffix)
	}

	return keys
}

// getSecret returns the value of the write-only "<field>_wo" from the configuration when set, falling back to the
// legacy "<field>" attribute otherwise.
func getSecret(d *schema.ResourceData, field string) string {
	raw, diags := d.GetRawConfigAt(cty.GetAttrPath(field + writeOnlySuffix))
	if !diags.HasError() {
		if v, ok := stringFromCty(raw); ok {
			return v
		}
	}

	return d.Get(field).(string)
}

func stringFromCty(v cty.Value) (string, bool) {
	if v.IsNull() || !v.IsKnown() || !v.Type().Equals(cty.String) {
		return "", false
	}

	return v.AsString(), true
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
