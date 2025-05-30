package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

// The AWX API returns '$encrypted$' in place of the password/ssh_key_data. We do not want to write that placeholder to the
// Terraform state file as it would break diffing and cause the SCM credential to be recreated on every apply.

func setSanitizedEncryptedValue(d *schema.ResourceData, fieldName string, value string) error {
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
	return setSanitizedEncryptedValue(d, fieldName, cred.Inputs[fieldName].(string))
}
