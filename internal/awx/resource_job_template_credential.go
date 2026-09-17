package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceJobTemplateCredentials() *schema.Resource {
	return &schema.Resource{
		Description:   "Associates a credential to a job template",
		CreateContext: resourceJobTemplateCredentialsCreate,
		DeleteContext: resourceJobTemplateCredentialsDelete,
		ReadContext:   resourceJobTemplateCredentialsRead,

		Schema: map[string]*schema.Schema{

			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the job template to associate the credential with",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the credential to associate with the job template",
			},
		},
	}
}

func resourceJobTemplateCredentialsCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)
	if _, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string)); err != nil {
		return utils.DiagNotFound("JobTemplate Credential", jobTemplateID, err)
	}

	result, err := client.JobTemplateService.AssociateCredentials(jobTemplateID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{})

	if err != nil {
		return utils.DiagCreate("JobTemplate AssociateCredentials", err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return nil
}

func resourceJobTemplateCredentialsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)
	credentialID := d.Get("credential_id").(int)

	credentials, err := client.JobTemplateService.ListJobTemplateCredentials(jobTemplateID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("JobTemplate Credential", jobTemplateID, err)
	}

	for _, credential := range credentials {
		if credential.ID == credentialID {
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceJobTemplateCredentialsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobTemplateID := d.Get("job_template_id").(int)
	res, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("JobTemplate Credential", jobTemplateID, err)
	}

	if _, err = client.JobTemplateService.DisAssociateCredentials(res.ID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{}); err != nil {
		return utils.DiagDelete("JobTemplate DisAssociateCredentials", jobTemplateID, err)
	}

	d.SetId("")
	return nil
}
