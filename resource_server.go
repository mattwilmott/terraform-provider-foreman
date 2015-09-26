package main

import (
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"location_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"environment_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"architecture_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"realm_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"puppet_proxy_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"puppet_class_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"operatingsystem_id": &schema.Schema{
				Type:     schema.TypeString, //Why isnt this an Int? API doco may be incorrect
				Optional: true,
			},
			"medium_id": &schema.Schema{
				Type:     schema.TypeString, //Why isnt this an Int as well? wtf
				Optional: true,
			},
			"ptable_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"compute_resource_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"root_pass": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"model_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hostgroup_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"puppet_ca_proxy_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_parameters_attributes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"build": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"provision_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"provision_report_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"capabilities": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"compute_profile_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"interface_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"compute_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

type hammerArgs struct {
	options    []string
	subcommand string
	args       []string
}

func hammerCLI(h *hammerArgs) (output []byte, err error) {
	cmd := exec.Command("/usr/bin/hammer", h.subcommand, "list")
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(stdout))
	return stdout, err
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	h := hammerArgs{subcommand: "host"}
	output, err := hammerCLI(&h)
	if err != nil {
		panic(err)
	}
	print(string(output))
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
