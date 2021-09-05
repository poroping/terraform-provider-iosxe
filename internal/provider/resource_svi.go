package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceSVI() *schema.Resource {
	return &schema.Resource{
		Description: "Manage an SVI interface.",

		CreateContext: resourceSVICreate,
		ReadContext:   resourceSVIRead,
		UpdateContext: resourceSVIUpdate,
		DeleteContext: resourceSVIDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Description: "Interface description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ip": {
				Description:  "Primary interface IP as CIDR.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"name": {
				Description: "Interface name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"secondary_ip": {
				Description: "Secondary IPs.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Description:  "Secondary interface IP as CIDR.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
					},
				},
			},
			"shutdown": {
				Description: "Interface status.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Default:     nil,
			},
			"vlanid": {
				Description: "VLANID.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceSVICreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.SVI{}
	params.Vlan.Name = id

	getCreateUpdateSVIObject(d, &params)

	err := client.CreateSVI(params)

	if err != nil {
		return diag.Errorf("error creating SVI. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceSVIRead(ctx, d, meta)
}

func resourceSVIRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.SVI{}
	params.Vlan.Name = id

	resp, err := client.ReadSVI(params)

	if err != nil {
		return diag.Errorf("error retrieving SVI. %s", err)
	}

	resourceSetSVI(d, &resp.Vlan)

	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceSVIUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.SVI{}
	params.Vlan.Name = id

	getCreateUpdateSVIObject(d, &params)

	err := client.UpdateSVI(params)

	if err != nil {
		return diag.Errorf("error updating SVI. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceSVIRead(ctx, d, meta)
}

func resourceSVIDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.SVI{}
	params.Vlan.Name = id

	err := client.DeleteSVI(params)

	if err != nil {
		return diag.Errorf("error deleting SVI. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetSVI(d *schema.ResourceData, resp *models.L3Interface) {
	if resp.Description != nil {
		d.Set("description", resp.Description)
	}
	resp.IP.Address.Primary.SetCIDR()
	d.Set("ip", resp.IP.Address.Primary.CIDR)
	name := int(resp.Name.(float64))
	d.Set("name", fmt.Sprintf("%s%d", models.SVIName, name))
	d.Set("secondary_ip", flattenSVISecondaryIPs(resp.IP.Address.Secondary))
	if resp.Shutdown != nil {
		d.Set("shutdown", true)
	} else {
		d.Set("shutdown", false)
	}
	d.Set("vlanid", resp.Name)
}

func flattenSVISecondaryIPs(input *[]models.SecondaryIPAddress) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	if l := input; l != nil {
		for _, v := range *input {
			output := map[string]interface{}{}
			v.SetCIDR()
			output["ip"] = v.CIDR
			results = append(results, output)
		}
	}
	return results
}

func getCreateUpdateSVIObject(d *schema.ResourceData, m *models.SVI) *models.SVI {
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.Vlan.Description = &s
		}
	}
	if v, ok := d.GetOk("ip"); ok {
		if s, ok := v.(string); ok {
			m.Vlan.IP.Address.Primary.CIDR = s
			m.Vlan.IP.Address.Primary.SetNetmask()
		}
	}
	if _, ok := d.GetOk("secondary_ip"); ok {
		o := expandSVISecondaryIPs(d, "secondary_ip")
		m.Vlan.IP.Address.Secondary = o
	}
	if v, ok := d.GetOk("shutdown"); ok {
		if b, ok := v.(bool); ok {
			if b {
				n := json.RawMessage(models.ExplicitNull)
				m.Vlan.Shutdown = &n
			}
		}
	}
	return m
}

func expandSVISecondaryIPs(d *schema.ResourceData, field string) *[]models.SecondaryIPAddress {
	l := d.Get(field).([]interface{})
	if len(l) == 0 {
		return nil
	}

	r := make([]models.SecondaryIPAddress, 0, len(l))
	for _, v := range l {
		temp := models.SecondaryIPAddress{}
		m := v.(map[string]interface{})
		temp.CIDR = m["ip"].(string)
		temp.SetNetmask()
		empty := json.RawMessage(models.ExplicitNull)
		temp.Secondary = &empty
		r = append(r, temp)
	}

	return &r
}
