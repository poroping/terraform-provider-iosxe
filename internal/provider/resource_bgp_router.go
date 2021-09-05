package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceBgpRouter() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a BGP router AS.",

		CreateContext: resourceBgpRouterCreate,
		ReadContext:   resourceBgpRouterRead,
		UpdateContext: resourceBgpRouterUpdate,
		DeleteContext: resourceBgpRouterDelete,

		Schema: map[string]*schema.Schema{
			"log_neighbor_changes": {
				Description: "Log neighbor changes.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"as": {
				Description: "Autonomous system number.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBgpRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("as").(int)

	params := models.BgpRouter{}
	params.Bgp.ID = id

	getCreateUpdateBgpRouterObject(d, &params)

	err := client.CreateBgpRouter(params)

	if err != nil {
		return diag.Errorf("error creating BgpRouter. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceBgpRouterRead(ctx, d, meta)
}

func resourceBgpRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("as").(int)

	params := models.BgpRouter{}
	params.Bgp.ID = id

	resp, err := client.ReadBgpRouter(params)

	if err != nil {
		return diag.Errorf("error retrieving BgpRouter. %s", err)
	}

	resourceSetBgpRouter(d, resp)

	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceBgpRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("as").(int)

	params := models.BgpRouter{}
	params.Bgp.ID = id

	getCreateUpdateBgpRouterObject(d, &params)

	err := client.UpdateBgpRouter(params)

	if err != nil {
		return diag.Errorf("error updating BgpRouter. %s", err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceBgpRouterRead(ctx, d, meta)
}

func resourceBgpRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("as").(int)

	params := models.BgpRouter{}
	params.Bgp.ID = id

	err := client.DeleteBgpRouter(params)

	if err != nil {
		return diag.Errorf("error deleting BgpRouter. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetBgpRouter(d *schema.ResourceData, resp *models.BgpRouter) {
	d.Set("as", resp.Bgp.ID)
	d.Set("log_neighbor_changes", resp.Bgp.Bgp.LogNeighborChanges)
}

func getCreateUpdateBgpRouterObject(d *schema.ResourceData, m *models.BgpRouter) *models.BgpRouter {
	if v, ok := d.GetOk("log_neighbor_changes"); ok {
		if b, ok := v.(bool); ok {
			m.Bgp.Bgp.LogNeighborChanges = b
		}
	}
	return m
}
