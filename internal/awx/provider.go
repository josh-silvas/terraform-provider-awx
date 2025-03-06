package awx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

// Provider returns a schema.Provider for AWX.
func Provider() *schema.Provider { //nolint:funlen
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_HOSTNAME", "http://localhost"),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable SSL verification of API calls",
			},
			"ca_pem": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Path to a CA Certificate in PEM format to be used to verify the server",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_USERNAME", "admin"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_PASSWORD", "password"),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_TOKEN", ""),
			},
			"http_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Optional. HTTP headers mapping keys to values used for accessing the AWX Api.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awx_credential_azure_key_vault":                          resourceCredentialAzureKeyVault(),
			"awx_credential_google_compute_engine":                    resourceCredentialGoogleComputeEngine(),
			"awx_credential_container_registry":                       resourceCredentialContainerRegistry(),
			"awx_credential_input_source":                             resourceCredentialInputSource(),
			"awx_credential":                                          resourceCredential(),
			"awx_credential_type":                                     resourceCredentialType(),
			"awx_credential_machine":                                  resourceCredentialMachine(),
			"awx_credential_scm":                                      resourceCredentialSCM(),
			"awx_credential_gitlab":                                   resourceCredentialGitlab(),
			"awx_credential_galaxy":                                   resourceCredentialGalaxy(),
			"awx_credential_vault":                                    resourceCredentialVault(),
			"awx_execution_environment":                               resourceExecutionEnvironment(),
			"awx_host":                                                resourceHost(),
			"awx_instance_group":                                      resourceInstanceGroup(),
			"awx_inventory_group":                                     resourceInventoryGroup(),
			"awx_inventory_source":                                    resourceInventorySource(),
			"awx_inventory":                                           resourceInventory(),
			"awx_inventory_instance_groups":                           resourceInventoryInstanceGroups(),
			"awx_job_template_credential":                             resourceJobTemplateCredentials(),
			"awx_job_template_instance_groups":                        resourceJobTemplateInstanceGroups(),
			"awx_job_template":                                        resourceJobTemplate(),
			"awx_job_template_launch":                                 resourceJobTemplateLaunch(),
			"awx_job_template_notification_template_error":            resourceJobTemplateNotificationTemplateError(),
			"awx_job_template_notification_template_started":          resourceJobTemplateNotificationTemplateStarted(),
			"awx_job_template_notification_template_success":          resourceJobTemplateNotificationTemplateSuccess(),
			"awx_notification_template":                               resourceNotificationTemplate(),
			"awx_organization":                                        resourceOrganization(),
			"awx_organization_galaxy_credential":                      resourceOrganizationsGalaxyCredentials(),
			"awx_organization_instance_groups":                        resourceOrganizationsInstanceGroups(),
			"awx_project":                                             resourceProject(),
			"awx_schedule":                                            resourceSchedule(),
			"awx_settings_ldap_team_map":                              resourceSettingsLDAPTeamMap(),
			"awx_setting":                                             resourceSetting(),
			"awx_team":                                                resourceTeam(),
			"awx_user":                                                resourceUser(),
			"awx_workflow_job_template_node_always":                   resourceWorkflowJobTemplateNodeAlways(),
			"awx_workflow_job_template_node_failure":                  resourceWorkflowJobTemplateNodeFailure(),
			"awx_workflow_job_template_node_success":                  resourceWorkflowJobTemplateNodeSuccess(),
			"awx_workflow_job_template_node":                          resourceWorkflowJobTemplateNode(),
			"awx_workflow_job_template":                               resourceWorkflowJobTemplate(),
			"awx_workflow_job_template_schedule":                      resourceWorkflowJobTemplateSchedule(),
			"awx_workflow_job_template_notification_template_error":   resourceWorkflowJobTemplateNotificationTemplateError(),
			"awx_workflow_job_template_notification_template_started": resourceWorkflowJobTemplateNotificationTemplateStarted(),
			"awx_workflow_job_template_notification_template_success": resourceWorkflowJobTemplateNotificationTemplateSuccess(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"awx_credential_azure_key_vault": dataSourceCredentialAzure(),
			"awx_credential":                 dataSourceCredentialByID(),
			"awx_credential_role":            dataSourceCredentialMachineRole(),
			"awx_credential_type":            dataSourceCredentialTypeByID(),
			"awx_credentials":                dataSourceCredentials(),
			"awx_execution_environment":      dataSourceExecutionEnvironment(),
			"awx_inventory_group":            dataSourceInventoryGroup(),
			"awx_inventory":                  dataSourceInventory(),
			"awx_inventory_role":             dataSourceInventoryRole(),
			"awx_job_template":               dataSourceJobTemplate(),
			"awx_job_template_role":          dataSourceJobTemplateRole(),
			"awx_notification_template":      dataSourceNotificationTemplate(),
			"awx_organization":               dataSourceOrganization(),
			"awx_organization_role":          dataSourceOrganizationRole(),
			"awx_organizations":              dataSourceOrganizations(),
			"awx_project":                    dataSourceProject(),
			"awx_project_role":               dataSourceProjectRole(),
			"awx_schedule":                   dataSourceSchedule(),
			"awx_workflow_job_template":      dataSourceWorkflowJobTemplate(),
			"awx_workflow_job_template_role": dataSourceWorkflowJobTemplateRole(),
			"awx_team":                       dataSourceTeam(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type HeadersRoundTripper struct {
	r       http.RoundTripper
	headers map[string]string
}

func (hrt HeadersRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	for header, value := range hrt.headers {
		r.Header.Set(header, value)
	}
	return hrt.r.RoundTrip(r)
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostname := d.Get("hostname").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	token := d.Get("token").(string)
	caPem := d.Get("ca_pem").(string)

	headers := map[string]string{}
	if httpHeaders, ok := d.GetOk("http_headers"); ok {
		for k, v := range httpHeaders.(map[string]interface{}) {
			headers[k] = v.(string)
		}
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if d.Get("insecure").(bool) {
		//nolint:gosec
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if caPem != "" {
		certPool := x509.NewCertPool()
		if caCertPem, err := os.ReadFile(caPem); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to read file",
				Detail:   fmt.Sprintf("Unable to read certificate file located at %s.", caPem),
			})
			return nil, diags
		} else if ok := certPool.AppendCertsFromPEM(caCertPem); !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to parse certificate.",
				Detail:   "Unable to parse certificate. Check that the certificate is in a valid PEM format.",
			})
			return nil, diags
		}
		customTransport.TLSClientConfig = &tls.Config{RootCAs: certPool, MinVersion: tls.VersionTLS12}
	}

	client := &http.Client{
		Transport: HeadersRoundTripper{r: customTransport, headers: headers},
	}

	var c *awx.AWX
	var err error
	if token != "" {
		c, err = awx.NewAWXToken(hostname, token, client)
	} else {
		c, err = awx.NewAWX(hostname, username, password, client)
	}
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create AWX client",
			Detail:   "Unable to auth user against AWX API: check the hostname, username and password",
		})
		return nil, diags
	}

	return c, diags
}
