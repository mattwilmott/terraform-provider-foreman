package main

import (
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
  "strings"
	"reflect"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
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
	volumes_attributes		{
					var name		  string
					var size_gb	  int
					var _delete	  string
					var datastore	string
				}
}

type userAccess struct {
	var username	string
	var password	string
	var url				string
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
			"interfaces_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"volumes_attributes": &schema.Schema{
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
func httpClient(rType string, d *data, u *userAccess, meta interface{}) error {
  //setup local vars
  r := strings.ToUpper(rType)
  lUserAccess := u
  lData := d

  //build and make request
	client := &http.Client{}
	req, err := http.NewRequest(r,lUserAccess.url,lData)
	//set basic auth if necessary
	if (u.username != nil){
	req.SetBasicAuth(lUserAccess.username,lUserAccess.password)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

  resp, err := client.Do(req)

	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
}


func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("name").(string))
        h := host{
          name: d.Get("name").(string),
        }

				u := userAccess{}

/* populate u struct instance */
				if v, ok := d.GetOk("username"); ok {
					u.username = v.(string)
				}
				if v, ok := d.GetOk("password"); ok {
					u.password = v.(string)
				}
				if v, ok := d.GetOk("url"); ok {
					u.url	= v.(string)
				}

/* build subtree level stuff first */
/* build compute_attributes now */
				if v, ok := d.GetOk("compute_attributes.cpus"); ok {
					h.compute_attributes.cpus = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.start"); ok {
					h.compute_attributes.cluster = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.memory_mb"); ok {
					h.compute_attributes.memory_mb = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.guest_id"); ok {
					h.compute_attributes.guest_id = v.(string)
				}
/* build volumes_attributes now */
				if v, ok := d.GetOk("volumes_attributes.name"); ok {
					h.volumes_attributes.name = v.(string)
				}
				if v, ok := d.GetOk("volumes_attributes.size_gb"); ok {
					h.volumes_attributes.size_gb = v.(string)
				}
				if v, ok := d.GetOk("volumes_attributes._delete"); ok {
					h.volumes_attributes._delete = v.(string)
				}
				if v, ok := d.GetOk("volumes_attributes.datastore"); ok {
					h.volumes_attributes.datastore = v.(string)
				}
/* build interfaces_attributes now */
				if v, ok := d.GetOk("interfaces_attributes.mac"); ok {
					h.interfaces_attributes.mac = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.ip"); ok {
					h.interfaces_attributes.ip = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.type"); ok {
					h.interfaces_attributes.type = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.name"); ok {
					h.interfaces_attributes.name = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.subnet_id"); ok {
					h.interfaces_attributes.subnet_id = v.(int)
				}
				if v, ok := d.GetOk("interfaces_attributes.domain_id"); ok {
					h.interfaces_attributes.domain_id = v.(int)
				}
				if v, ok := d.GetOk("interfaces_attributes.identifier"); ok {
					h.interfaces_attributes.identifier = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.managed"); ok {
					h.interfaces_attributes.managed = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.primary"); ok {
					h.interfaces_attributes.primary = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.provision"); ok {
					h.interfaces_attributes.provision = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.username"); ok {
					h.interfaces_attributes.username = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.password"); ok {
					h.interfaces_attributes.password = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.provider"); ok {
					h.interfaces_attributes.provider = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.virtual"); {
					h.interfaces_attributes.virtual = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.tag"); ok {
					h.interfaces_attributes.tag = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.attached_to"); ok {
					h.interfaces_attributes.attached_to = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.mode"); ok {
					h.interfaces_attributes.mode = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.attached_devices"); ok {
					h.interfaces_attributes.attached_devices = v.([]string)
				}
				if v, ok := d.GetOk("interfaces_attributes.bond_options"); ok {
					h.interfaces_attributes.attached_devices = v.(string)
				}
/* pupulate host_parameters_attributes now */
				if v, ok := d.GetOk("host_parameters_attributes.roles"); ok {
					h.host_parameters_attributes.roles = v.(string)
				}
				if v, ok := d.GetOk("host_parameters_attributes.puppet"); ok {
					h.host_parameters_attributes.puppet = v.(string)
				}
				if v, ok := d.GetOk("host_parameters_attributes.chef"); ok {
					h.host_parameters_attributes.chef = v.(string)
				}
				if v, ok := d.GetOk("host_parameters_attributes.JIRA_Ticket"); ok {
					h.host_parameters_attributes.JIRA_Ticket = v.(string)
				}
/* populate h struct instance for regular level data */
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
  jData, err := json.Marshal(h)
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
