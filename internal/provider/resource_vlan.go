package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a L2 VLAN.",

		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "VLAN name.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"vlanid": {
				Description: "VLAN ID.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.VlanList.ID = id

	getCreateUpdateVlanObject(d, &params)

	err := client.CreateVlan(params)

	if err != nil {
		return diag.Errorf("error creating Vlan. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceVlanRead(ctx, d, meta)
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.VlanList.ID = id

	resp, err := client.ReadVlan(params)

	if err != nil {
		return diag.Errorf("error retrieving Vlan. %s", err)
	}

	resourceSetVlan(d, &resp.VlanList)

	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.VlanList.ID = id

	getCreateUpdateVlanObject(d, &params)

	err := client.UpdateVlan(params)

	if err != nil {
		return diag.Errorf("error updating Vlan. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceVlanRead(ctx, d, meta)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("vlanid").(int)

	params := models.Vlan{}
	params.VlanList.ID = id

	err := client.DeleteVlan(params)

	if err != nil {
		return diag.Errorf("error deleting Vlan. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetVlan(d *schema.ResourceData, resp *models.VlanList) {
	d.Set("name", resp.Name)
}

func getCreateUpdateVlanObject(d *schema.ResourceData, m *models.Vlan) *models.Vlan {
	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			m.VlanList.Name = s
		}
	}
	return m
}
