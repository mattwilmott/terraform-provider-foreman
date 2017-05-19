package main

import (
	//"bytes"
	//"encoding/json"
	//"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mattwilmott/go-foreman"
	//"io/ioutil"
	//"net/http"
	//"strings"
	"log"
)

//reqHost is required to wrap the host for the foreman-api
type reqHost struct {
	//Lhost host `json:"host,omitempty"`
}

//Used for access authentication to foreman
type userAccess struct {
	username string
	password string
	url      string
}

type fRespDomain struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

//This sets up the schema, the interface between the tf file and the plugin
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
			"environment_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"location_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
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
				Type:     schema.TypeInt, //Why isnt this an Int? API doco may be incorrect
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
			"owner_type": &schema.Schema{
				Type:     schema.TypeString,
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
			/* See further down, this was the naive impl
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
			// Renamed progress_report_id in later versions
			//"provision_report_id": &schema.Schema{
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
			// List of network interface definitions
			"interfaces_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"subnet_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
						"domain_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
						"identifier": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"managed": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"primary": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"provision": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"provider": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"virtual": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"tag": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"attached_to": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"attached_devices": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							ForceNew: false,
						},
						"bond_options": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"compute_attributes": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"volumes_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"size_gb": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: false,
						},
						"_delete": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: false,
						},
						"datastore": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"compute_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"host_parameters_attributes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"puppet": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"chef": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
		},
	}
}

/* Wrong place for this shit
func getDomain(h *host, u *userAccess) string {
	dStruct := new(fRespDomain)
	//resp, err := httpClient("GET", h, u, "domains", false, "")
	client := meta.(*ForemanClient)
	resp, err := client.getDomains()
	if resp != nil {
		fResp := log.Printf("The server responded with: %v", resp)
		print(fResp)
		if strings.Contains(string(resp), "error") {
			err = errors.New(string(resp))
		}
	}

	if err != nil {
		panic(err)
	}

	unerr := json.Unmarshal(resp, &dStruct)
	if unerr != nil {
		panic(unerr)
	}
	return dStruct.Name
}
*/

// Private func to build out the Host Struct based on either the tf definition
func buildHostStruct(d *schema.ResourceData, meta interface{}) foreman.Host {
	h := foreman.Host{}
	if v, ok := d.GetOk("name"); ok {
		h.Name = v.(string)
	}
	caprefix := fmt.Sprintf("compute_attributes")
	if v, ok := d.GetOk(caprefix + ".cpus"); ok {
		h.Compute_attributes.Cpus = v.(string)
	}
	if v, ok := d.GetOk(caprefix + ".start"); ok {
		h.Compute_attributes.Start = v.(string)
	}
	if v, ok := d.GetOk(caprefix + ".cluster"); ok {
		h.Compute_attributes.Cluster = v.(string)
	}
	if v, ok := d.GetOk(caprefix + ".memory_mb"); ok {
		h.Compute_attributes.Memory_mb = v.(string)
	}
	if v, ok := d.GetOk(caprefix + ".guest_id"); ok {
		h.Compute_attributes.Guest_id = v.(string)
	}
	// build volumes_attributes now, this needs to be mapped so I build structs and then add them to the mapped value
	vaCount := d.Get("volumes_attributes.#").(int)
	if vaCount > 0 {
		h.Compute_attributes.Volumes_attributes_map = make(map[string]foreman.Volumes_attributes)
		for i := 0; i < vaCount; i++ {
			// setup iterator market and instantiate local struct to append to maps
			iStr := fmt.Sprintf("%d", i)
			lStruct := foreman.Volumes_attributes{}
			//setup lStruct values
			vaprefix := fmt.Sprintf("volumes_attributes.%d", i)
			if v, ok := d.GetOk(vaprefix + ".name"); ok {
				lStruct.Name = v.(string)
			}
			if v, ok := d.GetOk(vaprefix + ".size_gb"); ok {
				lStruct.Size_gb = v.(int)
			}
			//if v, ok := d.GetOk(vaprefix + "._delete"); ok {
			//	lStruct.Delete = v.(string)
			//}
			if v, ok := d.GetOk(vaprefix + ".datastore"); ok {
				lStruct.Datastore = v.(string)
			}
			//add in lStruct to the main host struct
			h.Compute_attributes.Volumes_attributes_map[iStr] = lStruct
		}
	}
	// build interfaces_attributes now
	iaCount := d.Get("interfaces_attributes.#").(int)
	if iaCount > 0 {
		for i := 0; i < iaCount; i++ {
			h.Interfaces_attributes_array = append(h.Interfaces_attributes_array, foreman.Interfaces_attributes{})
			prefix := fmt.Sprintf("interfaces_attributes.%d", i)
			if v, ok := d.GetOk(prefix + ".primary"); ok {
				h.Interfaces_attributes_array[i].Primary = v.(bool)
			}
			//Adding some logic to auto populate primary nic because of API deferral
			if v, ok := d.GetOk(prefix + ".mac"); ok {
				h.Interfaces_attributes_array[i].Mac = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".ip"); ok {
				h.Interfaces_attributes_array[i].Ip = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".type"); ok {
				h.Interfaces_attributes_array[i].Type = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".name"); ok {
				h.Interfaces_attributes_array[i].Name = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".subnet_id"); ok {
				h.Interfaces_attributes_array[i].Subnet_id = v.(int)
			}
			if v, ok := d.GetOk(prefix + ".domain_id"); ok {
				h.Interfaces_attributes_array[i].Domain_id = v.(int)
			}
			if v, ok := d.GetOk(prefix + ".identifier"); ok {
				h.Interfaces_attributes_array[i].Identifier = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".managed"); ok {
				h.Interfaces_attributes_array[i].Managed = v.(bool)
			}
			if v, ok := d.GetOk(prefix + ".provision"); ok {
				h.Interfaces_attributes_array[i].Provision = v.(bool)
			}
			if v, ok := d.GetOk(prefix + ".username"); ok {
				h.Interfaces_attributes_array[i].Username = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".password"); ok {
				h.Interfaces_attributes_array[i].Password = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".provider"); ok {
				h.Interfaces_attributes_array[i].Provider = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".virtual"); ok {
				h.Interfaces_attributes_array[i].Virtual = v.(bool)
			}
			if v, ok := d.GetOk(prefix + ".tag"); ok {
				h.Interfaces_attributes_array[i].Tag = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".attached_to"); ok {
				h.Interfaces_attributes_array[i].Attached_to = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".mode"); ok {
				h.Interfaces_attributes_array[i].Mode = v.(string)
			}
			if v, ok := d.GetOk(prefix + ".attached_devices"); ok {
				h.Interfaces_attributes_array[i].Attached_devices = v.([]string)
			}
			if v, ok := d.GetOk(prefix + ".bond_options"); ok {
				h.Interfaces_attributes_array[i].Bond_options = v.(string)
			}
			ifcaprefix := fmt.Sprintf("%s.compute_attributes", prefix)
			if v, ok := d.GetOk(ifcaprefix + ".network"); ok {
				h.Interfaces_attributes_array[i].Compute_attributes.Network = v.(string)
			}
			if v, ok := d.GetOk(ifcaprefix + ".type"); ok {
				h.Interfaces_attributes_array[i].Compute_attributes.Type = v.(string)
			}
		}
	}

	// populate host_parameters_attributes now
	hpaCount := d.Get("host_parameters_attributes.#").(int)
	if hpaCount > 0 {
		h.Host_parameters_attributes_map = make(map[string]foreman.Params_archetype)
		for i := 0; i < hpaCount; i++ {
			intCnt := 0
			prefix := fmt.Sprintf("host_parameters_attributes.%d", i)
			if v, ok := d.GetOk(prefix + ".roles"); ok {
				roleStruct := foreman.Params_archetype{}
				iStr := fmt.Sprintf("%d", intCnt)
				roleStruct.Name = "roles"
				roleStruct.Value = v.(string)
				h.Host_parameters_attributes_map[iStr] = roleStruct
				intCnt++
			}
			if v, ok := d.GetOk(prefix + ".puppet"); ok {
				pupStruct := foreman.Params_archetype{}
				iStr := fmt.Sprintf("%d", intCnt)
				pupStruct.Name = "puppet"
				pupStruct.Value = v.(string)
				h.Host_parameters_attributes_map[iStr] = pupStruct
				intCnt++
			}
			if v, ok := d.GetOk(prefix + ".chef"); ok {
				chefStruct := foreman.Params_archetype{}
				iStr := fmt.Sprintf("%d", intCnt)
				chefStruct.Name = "chef"
				chefStruct.Value = v.(string)
				h.Host_parameters_attributes_map[iStr] = chefStruct
				intCnt++
			}
		}
	}

	// populate h struct instance for regular level data
	if v, ok := d.GetOk("environment_id"); ok {
		h.Environment_id = v.(int)
	}
	if v, ok := d.GetOk("organization_id"); ok {
		h.Organization_id = v.(int)
	}
	if v, ok := d.GetOk("location_id"); ok {
		h.Location_id = v.(int)
	}
	if v, ok := d.GetOk("ip"); ok {
		h.Ip = v.(string)
	}
	if v, ok := d.GetOk("mac"); ok {
		h.Mac = v.(string)
	}
	if v, ok := d.GetOk("architecture_id"); ok {
		h.Architecture_id = v.(int)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		h.Domain_id = v.(int)
	}
	if v, ok := d.GetOk("realm_id"); ok {
		h.Realm_id = v.(int)
	}
	if v, ok := d.GetOk("puppet_proxy_id"); ok {
		h.Puppet_proxy_id = v.(int)
	}
	if v, ok := d.GetOk("puppet_class_ids"); ok {
		h.Puppetclass_ids = v.([]int)
	}
	if v, ok := d.GetOk("operatingsystem_id"); ok {
		h.Operatingsystem_id = v.(int)
	}
	if v, ok := d.GetOk("medium_id"); ok {
		h.Medium_id = v.(int)
	}
	if v, ok := d.GetOk("ptable_id"); ok {
		h.Ptable_id = v.(int)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		h.Subnet_id = v.(int)
	}
	if v, ok := d.GetOk("compute_resource_id"); ok {
		h.Compute_resource_id = v.(int)
	}
	if v, ok := d.GetOk("root_pass"); ok {
		h.Root_pass = v.(string)
	}
	if v, ok := d.GetOk("model_id"); ok {
		h.Model_id = v.(int)
	}
	if v, ok := d.GetOk("hostgroup_id"); ok {
		h.Hostgroup_id = v.(int)
	}
	if v, ok := d.GetOk("owner_id"); ok {
		h.Owner_id = v.(int)
	}
	if v, ok := d.GetOk("owner_type"); ok {
		if v.(string) == "User" || v.(string) == "Usergroup" {
			h.Owner_type = v.(string)
		}
	}
	if v, ok := d.GetOk("puppet_ca_proxy_id"); ok {
		h.Puppet_ca_proxy_id = v.(int)
	}
	if v, ok := d.GetOk("image_id"); ok {
		h.Image_id = v.(int)
	}
	if v, ok := d.GetOk("build"); ok {
		h.Build = v.(bool)
	}
	if v, ok := d.GetOk("enabled"); ok {
		h.Enabled = v.(bool)
	}
	if v, ok := d.GetOk("provision_method"); ok {
		h.Provision_method = v.(string)
	}
	if v, ok := d.GetOk("managed"); ok {
		h.Managed = v.(bool)
	}
	if v, ok := d.GetOk("progress_report_id"); ok {
		h.Progress_report_id = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok {
		h.Comment = v.(string)
	}
	if v, ok := d.GetOk("capabilities"); ok {
		h.Capabilities = v.([]interface{})
	}
	if v, ok := d.GetOk("compute_profile_id"); ok {
		h.Compute_profile_id = v.(int)
	}
	return h
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	h := buildHostStruct(d, meta)
	client := meta.(*ForemanClient).foremanconn

	// check debug flag
	debug := false
	if v, ok := d.GetOk("Debug"); ok {
		debug = v.(bool)
	}

	//resp, err := httpClient("POST", &h, &client, "hosts", debug, "")
	respHost, err := client.CreateHost(&h)
	if respHost != nil {
		log.Printf("Successfully created %s\n", respHost.Name)
		if debug {
			log.Printf("The server responded with: %+v\n", respHost)
		}
		//d.SetId(d.Get("name").(string))
		return nil
	} else if err != nil {
		log.Printf("ERROR: Foreman failed to create host - %s\n", h.Name)
		log.Printf("ERROR: %s", err)
		return err
	} else {
		log.Printf("Successfully sent created request for server %s but there was no data returned!\n", h.Name)
		panic("Unknown state whilst creating server resource details!")
	}
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	h := buildHostStruct(d, meta)
	client := meta.(*ForemanClient).foremanconn
	debug := true

	respHost, err := client.GetHost(&h)
	if respHost != nil {
		if debug {
			log.Printf("The server responded with: %+v", respHost)
		}
		d.SetId(d.Get("name").(string))
		d.Set("ip", h.Ip)
		d.Set("comment", h.Comment)

		return nil
	} else if err != nil {
		log.Printf("Failed to read server: %s\n", h.Name)
		log.Printf("Status Code: %s\n", err)
		if err.Error() == "HTTP Error 404 Not Found" {
			// Server not found rather than API error
			return nil
		}
		return err
	} else {
		log.Printf("Successfully read server %s but there was no data returned!\n", h.Name)
		panic("Unknown state whilst fetching server resource details!")
	}
}

func resourceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	h := buildHostStruct(d, meta)
	client := meta.(*ForemanClient).foremanconn
	//hChanges := new(host)
	//dom := getDomain(&h, &u)
	//fqdn := log.Printf("%s.%s", h.Name, dom)

	//if (fqdn != "") && (fqdn != dom) {
	//resp, err := httpClient("PUT", &h, &u, "hosts", false, fqdn)
	respHost, err := client.UpdateHost(&h)
	if respHost != nil {
		log.Printf("The server responded with: %v", respHost.Name)
	} else if err != nil {
		return err
	} else {
		log.Printf("Successfully updated server %s but there was no data returned!\n", h.Name)
		panic("Unknown state whilst updating server resource details!")
	}
	//}
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	/* commenting out till this can be properly tested
	h := buildHostStruct(d,m)
	dom := getDomain(&h,&u)
	fqdn := log.Printf("%s.%s",h.Name,dom)
	if (fqdn != "") && (fqdn != dom) {
	 resp, err := httpClient("DELETE", &h, &u, "hosts", false,fqdn)
	 if resp != nil {
		 fResp := log.Printf("The server responded with: %v",resp)
		 print(fResp)
	 	 if strings.Contains(string(resp),"error"){
	 		 err = errors.New(string(resp))
	 	 }
	 }
	 if err != nil {
		return err
	 }
	}
	*/
	return nil
}
