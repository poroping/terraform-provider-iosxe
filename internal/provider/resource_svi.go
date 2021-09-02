package provider

import (
	"context"
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
				Computed:    true,
				Optional:    true,
			},
			"ip": {
				Description:  "Primary interface IP.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"name": {
				Description: "Interface name.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"shutdown": {
				Description: "Interface status.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
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
	d.Set("description", resp.Description)
	d.Set("name", strconv.Itoa(resp.Name))
	d.Set("vlanid", resp.Name)
	d.Set("shutdown", resp.Shutdown)
	resp.IP.Address.Primary.SetCIDR()
	d.Set("ip", resp.IP.Address.Primary.CIDR)
}

func getCreateUpdateSVIObject(d *schema.ResourceData, m *models.SVI) *models.SVI {
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.Vlan.Description = s
		}
	}

	if v, ok := d.GetOk("ip"); ok {
		if s, ok := v.(string); ok {
			m.Vlan.IP.Address.Primary.CIDR = s
			m.Vlan.IP.Address.Primary.SetNetmask()
		}
	}
	return m
}
