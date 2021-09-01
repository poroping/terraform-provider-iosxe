package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func dataSourceL3Vlan() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample data source in the Terraform provider L3Vlan.",

		ReadContext: dataSourceL3VlanRead,

		Schema: map[string]*schema.Schema{
			"description": {
				Description: "Interface description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ip": {
				Description: "Primary interface IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Interface description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"shutdown": {
				Description: "Interface description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vlanid": {
				Description: "VLANID.",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func dataSourceL3VlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.L3Vlan{}
	params.Vlan.Name = id

	resp, err := client.ReadL3Vlan(params)

	if err != nil {
		return diag.Errorf("error retrieving L3 vlans. %s", err)
	}

	dataSetL3Vlan(d, &resp.Vlan)

	d.SetId(strconv.Itoa(id))

	return nil
}

func dataSetL3Vlan(d *schema.ResourceData, resp *models.Vlan) {
	if &resp.Description != nil {
		d.Set("description", resp.Description)
	}
	if &resp.Name != nil {
		d.Set("name", strconv.Itoa(resp.Name))
		d.Set("vlanid", resp.Name)
	}
	if &resp.Shutdown != nil {
		d.Set("shutdown", resp.Shutdown)
	}
	dataSetL3VlanPrimaryIP(d, resp)
}

func dataSetL3VlanPrimaryIP(d *schema.ResourceData, resp *models.Vlan) {
	var primaryIP string
	var primaryMask string
	if &resp.IP.Address.Primary.Address != nil {
		primaryIP = resp.IP.Address.Primary.Address
	}
	if &resp.IP.Address.Primary.Mask != nil {
		primaryMask = resp.IP.Address.Primary.Mask
	}
	if primaryIP != "" && primaryMask != "" {
		d.Set("ip", fmt.Sprintf("%s/%s", primaryIP, primaryMask))
	}
}
