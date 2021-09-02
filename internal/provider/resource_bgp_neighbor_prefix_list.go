package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/models"
)

func resourceBgpNeighborPrefixList() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a BGP Neighbor prefix-list.",

		CreateContext: resourceBgpNeighborPrefixListCreate,
		ReadContext:   resourceBgpNeighborPrefixListRead,
		UpdateContext: resourceBgpNeighborPrefixListUpdate,
		DeleteContext: resourceBgpNeighborPrefixListDelete,

		Schema: map[string]*schema.Schema{
			"direction": {
				Description:  "Direction of prefix-list application.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"in", "out"}, false),
			},
			"name": {
				Description: "Prefix-list name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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

func resourceBgpNeighborPrefixListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("name").(string)

	params := models.BgpNeighborPrefixList{}
	params.PrefixList.PrefixListName = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	getCreateUpdateBgpNeighborPrefixListObject(d, &params)

	err := client.CreateBgpNeighborPrefixList(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error creating BgpNeighborPrefixList. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborPrefixListRead(ctx, d, meta)
}

func resourceBgpNeighborPrefixListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("name").(string)

	params := models.BgpNeighborPrefixList{}
	params.PrefixList.PrefixListName = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	resp, err := client.ReadBgpNeighborPrefixList(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error retrieving BgpNeighborPrefixList. %s", err)
	}

	resourceSetBgpNeighborPrefixList(d, resp)

	d.SetId(id)

	return nil
}

func resourceBgpNeighborPrefixListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("name").(string)

	params := models.BgpNeighborPrefixList{}
	params.PrefixList.PrefixListName = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	getCreateUpdateBgpNeighborPrefixListObject(d, &params)

	err := client.UpdateBgpNeighborPrefixList(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error updating BgpNeighborPrefixList. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborPrefixListRead(ctx, d, meta)
}

func resourceBgpNeighborPrefixListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("name").(string)

	params := models.BgpNeighborPrefixList{}
	params.PrefixList.PrefixListName = id

	router := models.BgpRouter{}
	router.Bgp.ID = d.Get("routing_instance").(string)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = d.Get("neighbor_ip").(string)

	err := client.DeleteBgpNeighborPrefixList(router, neighbor, params)

	if err != nil {
		return diag.Errorf("error deleting BgpNeighborPrefixList. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetBgpNeighborPrefixList(d *schema.ResourceData, resp *models.BgpNeighborPrefixList) {
	d.Set("direction", resp.PrefixList.Inout)
	d.Set("name", resp.PrefixList.PrefixListName)
}

func getCreateUpdateBgpNeighborPrefixListObject(d *schema.ResourceData, m *models.BgpNeighborPrefixList) *models.BgpNeighborPrefixList {
	if v, ok := d.GetOk("direction"); ok {
		if s, ok := v.(string); ok {
			m.PrefixList.Inout = s
		}
	}
	return m
}
