package awx

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialsRead,
		Description: "Data source for fetching credentials from AWX",
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier for the credential",
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username for the credential",
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kind of credential",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the credential",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the credential",
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceCredentialsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)

	creds, err := client.CredentialsService.ListCredentials(map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   "Unable to fetch credentials from AWX API",
		})
		return diags
	}

	parsedCreds := make([]map[string]interface{}, 0)
	for _, c := range creds {
		parsedCreds = append(parsedCreds, map[string]interface{}{
			"id":          c.ID,
			"username":    c.Inputs["username"],
			"kind":        c.Kind,
			"description": c.Description,
			"name":        c.Name,
		})
	}

	if err := d.Set("credentials", parsedCreds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
