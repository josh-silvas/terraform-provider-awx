package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagTeamRole = "Team Role"

func dataSourceTeamRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRolesRead,
		Description: "Use this data source to get the details of a team role in AWX.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the team role.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the team role.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"team_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the team.",
			},
		},
	}
}

func dataSourceTeamRolesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	teamID := d.Get("team_id").(int)

	team, err := client.TeamService.GetTeamByID(teamID, params)
	if err != nil {
		return utils.DiagFetch(diagTeamRole, teamID, err)
	}

	rolesList := []*awx.ApplyRole{
		team.SummaryFields.ObjectRoles.AdminRole,
		team.SummaryFields.ObjectRoles.MemberRole,
		team.SummaryFields.ObjectRoles.ReadRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range rolesList {
			if v != nil && id == v.ID {
				d = setTeamRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range rolesList {
			if v != nil && name == v.Name {
				d = setTeamRoleData(d, v)
				return diags
			}
		}
	}

	return utils.DiagNotFound(diagTeamRole, teamID, nil)
}

func setTeamRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
