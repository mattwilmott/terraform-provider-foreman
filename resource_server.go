package main

import (
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
        "strings"
	"reflect"
	"fmt"
	"net/http"
)

type host struct {
  name				string
  environment_id		string
  ip				string
  mac				string
  architecture_id		int
  domain_id			int
  realm_id			int
  puppet_proxy_id		int
  puppetclass_ids		[]int
  operatingsystem_id		string
  medium_id			string
  ptable_id			int
  subnet_id			int
  compute_resource_id		int
  root_pass			string
  model_id			int
  hostgroup_id			int
  owner_id			int
  owner_type			string // must be either User or Usergroup
  puppet_ca_proxy_id		int
  image_id			int
  host_parameters_attributes	{
				  var roles string
				  var puppet string
				  var chef string
				  var JIRA_Ticket string
				}
  build				bool
  enabled			bool
  provision_method		string
  managed			bool
  progress_report_id		string
  comment			string
  capabilities			string
  compute_profile_id		int
  interfaces_attributes		{
				  var mac string
				  var ip string
				  var type string
				  var name string
				  var subnet_id int
				  var domain_id int
				  var identifier string
				  var managed bool
				  var primary bool
				  var provision bool
				  var username string //only for bmc
				  var password string //only for bmc
				  var provider string //only accepted IPMI
				  var virtual bool
				  var tag string
				  var attached_to string
				  var mode string // with validations
				  var attached_devices []string
				  var bond_options string
				  var compute_attributes []string{}
				}
  compute_attributes		{
				  var cpus string
				  var start string
				  var cluster string
				  var memory_mb string
				  var guest_id string
				}
}

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
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment-id": &schema.Schema{
				Type:     schema.TypeString,
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
			"puppetclass_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"operatingsystem_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"medium_id": &schema.Schema{
				Type:     schema.TypeString, 
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
				Type:     schema.TypeInt, // if nil it assumes bare-metal build ,Add some validation logic later when you have time.
				Optional: true,
			},
			"root_pass": &schema.Schema{
				Type:     schema.TypeString, // required if host is managed and not inherited from hostgroup or default password.
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
			"progress_report_id": &schema.Schema{
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
			"storage_attributes": &schema.Schema{
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

// Setup a function to make api calls
func httpClient(rType string, h *host, u *userAccess. meta interface{}) error {
  //setup local vars
  r := strings.ToLower(rType)
  lUserAccess := userAccess
  lHost := h

  //select and run request type.
  if (r == "get"){
    http.Get()
  } else if (r == "post") {
    http.Post()
  } else if (r == "put") {
    http.Put()
  } else if (r == "delete") {
    http.Delete()
  }
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("name").(string))
        h := host{
          name: d.Get("name").(string),
        }

/* populate h struct instance */
        if v, ok := d.GetOk("environment-id"); ok {
          h.environment_id = v.(string)
        }
        if v,ok := d.GetOk("ip"); ok{
          h.ip = v.(string)
        }
        if v,ok := d.GetOk("mac"); ok{
          h.mac = v.(string)
        }
        if v,ok := d.GetOk("architecture_id"); ok{
          h.architecture_id = v.(string)
        }
        if v,ok := d.GetOk("domain_id"); ok{
          h.domain_id = v.(int)
        }
        if v,ok := d.GetOk("realm_id"); ok{
          h.realm_id = v.(int)
        }
        if v,ok := d.GetOk("puppet_proxy_id"); ok{
          h.puppet_proxy_id = v.(int)
        }
        if v,ok := d.GetOk("puppet_class_ids"); ok{
          h.puppet_class_ids = v.([]int)
        }
        if v,ok := d.GetOk("operatingsystem_id"); ok{
          h.operatingsystem_id = v.(string)
        }
        if v,ok := d.GetOk("medium_id"); ok{
          h.medium_id = v.(string)
        }
        if v,ok := d.GetOk("ptable_id"); ok{
          h.ptable_id = v.(int)
        }
        if v,ok := d.GetOk("subnet_id"); ok{
          h.subnet_id = v.(int)
        }
        if v,ok := d.GetOk("computer_resource_id"); ok{
          h.compute_resource_id = v.(int)
        }
        if v,ok := d.GetOk("root_pass"); ok{
          h.root_pass = v.(string)
        }
        if v,ok := d.GetOk("model_id"); ok{
          h.model_id = v.(int)
        }
        if v,ok := d.GetOk("hostgroup_id"); ok{
          h.hostgroup_id = v.(int)
        }
        if v,ok := d.GetOk("owner_id"); ok{
          h.owner_id = v.(int)
        }
        if v,ok := d.GetOk("owner_type"); ok{
          h.owner_type = v.(int)
        }
        if v,ok := d.GetOk("puppet_ca_proxy_id"); ok{
          h.puppet_ca_proxy_id = v.(int)
        }
        if v,ok := d.GetOk("image_id"); ok{
          h.image_id = v.(string)
        }
        if v,ok := d.GetOk("host_parameters_attributes"); ok{
          h.host_parameters_attributes = v.([]string)
        }
        if v,ok := d.GetOk("build"); ok{
          h.build = v.(bool)
        }
        if v,ok := d.GetOk("enabled"); ok{
          h.enabled = v.(bool)
        }
        if v,ok := d.GetOk("provision_method"); ok{
          h.provision_method = v.(string)
        }
        if v,ok := d.GetOk("managed"); ok{
          h.managed = v.(bool)
        }
        if v,ok := d.GetOk("progress_report_id"); ok{
          h.progess_report_id = v.(string)
        }
        if v,ok := d.GetOk("comment"); ok{
          h.comment = v.(string)
        }
        if v,ok := d.GetOk("capabilities"); ok{
          h.capabilities = v.(string)
        }
        if v,ok := d.GetOk("compute_profile_id"); ok{
          h.compute_profile_id = v.(int)
        }

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
