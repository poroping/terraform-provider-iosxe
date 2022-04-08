package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceL2Vlan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a L2 VLAN.",

		CreateContext: resourceL2VlanCreate,
		ReadContext:   resourceL2VlanRead,
		UpdateContext: resourceL2VlanUpdate,
		DeleteContext: resourceL2VlanDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "VLAN name.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"vlanid": {
				Description:  "VLAN ID.",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4096),
			},
		},
	}
}

func resourceL2VlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.L2VlanList{}
	params.VlanList.ID = id

	getCreateUpdateL2VlanObject(d, &params.VlanList)

	err := client.CreateL2Vlan(params)

	if err != nil {
		return diag.Errorf("error creating Vlan. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceL2VlanRead(ctx, d, meta)
}

func resourceL2VlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.L2VlanList{}
	params.VlanList.ID = id

	resp, err := client.ReadL2Vlan(params)

	if err != nil {
		return diag.Errorf("error retrieving Vlan. %s", err)
	}

	resourceSetL2Vlan(d, &resp.VlanList)

	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceL2VlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.L2VlanList{}
	params.VlanList.ID = id

	getCreateUpdateL2VlanObject(d, &params.VlanList)

	err := client.UpdateL2Vlan(params)

	if err != nil {
		return diag.Errorf("error updating Vlan. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceL2VlanRead(ctx, d, meta)
}

func resourceL2VlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	id := d.Get("vlanid").(int)

	params := models.L2VlanList{}
	params.VlanList.ID = id

	err := client.DeleteL2Vlan(params)

	if err != nil {
		return diag.Errorf("error deleting Vlan. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetL2Vlan(d *schema.ResourceData, resp *models.L2Vlan) {
	d.Set("name", resp.Name)
}

func getCreateUpdateL2VlanObject(d *schema.ResourceData, m *models.L2Vlan) *models.L2Vlan {
	if v, ok := d.GetOk("name"); ok {
		if s, ok := v.(string); ok {
			m.Name = s
		}
	}
	return m
}
