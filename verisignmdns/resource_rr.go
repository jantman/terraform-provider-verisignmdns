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

                Schema: map[string]*schema.Schema{
                        "address": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                },
        }
}

func resourceRrCreate(d *schema.ResourceData, m interface{}) error {
        return resourceRrRead(d, m)
}

func resourceRrRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceRrUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceRrRead(d, m)
}

func resourceRrDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
