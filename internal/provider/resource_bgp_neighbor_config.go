package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceBgpNeighborConfig() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a BGP Neighbor config.",

		CreateContext: resourceBgpNeighborConfigCreate,
		ReadContext:   resourceBgpNeighborConfigRead,
		UpdateContext: resourceBgpNeighborConfigUpdate,
		DeleteContext: resourceBgpNeighborConfigDelete,

		Schema: map[string]*schema.Schema{
			"default_originate": {
				Description: "Originate default route.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"remove_private_as": {
				Description: "Remove private AS.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"soft_reconfiguration": {
				Description: "Soft reconfiguration.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"neighbor_ip": {
				Description:  "Neighbor IP.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"routing_instance": {
				Description: "Routing instance.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBgpNeighborConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("neighbor_ip").(string)

	params := models.BgpNeighborConfig{}

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	getCreateUpdateBgpNeighborConfigObject(d, &params)

	err := client.CreateBgpNeighborConfig(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error creating BgpNeighborConfig. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborConfigRead(ctx, d, meta)
}

func resourceBgpNeighborConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("neighbor_ip").(string)

	params := models.BgpNeighborConfig{}

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	resp, err := client.ReadBgpNeighborConfig(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error retrieving BgpNeighborConfig. %s", err)
	}

	resourceSetBgpNeighborConfig(d, resp)

	d.SetId(id)

	return nil
}

func resourceBgpNeighborConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("neighbor_ip").(string)

	params := models.BgpNeighborConfig{}

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	getCreateUpdateBgpNeighborConfigObject(d, &params)

	err := client.UpdateBgpNeighborConfig(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error updating BgpNeighborConfig. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborConfigRead(ctx, d, meta)
}

func resourceBgpNeighborConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("neighbor_ip").(string)

	params := models.BgpNeighborConfig{}

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = id

	err := client.DeleteBgpNeighborConfig(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error deleting BgpNeighborConfig. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetBgpNeighborConfig(d *schema.ResourceData, resp *models.BgpNeighborConfig) {
	d.Set("direction", resp.Config.Inout)
	d.Set("name", resp.Config.ConfigName)
}

func getCreateUpdateBgpNeighborConfigObject(d *schema.ResourceData, m *models.BgpNeighborConfig) *models.BgpNeighborConfig {
	if v, ok := d.GetOk("direction"); ok {
		if s, ok := v.(string); ok {
			m.Config.Inout = s
		}
	}
	return m
}
