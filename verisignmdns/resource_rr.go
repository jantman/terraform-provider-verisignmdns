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
                        "account_id": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "zone_name": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "record_name": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "record_type": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "record_data": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                },
        }
}

func resourceRrCreate(d *schema.ResourceData, m interface{}) error {
  client := m.(*api_client)
  data, err := client.create_rr(
    d.Get("account_id").(string),
    d.Get("zone_name").(string),
    d.Get("record_name").(string),
    d.Get("record_type").(string),
    d.Get("record_data").(string),
  )
  if err != nil {
    return err
  }
  d.SetId(data["resource_record_id"].(string))
  return resourceRrRead(d, m)
}

func resourceRrRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*api_client)
  data, err := client.get_rr(
    d.Get("accountId").(string),
    d.Get("zoneName").(string),
    d.Id(),
  )
  if err != nil {
    return err
  }
  d.Set("recordName", data["owner"].(string))
  d.Set("recordType", data["type"].(string))
  d.Set("recordData", data["rdata"].(string))
  return nil
}

func resourceRrUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceRrRead(d, m)
}

func resourceRrDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
