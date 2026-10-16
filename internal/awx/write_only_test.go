package awx

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// writeOnlyResources lists every resource that exposes write-only "<field>_wo" secrets.
var writeOnlyResources = []struct {
	name     string
	resource func() *schema.Resource
	secrets  map[string]string
}{
	{"awx_credential_machine", resourceCredentialMachine, credentialMachineSecrets},
	{"awx_credential_scm", resourceCredentialSCM, credentialSCMSecrets},
	{"awx_credential_azure_key_vault", resourceCredentialAzureKeyVault, credentialAzureKeyVaultSecrets},
	{"awx_credential_container_registry", resourceCredentialContainerRegistry, credentialContainerRegistrySecrets},
	{"awx_credential_galaxy", resourceCredentialGalaxy, credentialGalaxySecrets},
	{"awx_credential_gitlab", resourceCredentialGitlab, credentialGitlabSecrets},
	{"awx_credential_google_compute_engine", resourceCredentialGoogleComputeEngine, credentialGoogleComputeEngineSecrets},
	{"awx_credential_vault", resourceCredentialVault, credentialVaultSecrets},
	{"awx_user", resourceUser, userSecrets},
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

func TestWriteOnlyRequiredSecretNeedsExactlyOne(t *testing.T) {
	base := map[string]interface{}{
		"name":            "test",
		"organization_id": 1,
		"url":             "https://vault.example.com",
		"client":          "client",
		"tenant":          "tenant",
	}
	cases := []struct {
		name    string
		extra   map[string]interface{}
		wantErr bool
	}{
		{"legacy", map[string]interface{}{"secret": "s"}, false},
		{"write-only", map[string]interface{}{"secret_wo": "s"}, false},
		{"neither", nil, true},
		{"both", map[string]interface{}{"secret": "s", "secret_wo": "s"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{}
			for k, v := range base {
				raw[k] = v
			}
			for k, v := range tc.extra {
				raw[k] = v
			}

			diags := resourceCredentialAzureKeyVault().Validate(terraform.NewResourceConfigRaw(raw))
			if diags.HasError() != tc.wantErr {
				t.Errorf("Validate error = %v, want error %v", diags, tc.wantErr)
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

func TestCredentialSchemaIsValid(t *testing.T) {
	if err := resourceCredential().InternalValidate(nil, true); err != nil {
		t.Fatalf("invalid schema: %v", err)
	}
	if wo := resourceCredential().Schema["secret_inputs_wo"]; !wo.WriteOnly || !wo.Sensitive {
		t.Error("secret_inputs_wo: want a sensitive write-only attribute")
	}
}

func TestMergeSecretInputs(t *testing.T) {
	got, err := mergeSecretInputs(map[string]interface{}{"username": "admin"}, `{"password": "secret", "verify": true}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["username"] != "admin" || got["password"] != "secret" || got["verify"] != true {
		t.Errorf("merged inputs = %v", got)
	}

	if _, err := mergeSecretInputs(map[string]interface{}{"password": "x"}, `{"password": "y"}`); err == nil {
		t.Error("want an error for a key set in both inputs and secret_inputs_wo")
	}
	if _, err := mergeSecretInputs(map[string]interface{}{}, `["password"]`); err == nil {
		t.Error("want an error for a non-object JSON value")
	}
}

func TestStateCredentialInputs(t *testing.T) {
	remote := map[string]interface{}{
		"username":     "admin",
		"password":     "$encrypted$",
		"ssh_key_data": "$encrypted$",
	}
	prior := map[string]interface{}{
		"username":     "admin",
		"ssh_key_data": "legacy",
	}

	got := stateCredentialInputs(prior, remote)

	if _, ok := got["password"]; ok {
		t.Error("password: encrypted input not tracked in state should be left out")
	}
	if got["ssh_key_data"] != "$encrypted$" {
		t.Errorf("ssh_key_data = %v, want $encrypted$ (legacy input tracked in state)", got["ssh_key_data"])
	}
	if got["username"] != "admin" {
		t.Errorf("username = %v, want admin", got["username"])
	}
}
