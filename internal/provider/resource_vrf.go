package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceVRF() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a VRF.",

		CreateContext: resourceVRFCreate,
		ReadContext:   resourceVRFRead,
		UpdateContext: resourceVRFUpdate,
		DeleteContext: resourceVRFDelete,

		Schema: map[string]*schema.Schema{
			"address_family": {
				Description: "VRF address family.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_version": {
							Description: "Address family version.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"maximum_routes": {
							Description: "Maximum routes.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"maxiumum_routes_warning_only": {
							Description: "Maximum routes warning only.",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"description": {
				Description: "VRF description.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "VRF name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"rd": {
				Description: "VRF Route Distinguisher.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"route_target": {
				Description: "VRF Route-targets.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"community": {
							Description:  "Route-Target community.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"export", "import"}, false),
						},
						"rt": {
							Description: "Route-target.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceVRFCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.VRFDefinition{}
	params.VRFDefinition.Name = id

	getCreateUpdateVRFObject(d, &params.VRFDefinition)

	err := client.CreateVRF(params)

	if err != nil {
		return diag.Errorf("error creating Vlan. %s", err)
	}

	d.SetId(id)

	return resourceVRFRead(ctx, d, meta)
}

func resourceVRFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.VRFDefinition{}
	params.VRFDefinition.Name = id

	resp, err := client.ReadVRF(params)

	if err != nil {
		return diag.Errorf("error retrieving Vlan. %s", err)
	}

	resourceSetVRF(d, &resp.VRFDefinition)

	d.SetId(id)

	return nil
}

func resourceVRFUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.VRFDefinition{}
	params.VRFDefinition.Name = id

	getCreateUpdateVRFObject(d, &params.VRFDefinition)

	err := client.UpdateVRF(params)

	if err != nil {
		return diag.Errorf("error updating Vlan. %s", err)
	}

	d.SetId(id)

	return resourceVRFRead(ctx, d, meta)
}

func resourceVRFDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.VRFDefinition{}
	params.VRFDefinition.Name = id

	err := client.DeleteVRF(params)

	if err != nil {
		return diag.Errorf("error deleting Vlan. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetVRF(d *schema.ResourceData, resp *models.VRF) {
	if resp.Description != nil {
		d.Set("description", resp.Description)
	}
	if resp.RD != nil {
		d.Set("rd", resp.RD)
	}
	d.Set("name", resp.Name)
	d.Set("address_family", flattenVRFAddressFamily(resp.AddressFamily))
	d.Set("route_target", flattenVRFRouteTarget(resp.RouteTarget))
}

func flattenVRFAddressFamily(input *models.VRFAddressFamily) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	if l := input; l != nil {
		if input.Ipv4 != nil {
			output := map[string]interface{}{}
			output["ip_version"] = 4
			results = append(results, output)
		}
		if input.Ipv6 != nil {
			output := map[string]interface{}{}
			output["ip_version"] = 6
			results = append(results, output)
		}
	}
	return results
}

func flattenVRFRouteTarget(input *models.VRFRouteTargetList) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	if l := input; l != nil {
		for _, v := range input.Export {
			output := map[string]interface{}{}
			output["community"] = "export"
			output["rt"] = *v.AsnIP
			results = append(results, output)
		}
		for _, v := range input.Import {
			output := map[string]interface{}{}
			output["community"] = "import"
			output["rt"] = *v.AsnIP
			results = append(results, output)
		}
	}
	return results
}

func getCreateUpdateVRFObject(d *schema.ResourceData, m *models.VRF) *models.VRF {
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.Description = &s
		}
	}
	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			m.Name = s
		}
	}
	if v, ok := d.GetOk("rd"); ok {
		if s, ok := v.(string); ok {
			m.RD = &s
		}
	}
	if _, ok := d.GetOk("address_family"); ok {
		o := expandVRFAddressFamily(d, "address_family")
		m.AddressFamily = o
	}
	if _, ok := d.GetOk("route_target"); ok {
		o := expandVRFRouteTarget(d, "route_target")
		m.RouteTarget = o
	}
	return m
}

func expandVRFAddressFamily(d *schema.ResourceData, field string) *models.VRFAddressFamily {
	l := d.Get(field).([]interface{})
	if len(l) == 0 {
		return nil
	}

	r := models.VRFAddressFamily{}
	for _, v := range l {
		m := v.(map[string]interface{})
		if m["ip_version"].(int) == 4 {
			af4 := models.VRFAddressFamilyIpv4{}
			r.Ipv4 = &af4
		}
		if m["ip_version"].(int) == 6 {
			af6 := models.VRFAddressFamilyIpv6{}
			r.Ipv6 = &af6
		}
	}

	return &r
}

func expandVRFRouteTarget(d *schema.ResourceData, field string) *models.VRFRouteTargetList {
	l := d.Get(field).([]interface{})
	if len(l) == 0 {
		return nil
	}

	r := models.VRFRouteTargetList{}
	imp := []models.VRFRouteTarget{}
	exp := []models.VRFRouteTarget{}
	for _, v := range l {
		m := v.(map[string]interface{})
		if m["community"] == "import" {
			if v, ok := m["rt"].(string); ok {
				tmp := &models.VRFRouteTarget{
					AsnIP: &v,
				}
				imp = append(imp, *tmp)
			}
		}
		if m["community"] == "export" {
			if v, ok := m["rt"].(string); ok {
				tmp := &models.VRFRouteTarget{
					AsnIP: &v,
				}
				exp = append(exp, *tmp)
			}
		}
	}
	r.Import = imp
	r.Export = exp

	return &r
}
