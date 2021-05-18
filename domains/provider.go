package domains

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openapiclient "github.com/sigmadigitalza/godaddy-domains-client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_HOST", "https://api.ote-godaddy.com"),
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_KEY", ""),
			},
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_SECRET", ""),
			},
		},
		ConfigureContextFunc: configureContext,
		ResourcesMap:         map[string]*schema.Resource{
			"domains_record": resourceRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"domains_domain": dataSourceDomain(),
		},
	}
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	key := d.Get("key").(string)
	secret := d.Get("secret").(string)

	var diags diag.Diagnostics

	if (key == "") || (secret == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create GoDaddy Domains client",
			Detail:   "Missing credentials for the GoDaddy Domains client",
		})
	}

	configuration := openapiclient.Configuration{
		DefaultHeader: map[string]string{
			"Authorization": "sso-key " + key + ":" + secret,
		},
		UserAgent: "OpenAPI-Generator/1.0.0/go",
		Debug:     false,
		Servers: openapiclient.ServerConfigurations{
			{
				URL:         host,
				Description: "No description provided",
			},
		},
		OperationServers: map[string]openapiclient.ServerConfigurations{
		},
	}

	client := openapiclient.NewAPIClient(&configuration)

	return client, diags
}
