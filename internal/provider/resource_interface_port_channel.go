package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourcePortChannel() *schema.Resource {
	return &schema.Resource{
		Description: "Manage an PortChannel interface.",

		CreateContext: resourcePortChannelCreate,
		ReadContext:   resourcePortChannelRead,
		UpdateContext: resourcePortChannelUpdate,
		DeleteContext: resourcePortChannelDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Description: "Interface description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ip": {
				Description:  "Primary interface IP as CIDR.",
				Type:         schema.TypeString,
				Optional:     true,
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
			"vrf": {
				Description: "VRF.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourcePortChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannel{}
	params.PortChannel.Name = id

	getCreateUpdatePortChannelObject(d, &params)

	err := client.CreatePortChannel(params)

	if err != nil {
		return diag.Errorf("error creating PortChannel. %s", err)
	}

	d.SetId(id)

	return resourcePortChannelRead(ctx, d, meta)
}

func resourcePortChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannel{}
	params.PortChannel.Name = id
	resp, err := client.ReadPortChannel(params)

	if err != nil {
		return diag.Errorf("error retrieving PortChannel. %s", err)
	}

	resourceSetPortChannel(d, &resp.PortChannel)

	d.SetId(id)

	return nil
}

func resourcePortChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannel{}
	params.PortChannel.Name = id
	getCreateUpdatePortChannelObject(d, &params)

	err := client.UpdatePortChannel(params)

	if err != nil {
		return diag.Errorf("error updating PortChannel. %s", err)
	}

	d.SetId(id)

	return resourcePortChannelRead(ctx, d, meta)
}

func resourcePortChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("name").(string)

	params := models.PortChannel{}
	params.PortChannel.Name = id
	err := client.DeletePortChannel(params)

	if err != nil {
		return diag.Errorf("error deleting PortChannel. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetPortChannel(d *schema.ResourceData, resp *models.Interface) {
	if resp.Description != nil {
		d.Set("description", resp.Description)
	}
	if resp.IP != nil {
		if resp.IP.Address.Primary != nil {
			resp.IP.Address.Primary.SetCIDR()
			d.Set("ip", resp.IP.Address.Primary.CIDR)
		}
		d.Set("secondary_ip", flattenPortChannelSecondaryIPs(resp.IP.Address.Secondary))
	}
	d.Set("name", resp.Name)
	if resp.Shutdown != nil {
		d.Set("shutdown", true)
	} else {
		d.Set("shutdown", false)
	}
}

func flattenPortChannelSecondaryIPs(input *[]models.SecondaryIPAddress) []map[string]interface{} {
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

func getCreateUpdatePortChannelObject(d *schema.ResourceData, m *models.PortChannel) *models.PortChannel {
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.PortChannel.Description = &s
		}
	}
	m.PortChannel.IP = &models.InterfaceIP{
		Address: &models.Address{},
	}
	if v, ok := d.GetOk("ip"); ok {
		if s, ok := v.(string); ok {
			ip := &models.IPAddress{
				CIDR: s,
			}
			m.PortChannel.IP.Address.Primary = ip
			m.PortChannel.IP.Address.Primary.SetNetmask()
		}
	}
	if _, ok := d.GetOk("secondary_ip"); ok {
		o := expandPortChannelSecondaryIPs(d, "secondary_ip")
		m.PortChannel.IP.Address.Secondary = o
	}
	if v, ok := d.GetOk("shutdown"); ok {
		if b, ok := v.(bool); ok {
			if b {
				n := models.CiscoEnabled
				m.PortChannel.Shutdown = &n
			}
		}
	}
	if v, ok := d.GetOk("vrf"); ok {
		if s, ok := v.(string); ok {
			o := &models.InterfaceVrf{
				Forwarding: &s,
			}
			m.PortChannel.Vrf = o
		}
	}
	return m
}

func expandPortChannelSecondaryIPs(d *schema.ResourceData, field string) *[]models.SecondaryIPAddress {
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
