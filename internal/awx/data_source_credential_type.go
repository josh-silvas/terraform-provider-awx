package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceCredentialTypeByID() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialTypeByIDRead,
		Description: "Use this data source to get the details of a credential type by ID.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the credential type",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the credential type",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the credential type",
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The kind of the credential type",
			},
			"inputs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The inputs of the credential type",
			},
			"injectors": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The injectors of the credential type",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceCredentialTypeByIDRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id := d.Get("id").(int)
	credType, err := client.CredentialTypeService.GetCredentialTypeByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential type",
			Detail:   fmt.Sprintf("Unable to fetch credential type with ID: %d. Error: %s", id, err.Error()),
		})
	}

	if err := d.Set("name", credType.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", credType.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kind", credType.Kind); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inputs", credType.Inputs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("injectors", credType.Injectors); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))

	return diags
}
