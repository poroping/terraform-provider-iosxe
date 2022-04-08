package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage an Vlan interface.",

		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,

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
				Description:  "VLANID.",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4096),
			},
			"vrf": {
				Description: "VRF.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.Vlan.Name = strconv.Itoa(id)

	getCreateUpdateVlanObject(d, &params)

	err := client.CreateVlan(params)

	if err != nil {
		return diag.Errorf("error creating Vlan. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceVlanRead(ctx, d, meta)
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.Vlan.Name = strconv.Itoa(id)
	resp, err := client.ReadVlan(params)

	if err != nil {
		return diag.Errorf("error retrieving Vlan. %s", err)
	}

	resourceSetVlan(d, &resp.Vlan)

	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.Vlan.Name = strconv.Itoa(id)
	getCreateUpdateVlanObject(d, &params)

	err := client.UpdateVlan(params)

	if err != nil {
		return diag.Errorf("error updating Vlan. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceVlanRead(ctx, d, meta)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.Vlan.Name = strconv.Itoa(id)
	err := client.DeleteVlan(params)

	if err != nil {
		return diag.Errorf("error deleting Vlan. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetVlan(d *schema.ResourceData, resp *models.Interface) {
	if resp.Description != nil {
		d.Set("description", resp.Description)
	}
	resp.IP.Address.Primary.SetCIDR()
	d.Set("ip", resp.IP.Address.Primary.CIDR)
	d.Set("name", fmt.Sprintf("%s%v", models.VlanName, resp.Name))
	d.Set("secondary_ip", flattenVlanSecondaryIPs(resp.IP.Address.Secondary))
	if resp.Shutdown != nil {
		d.Set("shutdown", true)
	} else {
		d.Set("shutdown", false)
	}
	id, err := strconv.Atoi(resp.Name)
	if err != nil {
		log.Printf("[WARN] Vlan ID not castable to int")
	}
	d.Set("vlanid", id)
	if resp.Vrf != nil {
		d.Set("vrf", resp.Vrf.Forwarding)
	}
}

func flattenVlanSecondaryIPs(input *[]models.SecondaryIPAddress) []map[string]interface{} {
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

func getCreateUpdateVlanObject(d *schema.ResourceData, m *models.Vlan) *models.Vlan {
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.Vlan.Description = &s
		}
	}
	m.Vlan.IP = &models.InterfaceIP{
		Address: &models.Address{},
	}
	if v, ok := d.GetOk("ip"); ok {
		if s, ok := v.(string); ok {
			ip := &models.IPAddress{
				CIDR: s,
			}
			m.Vlan.IP.Address.Primary = ip
			m.Vlan.IP.Address.Primary.SetNetmask()
		}
	}
	if _, ok := d.GetOk("secondary_ip"); ok {
		o := expandVlanSecondaryIPs(d, "secondary_ip")
		m.Vlan.IP.Address.Secondary = o
	}
	if v, ok := d.GetOk("shutdown"); ok {
		if b, ok := v.(bool); ok {
			if b {
				n := models.CiscoEnabled
				m.Vlan.Shutdown = &n
			}
		}
	}
	if v, ok := d.GetOk("vrf"); ok {
		if s, ok := v.(string); ok {
			o := &models.InterfaceVrf{
				Forwarding: &s,
			}
			m.Vlan.Vrf = o
		}
	}
	return m
}

func expandVlanSecondaryIPs(d *schema.ResourceData, field string) *[]models.SecondaryIPAddress {
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
		temp.Secondary = &models.CiscoEnabled
		r = append(r, temp)
	}

	return &r
}
