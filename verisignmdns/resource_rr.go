package verisignmdns

import (
        "github.com/hashicorp/terraform/helper/schema"
)

func resourceRr() *schema.Resource {
        return &schema.Resource{
                Create: resourceRrCreate,
                Read:   resourceRrRead,
                Update: resourceRrUpdate,
                Delete: resourceRrDelete,

                SchemaVersion: 1,
                Schema: map[string]*schema.Schema{
                        "record_name": &schema.Schema{
                                Type:        schema.TypeString,
                                Required:    true,
                                Description: "The name of the record, i.e. the FQDN without a trailing dot.",
                                ForceNew:    true,
                        },
                        "record_type": &schema.Schema{
                                Type:        schema.TypeString,
                                Required:    true,
                                Description: "The type of record, i.e. A, AAAA, CNAME, etc.",
                                ForceNew:    true,
                        },
                        "record_data": &schema.Schema{
                                Type:        schema.TypeString,
                                Required:    true,
                                Description: "The value of the record.",
                        },
                        "account_id": &schema.Schema{
                                Type:        schema.TypeString,
                                Description: "The Account ID the record exists in (configured as part of the provider).",
                                Computed:    true,
                        },
                        "zone_name": &schema.Schema{
                                Type:        schema.TypeString,
                                Description: "The Zone Name the record exists in (configured as part of the provider).",
                                Computed:    true,
                        },
                },
                Importer: &schema.ResourceImporter{
                    State: schema.ImportStatePassthrough,
                },
        }
}

func resourceRrCreate(d *schema.ResourceData, m interface{}) error {
  client := m.(*api_client)
  data, err := client.create_rr(
    d.Get("record_name").(string),
    d.Get("record_type").(string),
    d.Get("record_data").(string),
  )
  if err != nil {
    return err
  }
  d.SetId(data["resource_record_id"].(string))
  d.Set("recordName", data["owner"].(string))
  d.Set("recordType", data["type"].(string))
  d.Set("recordData", data["rdata"].(string))
  d.Set("zone_name", client.zone_name)
  d.Set("account_id", client.account_id)
  return nil
}

func resourceRrRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*api_client)
  data, err := client.get_rr(
    d.Id(),
  )
  if err != nil {
    return err
  }
  d.Set("recordName", data["owner"].(string))
  d.Set("recordType", data["type"].(string))
  d.Set("recordData", data["rdata"].(string))
  d.Set("zone_name", client.zone_name)
  d.Set("account_id", client.account_id)
  return nil
}

func resourceRrUpdate(d *schema.ResourceData, m interface{}) error {
  client := m.(*api_client)
  data, err := client.update_rr(
    d.Id(),
    d.Get("record_data").(string),
  )
  if err != nil {
    return err
  }
  d.SetId(data["resource_record_id"].(string))
  d.Set("recordName", data["owner"].(string))
  d.Set("recordType", data["type"].(string))
  d.Set("recordData", data["rdata"].(string))
  d.Set("zone_name", client.zone_name)
  d.Set("account_id", client.account_id)
  return nil
}

func resourceRrDelete(d *schema.ResourceData, m interface{}) error {
  client := m.(*api_client)
  err := client.delete_rr(
    d.Id(),
  )
  if err != nil {
    return err
  }
  d.SetId("")
  return nil
}
