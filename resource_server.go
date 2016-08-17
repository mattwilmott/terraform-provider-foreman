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
			"server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"location-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization-id": &schema.Schema{
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
			"architecture-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"realm-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"puppet-proxy-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"puppet-class-ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"operatingsystem-id": &schema.Schema{
				Type:     schema.TypeString, //Why isnt this an Int? API doco may be incorrect
				Optional: true,
			},
			"medium-id": &schema.Schema{
				Type:     schema.TypeString, //Why isnt this an Int as well? wtf
				Optional: true,
			},
			"partition-table-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"compute-resource-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"root-pass": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"model-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostgroup-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"puppet-ca-proxy-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"host-parameters-attributes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"build": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provision-method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provision-report-id": &schema.Schema{
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
			"compute-profile-id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"interface-attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"storage-attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"compute-attributes": &schema.Schema{
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
        if v,ok := d.GetOk("architecture-id"); ok{
          h.architecture_id = v.(string)
        }
        if v,ok := d.GetOk("domain-id"); ok{
          h.domain_id = v.(int)
        }
        if v,ok := d.GetOk("realm-id"); ok{
          h.realm_id = v.(int)
        }
        if v,ok := d.GetOk("puppet-proxy-id"); ok{
          h.puppet_proxy_id = v.(int)
        }
        if v,ok := d.GetOk("puppet-class-ids"); ok{
          h.puppet_class_ids = v.([]int)
        }
        if v,ok := d.GetOk("operatingsystem-id"); ok{
          h.operatingsystem_id = v.(string)
        }
        if v,ok := d.GetOk("medium-id"); ok{
          h.medium_id = v.(string)
        }
        if v,ok := d.GetOk("ptable-id"); ok{
          h.ptable_id = v.(int)
        }
        if v,ok := d.GetOk("subnet-id"); ok{
          h.subnet_id = v.(int)
        }
        if v,ok := d.GetOk("computer-resource-id"); ok{
          h.compute_resource_id = v.(int)
        }
        if v,ok := d.GetOk("root-pass"); ok{
          h.root_pass = v.(string)
        }
        if v,ok := d.GetOk("model-id"); ok{
          h.model_id = v.(int)
        }
        if v,ok := d.GetOk("hostgroup-id"); ok{
          h.hostgroup_id = v.(int)
        }
        if v,ok := d.GetOk("owner-id"); ok{
          h.owner_id = v.(int)
        }
        if v,ok := d.GetOk("owner-type"); ok{
          h.owner_type = v.(int)
        }
        if v,ok := d.GetOk("puppet-ca-proxy-id"); ok{
          h.puppet_ca_proxy_id = v.(int)
        }
        if v,ok := d.GetOk("image-id"); ok{
          h.image_id = v.(string)
        }
        if v,ok := d.GetOk("host-parameters-attributes"); ok{
          h.host_parameters_attributes = v.([]string)
        }
        if v,ok := d.GetOk("build"); ok{
          h.build = v.(bool)
        }
        if v,ok := d.GetOk("enabled"); ok{
          h.enabled = v.(bool)
        }
        if v,ok := d.GetOk("provision-method"); ok{
          h.provision_method = v.(string)
        }
        if v,ok := d.GetOk("managed"); ok{
          h.managed = v.(bool)
        }
        if v,ok := d.GetOk("progress-report-id"); ok{
          h.progess_report_id = v.(string)
        }
        if v,ok := d.GetOk("comment"); ok{
          h.comment = v.(string)
        }
        if v,ok := d.GetOk("capabilities"); ok{
          h.capabilities = v.(string)
        }
        if v,ok := d.GetOk("compute-profile-id"); ok{
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
