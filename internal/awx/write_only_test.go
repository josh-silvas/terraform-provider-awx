package awx

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// writeOnlyResources lists every resource that exposes write-only "<field>_wo" secrets.
var writeOnlyResources = []struct {
	name     string
	resource func() *schema.Resource
	secrets  map[string]string
}{
	{"awx_credential_machine", resourceCredentialMachine, credentialMachineSecrets},
	{"awx_credential_scm", resourceCredentialSCM, credentialSCMSecrets},
}

func TestWriteOnlyResourceSchemasAreValid(t *testing.T) {
	for _, tc := range writeOnlyResources {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.resource().InternalValidate(nil, true); err != nil {
				t.Fatalf("invalid schema: %v", err)
			}
		})
	}
}

func TestWriteOnlyResourceSecrets(t *testing.T) {
	for _, tc := range writeOnlyResources {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.resource().Schema
			for field := range tc.secrets {
				if s[field] == nil || !s[field].Sensitive {
					t.Errorf("%s: want an existing sensitive attribute", field)
				}
				wo := s[field+writeOnlySuffix]
				if wo == nil || !wo.WriteOnly || !wo.Sensitive {
					t.Errorf("%s%s: want a sensitive write-only attribute", field, writeOnlySuffix)
				}
				if v := s[field+writeOnlyVersionSuffix]; v == nil || v.WriteOnly {
					t.Errorf("%s%s: want a regular (stored) attribute", field, writeOnlyVersionSuffix)
				}
			}
		})
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
