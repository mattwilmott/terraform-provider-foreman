package main

import (
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
  "strings"
	//"reflect"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	//"os"
	"bytes"
	"strconv"
)
type host_parameters_attributes	struct {
  Roles 			string	`json:"roles,omitempty"`
	Puppet 			string	`json:"puppet,omitempty"`
	Chef 				string	`json:"chef,omitempty"`
	JIRA_Ticket string	`json:",omitempty"`
}
type interfaces_attributes	struct	{
	Mac 								string	`json:"mac,omitempty"`
	Ip 									string	`json:"ip,omitempty"`
	//type 								string
	Name 								string	`json:"name,omitempty"`
	Subnet_id 					int			`json:"subnet_id,omitempty"`
	Domain_id 					int			`json:"domain_id,omitempty"`
	Identifier 					string	`json:"identifier,omitempty"`
	Managed 						bool		`json:"managed,omitempty"`
	Primary 						bool		`json:"primary,omitempty"`
	Provision 					bool		`json:"provision,omitempty"`
	Username 						string	`json:"username,omitempty"`//only for bmc
	Password 						string	`json:"password,omitempty"` //only for bmc
	Provider 						string	`json:"provider,omitempty"` //only accepted IPMI
	Virtual 						bool		`json:"virtual,omitempty"`
	Tag 								string	`json:"tag,omitempty"`
	Attached_to 				string	`json:"attached_to,omitempty"`
	Mode 								string	`json:"mode,omitempty"` // with validations
	Attached_devices 		[]string
	Bond_options 				string	`json:"bond_options,omitempty"`
	compute_attributes
}
type compute_attributes	struct {
	Cpus 			string	`json:"cpus,omitempty"`
	//start 		string
	Cluster 	string	`json:"cluster,omitempty"`
	Memory_mb string	`json:"memory_mb,omitempty"`
	Guest_id 	string	`json:"guest_id,omitempty"`
}
type	volumes_attributes struct {
	Name		  string	`json:"name,omitempty"`
	Size_gb	  int			`json:"size_gb,omitempty"`
	_delete	  string	`json:",omitempty"`
	Datastore	string	`json:"datastore,omitempty"`
}

type host struct {
  Name									string	`json:"name,omitempty"`
  Environment_id				string	`json:"environment_id,omitempty"`
  Ip										string	`json:"ip,omitempty"`
  Mac										string	`json:"mac,omitempty"`
  Architecture_id				int			`json:"architecture_id,omitempty"`
  Domain_id							int			`json:"domain_id,omitempty"`
  Realm_id							int			`json:"realm_id,omitempty"`
  Puppet_proxy_id				int			`json:"puppet_proxy_id,omitempty"`
  Puppetclass_ids				[]int		`json:"puppetclass_ids,omitempty"`
  Operatingsystem_id		string	`json:"operatingsystem_id,omitempty"`
  Medium_id							string	`json:"medium_id,omitempty"`
  Ptable_id							int			`json:"ptable_id,omitempty"`
  Subnet_id							int			`json:"subnet_id,omitempty"`
  Compute_resource_id		int			`json:"compute_resource_id,omitempty"`
  Root_pass							string	`json:"root_pass,omitempty"`
  Model_id							int			`json:"model_id,omitempty"`
  Hostgroup_id					int			`json:"hostgroup_id,omitempty"`
  Owner_id							int			`json:"owner_id,omitempty"`
  Owner_type						string	`json:"owner_type,omitempty"` // must be either User or Usergroup
  Puppet_ca_proxy_id		int			`json:"puppet_ca_proxy_id,omitempty"`
  Image_id							int			`json:"image_id,omitempty"`
  Build									bool		`json:"build,omitempty"`
  Enabled								bool		`json:"enabled,omitempty"`
  Provision_method			string	`json:"provision_method,omitempty"`
  Managed								bool		`json:"managed,omitempty"`
  Progress_report_id		string	`json:"progress_report_id,omitempty"`
  Comment								string	`json:"comment,omitempty"`
  Capabilities					string	`json:"capabilities,omitempty"`
  Compute_profile_id		int			`json:"compute_profile_id,omitempty"`
	host_parameters_attributes
  interfaces_attributes
  compute_attributes
	volumes_attributes
}

type userAccess struct {
	username	string
	password	string
	url				string
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"debug": &schema.Schema{
				Type:			schema.TypeBool,
				Optional:	true,
			},
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
			"environment_id": &schema.Schema{
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
			/*
			"host_parameters_attributes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			*/
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
			"host_parameters_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

// Setup a function to make api calls
func httpClient(rType string, d *host, u *userAccess, debug bool) error {
	println("JPB - Made it to httpClient")
  //setup local vars
  r := strings.ToUpper(rType)
  lUserAccess := u


  b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(d)

	//panic(b)
  //build and make request
	client := &http.Client{}
	req, err := http.NewRequest(r,lUserAccess.url,b)

	if err != nil {
		panic(err)
	}
	println("JPB - Setup request and client successfully")
	//set basic auth if necessary
	if u.username != "" {
	req.SetBasicAuth(lUserAccess.username,lUserAccess.password)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	println("JPB - Setup basic auth headers")
   //enable debugging data
	if debug {
		panic(req)
	}
  resp, err := client.Do(req)
	println("JPB - Request made to server")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	println("JPB - Reading content")
	if content != nil {
		fmt.Println("%v",content)
	}

	return nil
}


func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	println("JPB - Made it to create method")
	d.SetId(d.Get("name").(string))
        h := new(host)

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
				println("JPB - Built u struct instance")

/* build subtree level stuff first */
/* build compute_attributes now */
				if v, ok := d.GetOk("name"); ok {
					h.Name = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.cpus"); ok {
					h.compute_attributes.Cpus = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.start"); ok {
					h.compute_attributes.Cluster = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.memory_mb"); ok {
					h.compute_attributes.Memory_mb = v.(string)
				}
				if v, ok := d.GetOk("compute_attributes.guest_id"); ok {
					h.compute_attributes.Guest_id = v.(string)
				}
				println("JPB - Built compute_attrs struct instance")
/* build volumes_attributes now */
				if v, ok := d.GetOk("volumes_attributes.name"); ok {
					h.volumes_attributes.Name = v.(string)
				}
				println("JPB - Added volumes_attributes.name ")
				println("JPB - About to add volumes_attributes.size_gb")
				if v, ok := d.GetOk("volumes_attributes.size_gb"); ok {
					i, _ := strconv.Atoi(v.(string))
					h.volumes_attributes.Size_gb = i
				}
				println("JPB - Added volumes_attributes.size_gb")
				if v, ok := d.GetOk("volumes_attributes._delete"); ok {
					h.volumes_attributes._delete = v.(string)
				}
				println("JPB - Added volumes_attributes._delete")
				if v, ok := d.GetOk("volumes_attributes.datastore"); ok {
					h.volumes_attributes.Datastore = v.(string)
				}
				println("JPB - Added volumes_attributes.datastore")
				println("JPB - Built u volumes_attributes instance")
/* build interfaces_attributes now */
				if v, ok := d.GetOk("interfaces_attributes.mac"); ok {
					h.interfaces_attributes.Mac = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.ip"); ok {
					h.interfaces_attributes.Ip = v.(string)
				}
				/* removing the following bit since type is a keyword
				if v, ok := d.GetOk("interfaces_attributes.type"); ok {
					h.interfaces_attributes.type = v.(string)
				}
				*/
				if v, ok := d.GetOk("interfaces_attributes.name"); ok {
					h.interfaces_attributes.Name = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.subnet_id"); ok {
					h.interfaces_attributes.Subnet_id = v.(int)
				}
				if v, ok := d.GetOk("interfaces_attributes.domain_id"); ok {
					h.interfaces_attributes.Domain_id = v.(int)
				}
				if v, ok := d.GetOk("interfaces_attributes.identifier"); ok {
					h.interfaces_attributes.Identifier = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.managed"); ok {
					h.interfaces_attributes.Managed = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.primary"); ok {
					h.interfaces_attributes.Primary = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.provision"); ok {
					h.interfaces_attributes.Provision = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.username"); ok {
					h.interfaces_attributes.Username = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.password"); ok {
					h.interfaces_attributes.Password = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.provider"); ok {
					h.interfaces_attributes.Provider = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.virtual"); ok{
					h.interfaces_attributes.Virtual = v.(bool)
				}
				if v, ok := d.GetOk("interfaces_attributes.tag"); ok {
					h.interfaces_attributes.Tag = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.attached_to"); ok {
					h.interfaces_attributes.Attached_to = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.mode"); ok {
					h.interfaces_attributes.Mode = v.(string)
				}
				if v, ok := d.GetOk("interfaces_attributes.attached_devices"); ok {
					h.interfaces_attributes.Attached_devices = v.([]string)
				}
				if v, ok := d.GetOk("interfaces_attributes.bond_options"); ok {
					h.interfaces_attributes.Bond_options = v.(string)
				}
				println("JPB - Built interfaces struct instance")
/* pupulate host_parameters_attributes now */
				if v, ok := d.GetOk("host_parameters_attributes.roles"); ok {
					h.host_parameters_attributes.Roles = v.(string)
				}
				if v, ok := d.GetOk("host_parameters_attributes.puppet"); ok {
					h.host_parameters_attributes.Puppet = v.(string)
				}
				if v, ok := d.GetOk("host_parameters_attributes.chef"); ok {
					h.host_parameters_attributes.Chef = v.(string)
				}
				if v, ok := d.GetOk("host_parameters_attributes.JIRA_Ticket"); ok {
					h.host_parameters_attributes.JIRA_Ticket = v.(string)
				}
				println("JPB - Built host_parameters_attributes struct instance")
/* populate h struct instance for regular level data */
        if v, ok := d.GetOk("environment-id"); ok {
          h.Environment_id = v.(string)
        }
        if v,ok := d.GetOk("ip"); ok{
          h.Ip = v.(string)
        }
        if v,ok := d.GetOk("mac"); ok{
          h.Mac = v.(string)
        }
        if v,ok := d.GetOk("architecture_id"); ok{
          h.Architecture_id = v.(int)
        }
        if v,ok := d.GetOk("domain_id"); ok{
          h.Domain_id = v.(int)
        }
        if v,ok := d.GetOk("realm_id"); ok{
          h.Realm_id = v.(int)
        }
        if v,ok := d.GetOk("puppet_proxy_id"); ok{
          h.Puppet_proxy_id = v.(int)
        }
        if v,ok := d.GetOk("puppet_class_ids"); ok{
          h.Puppetclass_ids = v.([]int)
        }
        if v,ok := d.GetOk("operatingsystem_id"); ok{
          h.Operatingsystem_id = v.(string)
        }
        if v,ok := d.GetOk("medium_id"); ok{
          h.Medium_id = v.(string)
        }
        if v,ok := d.GetOk("ptable_id"); ok{
          h.Ptable_id = v.(int)
        }
        if v,ok := d.GetOk("subnet_id"); ok{
          h.Subnet_id = v.(int)
        }
        if v,ok := d.GetOk("computer_resource_id"); ok{
          h.Compute_resource_id = v.(int)
        }
        if v,ok := d.GetOk("root_pass"); ok{
          h.Root_pass = v.(string)
        }
        if v,ok := d.GetOk("model_id"); ok{
          h.Model_id = v.(int)
        }
        if v,ok := d.GetOk("hostgroup_id"); ok{
          h.Hostgroup_id = v.(int)
        }
        if v,ok := d.GetOk("owner_id"); ok{
          h.Owner_id = v.(int)
        }
        if v,ok := d.GetOk("owner_type"); ok{
					if v.(string) == "User" || v.(string) == "Usergroup" {
					  h.Owner_type = v.(string)
					}
        }
        if v,ok := d.GetOk("puppet_ca_proxy_id"); ok{
          h.Puppet_ca_proxy_id = v.(int)
        }
        if v,ok := d.GetOk("image_id"); ok{
          h.Image_id = v.(int)
        }
        if v,ok := d.GetOk("build"); ok{
          h.Build = v.(bool)
        }
        if v,ok := d.GetOk("enabled"); ok{
          h.Enabled = v.(bool)
        }
        if v,ok := d.GetOk("provision_method"); ok{
          h.Provision_method = v.(string)
        }
        if v,ok := d.GetOk("managed"); ok{
          h.Managed = v.(bool)
        }
        if v,ok := d.GetOk("progress_report_id"); ok{
          h.Progress_report_id = v.(string)
        }
        if v,ok := d.GetOk("comment"); ok{
          h.Comment = v.(string)
        }
        if v,ok := d.GetOk("capabilities"); ok{
          h.Capabilities = v.(string)
        }
        if v,ok := d.GetOk("compute_profile_id"); ok{
          h.Compute_profile_id = v.(int)
        }
				println("JPB - Built h struct instance")
	/* check debug flag */
	println("JPB -  setting up debug flag")
	debug := false
	if v, ok := d.GetOk("debug"); ok {
		debug = v.(bool)
	}
  println("JPB - Debug complete calling httpClient")
	httpClient("POST", &h, &u, debug)
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
