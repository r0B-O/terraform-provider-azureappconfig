package azureappconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConfigurationSetting() *schema.Resource {
	return &schema.Resource{
		Create: onCreate,
		Read:   onRead,
		Update: onUpdate,
		Delete: onDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"config_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func onCreate(d *schema.ResourceData, m interface{}) error {
	return do("create", d, m)
}

func onRead(d *schema.ResourceData, m interface{}) error {
	return do("read", d, m)
}

func onUpdate(d *schema.ResourceData, m interface{}) error {
	return do("update", d, m)
}

func onDelete(d *schema.ResourceData, m interface{}) error {
	return do("delete", d, m)
}

func do(event string, d *schema.ResourceData, m interface{}) error {
	log.Printf("Executing: %s %s %s %s", d.Get("config_name"), d.Get("key"), event, d.Get("value"))

	if event == "create" {
		azureCommand := fmt.Sprintf("az appconfig kv set -n %s --key %s --value %s --label %s --yes", d.Get("config_name"), d.Get("key"), d.Get("value"), d.Get("label"))
		cmd := exec.Command("bash", "-c", azureCommand)

		result, err := cmd.Output()

		if err == nil {
			var resource map[string]interface{}
			err = json.Unmarshal([]byte(result), &resource)
			if err == nil {
				if event == "delete" {
					d.SetId("")
				} else {
					key := d.Get("id_key").(string)
					d.Set("resource", resource)
					d.SetId(resource[key].(string))
				}
			}
		}
		return err
	}

	if event == "read" {
		azureCommand := fmt.Sprintf("az appconfig kv list --all")
		cmd := exec.Command("bash", "-c", azureCommand)

		result, err := cmd.Output()

		if err == nil {
			var resource map[string]interface{}
			err = json.Unmarshal([]byte(result), &resource)
			if err == nil {
				if event == "delete" {
					d.SetId("")
				} else {
					key := d.Get("id_key").(string)
					d.Set("resource", resource)
					d.SetId(resource[key].(string))
				}
			}
		}
		return err
	}

	if event == "update" {
		azureCommand := fmt.Sprintf("az appconfig kv set -n %s --key %s --value %s --label %s --yes", d.Get("config_name"), d.Get("key"), d.Get("value"), d.Get("label"))
		cmd := exec.Command("bash", "-c", azureCommand)

		result, err := cmd.Output()

		if err == nil {
			var resource map[string]interface{}
			err = json.Unmarshal([]byte(result), &resource)
			if err == nil {
				if event == "delete" {
					d.SetId("")
				} else {
					key := d.Get("id_key").(string)
					d.Set("resource", resource)
					d.SetId(resource[key].(string))
				}
			}
		}
		return err

	}

	if event == "delete" {
		azureCommand := fmt.Sprintf("az appconfig kv delete --key %s", d.Get("key"))
		cmd := exec.Command("bash", "-c", azureCommand)

		result, err := cmd.Output()

		if err == nil {
			var resource map[string]interface{}
			err = json.Unmarshal([]byte(result), &resource)
			if err == nil {
				if event == "delete" {
					d.SetId("")
				} else {
					key := d.Get("id_key").(string)
					d.Set("resource", resource)
					d.SetId(resource[key].(string))
				}
			}
		}
		return err
	}
	return nil
}
