package domains

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openapiclient "github.com/sigmadigitalza/godaddy-domains-client"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*openapiclient.APIClient)
	domainName := d.Get("domain").(string)

	var diags diag.Diagnostics

	domainDetail, resp, err := client.V1Api.Get(ctx, domainName).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if err := d.Set("domain", domainDetail.Domain); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%f", domainDetail.DomainId))

	return diags
}
