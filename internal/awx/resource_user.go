package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

//nolint:funlen
func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource allows you to create, update, and delete a user in AWX.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		UpdateContext: resourceUserUpdate,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the user",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the user",
			},
			"first_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first name of the user",
			},
			"last_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The last name of the user",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email of the user",
			},
			"is_superuser": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The superuser status of the user",
			},
			"is_system_auditor": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The system auditor status of the user",
			},
			"role_entitlement": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of role IDs of the role entitlements",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeInt,
							Required: true,
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

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	userName := d.Get("username").(string)

	result, err := client.UserService.CreateUser(map[string]interface{}{
		"username":          userName,
		"password":          d.Get("password").(string),
		"first_name":        d.Get("first_name").(string),
		"last_name":         d.Get("last_name").(string),
		"email":             d.Get("email").(string),
		"is_superuser":      d.Get("is_superuser").(bool),
		"is_system_auditor": d.Get("is_system_auditor").(bool),
	}, map[string]string{})
	if err != nil {
		return utils.DiagCreate("Username", err)
	}

	d.SetId(strconv.Itoa(result.ID))

	if rent, entOk := d.GetOk("role_entitlement"); entOk {
		entset := rent.(*schema.Set).List()
		if err := roleUserEntitlementUpdate(m, result.ID, entset, false); err != nil {
			return utils.DiagCreate("Role entitlement", err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func roleUserEntitlementUpdate(m interface{}, userID int, roles []interface{}, remove bool) error {
	client := m.(*awx.AWX)
	awxService := client.UserService

	for _, v := range roles {
		emap := v.(map[string]interface{})
		payload := map[string]interface{}{
			"id": emap["role_id"],
		}
		if remove {
			payload["disassociate"] = true // presence of key triggers removal
		}

		_, err := awxService.UpdateUserRoleEntitlement(userID, payload, make(map[string]string))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("role_entitlement") {
		oi, ni := d.GetChange("role_entitlement")
		if oi == nil {
			oi = new(schema.Set)
		}
		if ni == nil {
			ni = new(schema.Set)
		}
		oe := oi.(*schema.Set)
		ne := ni.(*schema.Set)

		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		err := roleUserEntitlementUpdate(m, id, remove, true)
		if err != nil {
			return utils.DiagUpdate("User Role Entitlement", id, err)
		}
		err = roleUserEntitlementUpdate(m, id, add, false)
		if err != nil {
			return utils.DiagUpdate("User Role Entitlement", id, err)
		}
	}
	if _, err := client.UserService.UpdateUser(id, map[string]interface{}{
		"username":          d.Get("username").(string),
		"password":          d.Get("password").(string),
		"first_name":        d.Get("first_name").(string),
		"last_name":         d.Get("last_name").(string),
		"email":             d.Get("email").(string),
		"is_superuser":      d.Get("is_superuser").(bool),
		"is_system_auditor": d.Get("is_system_auditor").(bool),
	}, nil); err != nil {
		return utils.DiagUpdate("User", id, err)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := client.UserService.GetUserByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("User", id, err)
	}
	entitlements, _, err := client.UserService.ListUserRoleEntitlements(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("User Roles", id, err)
	}

	if err := d.Set("username", res.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedValue(d, "password", res.Password); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("first_name", res.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", res.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", res.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_superuser", res.IsSuperUser); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system_auditor", res.IsSystemAuditor); err != nil {
		return diag.FromErr(err)
	}

	var entlist []interface{}
	for _, v := range entitlements {
		elem := make(map[string]interface{})
		elem["role_id"] = v.ID
		entlist = append(entlist, elem)
	}
	f := schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_id": {Type: schema.TypeInt},
		}})

	ent := schema.NewSet(f, entlist)

	if err := d.Set("role_entitlement", ent); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.UserService
	id, diags := utils.StateIDToInt("Delete User", d)

	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteUser(id); err != nil {
		return utils.DiagDelete("User", id, err)
	}
	d.SetId("")
	return diags
}
