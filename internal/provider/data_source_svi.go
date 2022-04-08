package provider

// import (
// 	"context"
// 	"strconv"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/poroping/go-ios-xe-sdk/client"
// 	"github.com/poroping/go-ios-xe-sdk/models"
// )

// func dataSourceSVI() *schema.Resource {
// 	return &schema.Resource{
// 		// This description is used by the documentation generator and the language server.
// 		Description: "Sample data source in the Terraform provider SVI.",

// 		ReadContext: dataSourceSVIRead,

// 		Schema: map[string]*schema.Schema{
// 			"description": {
// 				Description: "Interface description.",
// 				Type:        schema.TypeString,
// 				Computed:    true,
// 			},
// 			"ip": {
// 				Description: "Primary interface IP.",
// 				Type:        schema.TypeString,
// 				Computed:    true,
// 			},
// 			"name": {
// 				Description: "Interface name.",
// 				Type:        schema.TypeString,
// 				Computed:    true,
// 			},
// 			"shutdown": {
// 				Description: "Interface status.",
// 				Type:        schema.TypeBool,
// 				Computed:    true,
// 			},
// 			"vlanid": {
// 				Description: "VLANID.",
// 				Type:        schema.TypeInt,
// 				Required:    true,
// 			},
// 		},
// 	}
// }

// func dataSourceSVIRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*client.Client)
// 	id := d.Get("vlanid").(int)

// 	params := models.SVI{}
// 	params.Vlan.Name = id

// 	resp, err := client.ReadSVI(params)

// 	if err != nil {
// 		return diag.Errorf("error retrieving SVI. %s", err)
// 	}

// 	dataSetSVI(d, &resp.Vlan)

// 	d.SetId(strconv.Itoa(id))

// 	return nil
// }

// func dataSetSVI(d *schema.ResourceData, resp *models.L3Interface) {
// 	d.Set("description", resp.Description)
// 	d.Set("name", resp.Name)
// 	d.Set("vlanid", resp.Name)
// 	d.Set("shutdown", resp.Shutdown)
// 	resp.IP.Address.Primary.SetCIDR()
// 	d.Set("ip", resp.IP.Address.Primary.CIDR)
// }
