package awx

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCredentialMachineSchemaIsValid(t *testing.T) {
	if err := resourceCredentialMachine().InternalValidate(nil, true); err != nil {
		t.Fatalf("invalid schema: %v", err)
	}
}

func TestCredentialMachineWriteOnlySecrets(t *testing.T) {
	s := resourceCredentialMachine().Schema

	for field := range credentialMachineSecrets {
		wo := s[field+writeOnlySuffix]
		if wo == nil || !wo.WriteOnly || !wo.Sensitive {
			t.Errorf("%s%s: want a sensitive write-only attribute", field, writeOnlySuffix)
		}
		if v := s[field+writeOnlyVersionSuffix]; v == nil || v.WriteOnly {
			t.Errorf("%s%s: want a regular (stored) attribute", field, writeOnlyVersionSuffix)
		}
	}
}

func TestGetSecretFallsBackToLegacyAttribute(t *testing.T) {
	r := resourceCredentialMachine()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":            "test",
		"organization_id": 1,
		"password":        "legacy",
	})

	if got := getSecret(d, "password"); got != "legacy" {
		t.Errorf("getSecret = %q, want legacy", got)
	}
}

func TestStringFromCty(t *testing.T) {
	cases := []struct {
		name   string
		in     cty.Value
		want   string
		wantOK bool
	}{
		{"set", cty.StringVal("secret"), "secret", true},
		{"empty", cty.StringVal(""), "", true},
		{"null", cty.NullVal(cty.String), "", false},
		{"unknown", cty.UnknownVal(cty.String), "", false},
		{"dynamic", cty.DynamicVal, "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stringFromCty(tc.in)
			if got != tc.want || ok != tc.wantOK {
				t.Errorf("stringFromCty = (%q, %v), want (%q, %v)", got, ok, tc.want, tc.wantOK)
			}
		})
	}
}
