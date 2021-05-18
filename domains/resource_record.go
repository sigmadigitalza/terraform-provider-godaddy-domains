package domains

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openapiclient "github.com/sigmadigitalza/godaddy-domains-client"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext: resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": &schema.Schema{
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*openapiclient.APIClient)

	domainName := d.Get("domain").(string)
	data := d.Get("data").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)

	var diags diag.Diagnostics

	records := []openapiclient.DNSRecord{
		*openapiclient.NewDNSRecord(data, name, recordType),
	}

	_, err := client.V1Api.RecordAdd(ctx, domainName).Records(records).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(generateId(domainName, name, recordType))

	resourceRecordRead(ctx, d, m)

	return diags
}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*openapiclient.APIClient)

	domainName := d.Get("domain").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)

	var diags diag.Diagnostics

	dnsRecords, resp, err := client.V1Api.RecordGet(ctx, domainName, recordType, name).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if len(dnsRecords) == 0 {
		d.SetId("")
		return diags
	}

	record := dnsRecords[0]

	return hydrate(diags, record, d)
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*openapiclient.APIClient)

	domainName := d.Get("domain").(string)
	data := d.Get("data").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)

	records := []openapiclient.DNSRecordCreateTypeName{
		*openapiclient.NewDNSRecordCreateTypeName(data),
	}

	_, err := client.V1Api.RecordReplaceTypeName(ctx, domainName, recordType, name).Records(records).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRecordRead(ctx, d, m)
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*openapiclient.APIClient)

	domainName := d.Get("domain").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)

	var diags diag.Diagnostics

	_, err := client.V1Api.RecordDeleteTypeName(ctx, domainName, recordType, name).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func hydrate(diags diag.Diagnostics, record openapiclient.DNSRecord, d *schema.ResourceData) diag.Diagnostics {
	if err := d.Set("data", record.Data); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", record.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", record.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ttl", record.Ttl); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func generateId(domainName string, recordName string, recordType string) string {
	return domainName + ":" + recordName + ":" + recordType
}
