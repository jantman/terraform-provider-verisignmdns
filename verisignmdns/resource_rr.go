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
                        "account_id": &schema.Schema{
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "zone_name": &schema.Schema{
                                Type:     schema.TypeString,
                                Optional: true,
                        },
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
  return resourceRrRead(d, m)
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
        return resourceRrRead(d, m)
}

func resourceRrDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
