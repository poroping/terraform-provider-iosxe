package provider

import (
	"context"
	"encoding/json"

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
			"as": {
				Description: "Autonomous system number.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			// "cluster_id": {
			// 	Description: "Cluster ID.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			"default_originate": {
				Description: "Originate default route.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Default:     nil,
			},
			"description": {
				Description: "Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"disable_connected_check": {
				Description: "Whatever this means.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Default:     nil,
			},
			"ebgp_multihop": {
				Description: "EBG multi-hop.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ip": {
				Description:  "Neighbor IP.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"local_as": {
				Description: "Local AS.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"prefix_list": {
				Description: "Prefix-list settings.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    2,
				Elem: &schema.Resource{
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
						},
					},
				},
			},
			"remote_as": {
				Description: "Remote AS.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"remove_private_as": {
				Description: "Remove private AS.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Default:     nil,
			},
			"shutdown": {
				Description: "Shutdown.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Default:     nil,
			},
			"soft_reconfiguration": {
				Description: "Soft reconfiguration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"timers": {
				Description: "BGP timers.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keepalive_interval": {
							Description: "Keepalive interval.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"holdtime": {
							Description: "Hold down time.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"minimum_neighbor_hold": {
							Description: "Min hold time from neighbor.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			// "update_source": {
			// 	Description: "Update source interface.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			"vrf": {
				Description: "VRF.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func resourceBgpNeighborCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	as := d.Get("as").(int)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = id

	// create neighbor
	getCreateUpdateBgpNeighborObject(d, &neighbor)

	err := client.CreateBgpNeighbor(as, neighbor)

	if err != nil {
		return diag.Errorf("error creating BgpNeighbor. %s", err)
	}

	// create neighbor config
	neighborConf := models.BgpNeighborConfig{}
	neighborConf.NeighborConfig.ID = id
	getCreateUpdateBgpNeighborConfigObject(d, &neighborConf)

	err = client.CreateBgpNeighborConfig(as, neighborConf)

	if err != nil {
		return diag.Errorf("error creating BgpNeighborConfig. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborRead(ctx, d, meta)
}

func resourceBgpNeighborRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	as := d.Get("as").(int)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = id

	// read neighbor
	resp, err := client.ReadBgpNeighbor(as, neighbor)

	if err != nil {
		return diag.Errorf("error retrieving BgpNeighbor. %s", err)
	}

	resourceSetBgpNeighbor(d, resp)

	// read neighbor config
	neighborConf := models.BgpNeighborConfig{}
	neighborConf.NeighborConfig.ID = id
	// getCreateUpdateBgpNeighborConfigObject(d, &neighborConf)

	resp2, err := client.ReadBgpNeighborConfig(as, neighborConf)

	if err != nil {
		return diag.Errorf("error retrieving BgpNeighborConfig. %s", err)
	}

	// todo: handle 404s

	resourceSetBgpNeighborConfig(d, resp2)

	d.SetId(id)

	return nil
}

func resourceBgpNeighborUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	as := d.Get("as").(int)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = id

	// update neighbor
	getCreateUpdateBgpNeighborObject(d, &neighbor)

	err := client.CreateBgpNeighbor(as, neighbor)

	if err != nil {
		return diag.Errorf("error updating BgpNeighbor. %s", err)
	}

	// update neighbor config
	neighborConf := models.BgpNeighborConfig{}
	neighborConf.NeighborConfig.ID = id
	getCreateUpdateBgpNeighborConfigObject(d, &neighborConf)

	// err = client.UpdateBgpNeighborConfig(as, neighborConf)
	err = client.CreateBgpNeighborConfig(as, neighborConf)

	if err != nil {
		return diag.Errorf("error updating BgpNeighborConfig. %s", err)
	}

	d.SetId(id)

	return resourceBgpNeighborRead(ctx, d, meta)
}

func resourceBgpNeighborDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id := d.Get("ip").(string)

	as := d.Get("as").(int)

	neighbor := models.BgpNeighbor{}
	neighbor.Neighbor.ID = id

	// delete neighbor (we can just remove parent here it will remove all child config)
	err := client.DeleteBgpNeighbor(as, neighbor)

	if err != nil {
		return diag.Errorf("error deleting BgpNeighbor. %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSetBgpNeighbor(d *schema.ResourceData, resp *models.BgpNeighbor) {
	// d.Set("cluster_id", resp.Neighbor.ClusterID)
	if resp.Neighbor.Description != nil {
		d.Set("description", resp.Neighbor.Description)
	}
	if resp.Neighbor.DisableConnectedCheck != nil {
		d.Set("disable_connected_check", true)
	} else {
		d.Set("disable_connected_check", false)
	}
	if resp.Neighbor.EbgpMultihop != nil {
		d.Set("ebgp_multihop", resp.Neighbor.EbgpMultihop.MaxHop)
	}
	if resp.Neighbor.LocalAs != nil {
		d.Set("local_as", resp.Neighbor.LocalAs.AsNo)
	}
	d.Set("remote_as", resp.Neighbor.RemoteAs)
	if resp.Neighbor.Shutdown != nil {
		d.Set("shutdown", true)
	} else {
		d.Set("shutdown", false)
	}
	d.Set("timers", flattenBgpNeighborTimers(resp.Neighbor.Timers))
	// if resp.Neighbor.UpdateSource != nil {
	// 	d.Set("update_source", "TOBEIMPLEMENTED") // do something here cause interface name is variable key
	// }
}

func resourceSetBgpNeighborConfig(d *schema.ResourceData, resp *models.BgpNeighborConfig) {
	if resp.NeighborConfig.DefaultOriginate != nil {
		d.Set("default_originate", true)
	} else {
		d.Set("default_originate", false)
	}
	// d.Set("local_as", &resp.NeighborConfig.LocalAs.AsNo)
	if resp.NeighborConfig.RemovePrivateAs != nil {
		d.Set("remove_private_as", true)
	} else {
		d.Set("remove_private_as", false)
	}
	d.Set("prefix_list", flattenBgpNeighborConfigPrefixList(&resp.NeighborConfig.PrefixList))
	if resp.NeighborConfig.SoftReconfiguration != nil {
		d.Set("soft_reconfiguration", resp.NeighborConfig.SoftReconfiguration)
	}
}

func flattenBgpNeighborTimers(input *models.Timers) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	if l := input; l != nil {
		for _, v := range []models.Timers{*input} {
			output := map[string]interface{}{}
			output["keepalive-interval"] = v.KeepaliveInterval
			output["holdtime"] = v.Holdtime
			output["minimum-neighbor-hold"] = v.MinimumNeighborHold
			results = append(results, output)
		}
	}
	return results
}

func flattenBgpNeighborConfigPrefixList(input *[]models.PrefixList) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	if l := input; l != nil {
		for _, v := range *input {
			output := map[string]interface{}{}
			output["direction"] = v.Inout
			output["name"] = v.PrefixListName
			results = append(results, output)
		}
	}
	return results
}

func getCreateUpdateBgpNeighborObject(d *schema.ResourceData, m *models.BgpNeighbor) *models.BgpNeighbor {
	// if v, ok := d.GetOk("cluster_id"); ok {
	// 	if s, ok := v.(string); ok {
	// 		m.Neighbor.ClusterID = s
	// 	}
	// }
	if v, ok := d.GetOk("description"); ok {
		if s, ok := v.(string); ok {
			m.Neighbor.Description = &s
		}
	}
	if v, ok := d.GetOk("disable_connected_check"); ok {
		if b, ok := v.(bool); ok {
			if b {
				n := json.RawMessage(models.ExplicitNull)
				m.Neighbor.DisableConnectedCheck = &n
			}
		}
	}
	if v, ok := d.GetOk("ebgp_multihop"); ok {
		if i, ok := v.(int); ok {
			m.Neighbor.EbgpMultihop.MaxHop = &i
		}
	}
	if v, ok := d.GetOk("local_as"); ok {
		if i, ok := v.(int); ok {
			m.Neighbor.LocalAs.AsNo = &i
		}
	}
	if v, ok := d.GetOk("remote_as"); ok {
		if i, ok := v.(int); ok {
			m.Neighbor.RemoteAs = i
		}
	}
	if v, ok := d.GetOk("shutdown"); ok {
		if b, ok := v.(bool); ok {
			if b {
				n := json.RawMessage(models.ExplicitNull)
				m.Neighbor.Shutdown = &n
			}
		}
	}
	if _, ok := d.GetOk("timers"); ok {
		o := expandBgpNeighborTimers(d, "timers")
		m.Neighbor.Timers = o
	}
	// if _, ok := d.GetOk("update_source"); ok {
	// 	o := expandBgpNeighborUpdateSource(d, "update_source")
	// 	m.Neighbor.UpdateSource = o
	// }
	return m
}

func getCreateUpdateBgpNeighborConfigObject(d *schema.ResourceData, m *models.BgpNeighborConfig) *models.BgpNeighborConfig {
	if v, ok := d.GetOk("default_originate"); ok {
		if b, ok := v.(bool); ok {
			if b {
				m.NeighborConfig.DefaultOriginate = &struct{}{}
			}
		}
	}
	if v, ok := d.GetOk("local_as"); ok {
		if i, ok := v.(int); ok {
			m.NeighborConfig.LocalAs.AsNo = i
		}
	}
	if _, ok := d.GetOk("prefix_list"); ok {
		o := expandBgpNeighborConfigPrefixList(d, "prefix_list")
		m.NeighborConfig.PrefixList = *o
	}
	if v, ok := d.GetOk("remove_private_as"); ok {
		if b, ok := v.(bool); ok {
			if b {
				m.NeighborConfig.RemovePrivateAs = &struct{}{}
			}
		}
	}
	if v, ok := d.GetOk("soft_reconfiguration"); ok {
		if s, ok := v.(string); ok {
			m.NeighborConfig.SoftReconfiguration = &s
		}
	}
	return m
}

func expandBgpNeighborConfigPrefixList(d *schema.ResourceData, field string) *[]models.PrefixList {
	l := d.Get(field).([]interface{})
	if len(l) == 0 {
		return nil
	}

	r := make([]models.PrefixList, 0, len(l))
	// probs could json.marshal but renamed fields cause dumb names
	for _, v := range l {
		temp := models.PrefixList{}
		m := v.(map[string]interface{})
		temp.Inout = m["direction"].(string)
		temp.PrefixListName = m["name"].(string)

		r = append(r, temp)
	}

	return &r
}

func expandBgpNeighborTimers(d *schema.ResourceData, field string) *models.Timers {
	l := d.Get(field).([]interface{})
	if len(l) == 0 {
		return nil
	}

	r := make([]models.Timers, 0, len(l))
	for _, v := range l {
		temp := models.Timers{}
		m := v.(map[string]interface{})
		temp.KeepaliveInterval = m["keepalive_interval"].(int)
		temp.Holdtime = m["holdtime"].(int)
		temp.MinimumNeighborHold = m["minimum-neighbor-hold"].(int)
		r = append(r, temp)
	}

	return &r[0]
}

// func getCreateUpdateBgpNeighborPrefixListObject(d *schema.ResourceData, m *models.BgpNeighborPrefixList) *models.BgpNeighborPrefixList {
// 	if v, ok := d.GetOk("prefix_list.direction"); ok {
// 		if s, ok := v.(string); ok {
// 			m.PrefixList.Inout = s
// 		}
// 	}
// 	if v, ok := d.GetOk("prefix_list.name"); ok {
// 		if s, ok := v.(string); ok {
// 			m.PrefixList.PrefixListName = s
// 		}
// 	}
// 	return m
// }
