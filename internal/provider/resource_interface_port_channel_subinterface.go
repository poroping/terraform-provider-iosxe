package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourcePortChannelSubinterface() *schema.Resource {
	return &schema.Resource{
		Description: "Manage an PortChannelSubinterface interface.",

		CreateContext: resourcePortChannelSubinterfaceCreate,
		ReadContext:   resourcePortChannelSubinterfaceRead,
		UpdateContext: resourcePortChannelSubinterfaceUpdate,
		DeleteContext: resourcePortChannelSubinterfaceDelete,

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
				Required:    true,
				ForceNew:    true,
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

func resourcePortChannelSubinterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannelSubinterface{}
	params.PortChannelSubinterface.Name = id

	getCreateUpdatePortChannelSubinterfaceObject(d, &params)

	err := client.CreatePortChannelSubinterface(params)

	if err != nil {
		return diag.Errorf("error creating PortChannelSubinterface. %s", err)
	}

	d.SetId(id)

	return resourcePortChannelSubinterfaceRead(ctx, d, meta)
}

func resourcePortChannelSubinterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannelSubinterface{}
	params.PortChannelSubinterface.Name = id
	resp, err := client.ReadPortChannelSubinterface(params)

	if err != nil {
		return diag.Errorf("error retrieving PortChannelSubinterface. %s", err)
	}

	resourceSetPortChannelSubinterface(d, &resp.PortChannelSubinterface)

	d.SetId(id)

	return nil
}

func resourcePortChannelSubinterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannelSubinterface{}
	params.PortChannelSubinterface.Name = id
	getCreateUpdatePortChannelSubinterfaceObject(d, &params)

	err := client.UpdatePortChannelSubinterface(params)

	if err != nil {
		return diag.Errorf("error updating PortChannelSubinterface. %s", err)
	}

	d.SetId(id)

	return resourcePortChannelSubinterfaceRead(ctx, d, meta)
}

func resourcePortChannelSubinterfaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannelSubinterface{}
	params.PortChannelSubinterface.Name = id
	err := client.DeletePortChannelSubinterface(params)

	if err != nil {
		return diag.Errorf("error deleting PortChannelSubinterface. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetPortChannelSubinterface(d *schema.ResourceData, resp *models.Interface) {
	if resp.Description != nil {
		d.Set("description", resp.Description)
	}
	resp.IP.Address.Primary.SetCIDR()
	d.Set("ip", resp.IP.Address.Primary.CIDR)
	d.Set("name", resp.Name)
	d.Set("secondary_ip", flattenPortChannelSubinterfaceSecondaryIPs(resp.IP.Address.Secondary))
	if resp.Shutdown != nil {
		d.Set("shutdown", true)
	} else {
		d.Set("shutdown", false)
	}
	d.Set("vlanid", resp.Encapsulation.Dot1Q.VlanID)
	if resp.Vrf != nil {
		d.Set("vrf", resp.Vrf.Forwarding)
	}
}

func flattenPortChannelSubinterfaceSecondaryIPs(input *[]models.SecondaryIPAddress) []map[string]interface{} {
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

func getCreateUpdatePortChannelSubinterfaceObject(d *schema.ResourceData, m *models.PortChannelSubinterface) *models.PortChannelSubinterface {
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.PortChannelSubinterface.Description = &s
		}
	}
	m.PortChannelSubinterface.IP = &models.InterfaceIP{
		Address: &models.Address{},
	}
	if v, ok := d.GetOk("ip"); ok {
		if s, ok := v.(string); ok {
			ip := &models.IPAddress{
				CIDR: s,
			}
			m.PortChannelSubinterface.IP.Address.Primary = ip
			m.PortChannelSubinterface.IP.Address.Primary.SetNetmask()
		}
	}
	if _, ok := d.GetOk("secondary_ip"); ok {
		o := expandPortChannelSubinterfaceSecondaryIPs(d, "secondary_ip")
		m.PortChannelSubinterface.IP.Address.Secondary = o
	}
	if v, ok := d.GetOk("shutdown"); ok {
		if b, ok := v.(bool); ok {
			if b {
				n := models.CiscoEnabled
				m.PortChannelSubinterface.Shutdown = &n
			}
		}
	}
	if v, ok := d.GetOk("vrf"); ok {
		if s, ok := v.(string); ok {
			o := &models.InterfaceVrf{
				Forwarding: &s,
			}
			m.PortChannelSubinterface.Vrf = o
		}
	}
	if v, ok := d.GetOk("vlanid"); ok {
		if i, ok := v.(int); ok {
			i2 := int64(i)
			m.PortChannelSubinterface.Encapsulation = &models.InterfaceEncapsulation{
				Dot1Q: &models.InterfaceEncapsulationDot1Q{
					VlanID: &i2,
				},
			}
		}
	}
	return m
}

func expandPortChannelSubinterfaceSecondaryIPs(d *schema.ResourceData, field string) *[]models.SecondaryIPAddress {
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
