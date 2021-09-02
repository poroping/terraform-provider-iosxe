package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a L2 VLAN.",

		CreateContext: resourceBgpNeighborCreate,
		ReadContext:   resourceBgpNeighborRead,
		UpdateContext: resourceBgpNeighborUpdate,
		DeleteContext: resourceBgpNeighborDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "Cluster ID.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"description": {
				Description: "Description.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"disable_connected_check": {
				Description: "Whatever this means.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"ip": {
				Description:  "Neighbor IP.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"remote_as": {
				Description: "Remote AS.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"shutdown": {
				Description: "Shutdown.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
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

func resourceBgpNeighborCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	params := models.BgpNeighbor{}
	params.Neighbor.ID = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	getCreateUpdateBgpNeighborObject(d, &params)

	err := client.CreateBgpNeighbor(router, params)

	if err != nil {
		return diag.Errorf("error creating BgpNeighbor. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborRead(ctx, d, meta)
}

func resourceBgpNeighborRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	params := models.BgpNeighbor{}
	params.Neighbor.ID = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	resp, err := client.ReadBgpNeighbor(router, params)

	if err != nil {
		return diag.Errorf("error retrieving BgpNeighbor. %s", err)
	}

	resourceSetBgpNeighbor(d, resp)

	d.SetId(id)

	return nil
}

func resourceBgpNeighborUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	params := models.BgpNeighbor{}
	params.Neighbor.ID = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	getCreateUpdateBgpNeighborObject(d, &params)

	err := client.UpdateBgpNeighbor(router, params)

	if err != nil {
		return diag.Errorf("error updating BgpNeighbor. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborRead(ctx, d, meta)
}

func resourceBgpNeighborDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	params := models.BgpNeighbor{}
	params.Neighbor.ID = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	err := client.DeleteBgpNeighbor(router, params)

	if err != nil {
		return diag.Errorf("error deleting BgpNeighbor. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetBgpNeighbor(d *schema.ResourceData, resp *models.BgpNeighbor) {
	d.Set("cluster_id", resp.Neighbor.ClusterID)
	d.Set("description", resp.Neighbor.Description)
	d.Set("disable_connected_check", resp.Neighbor.DisableConnectedCheck)
	d.Set("remote_as", resp.Neighbor.RemoteAs)
	d.Set("shutdown", resp.Neighbor.Shutdown)
}

func getCreateUpdateBgpNeighborObject(d *schema.ResourceData, m *models.BgpNeighbor) *models.BgpNeighbor {
	if v, ok := d.GetOk("cluster_id"); ok {
		if s, ok := v.(string); ok {
			m.Neighbor.ClusterID = s
		}
	}
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.Neighbor.Description = s
		}
	}
	if v, ok := d.GetOk("disable_connected_check"); ok {
		if s, ok := v.(string); ok {
			m.Neighbor.DisableConnectedCheck = s
		}
	}
	if v, ok := d.GetOk("remote_as"); ok {
		if i, ok := v.(int); ok {
			m.Neighbor.RemoteAs = i
		}
	}
	if v, ok := d.GetOk("shutdown"); ok {
		if s, ok := v.(string); ok {
			m.Neighbor.Shutdown = s
		}
	}
	return m
}
