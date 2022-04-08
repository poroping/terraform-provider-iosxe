package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poroping/go-ios-xe-sdk/client"
	"github.com/poroping/go-ios-xe-sdk/config"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("TF_IOSXE_HOST", nil),
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("TF_IOSXE_USERNAME", nil),
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("TF_IOSXE_PASSWORD", nil),
				},
				"insecure": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("TF_IOSXE_INSECURE", false),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				// "iosxe_interface_vlan": dataSourceVlan(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"iosxe_interface_port_channel":              resourcePortChannel(),
				"iosxe_interface_port_channel_subinterface": resourcePortChannelSubinterface(),
				"iosxe_interface_vlan":                      resourceVlan(),
				// "iosxe_l3_interface":                        resourceL3Interface(),
				"iosxe_l2_vlan":      resourceL2Vlan(),
				"iosxe_bgp_router":   resourceBgpRouter(),
				"iosxe_bgp_neighbor": resourceBgpNeighbor(),
				"iosxe_vrf":          resourceVRF(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	Client *client.CiscoIOSXEClient
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		host := d.Get("host").(string)
		insecure := d.Get("insecure").(bool)
		userAgent := p.UserAgent("terraform-provider-iosxe", version)
		var diags diag.Diagnostics
		cfg := config.Config{
			Username:  username,
			Password:  password,
			Host:      host,
			Insecure:  insecure,
			UserAgent: userAgent,
		}
		c, err := client.NewClient(cfg)
		diags = append(diags, diag.FromErr(err)...)
		apiClient := &apiClient{
			Client: c,
		}

		return apiClient, diags
	}
}
