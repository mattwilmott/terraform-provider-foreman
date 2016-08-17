package main

import (
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
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
				  var compute_attributes {}
				}
  compute_attributes		{
				  var cpus string
				  var start string
				  var cluster string
				  var memory_mb string
				  var guest_id string
				}
}

type compute_attributes struct {
  CPUS			string
  START			string
  CLUSTER		string
  MEMORY_MB		string
  GUEST_ID		string
} 

type network_if struct {
  COMPUTE_TYPE		string
  COMPUTE_NETWORK	string	
}

type volume struct {
  DATASTORE		string
  SIZE_GB		string
  THIN			string
  EAGER_ZERO		string
}

type machine struct {
  USERNAME		string
  PASSWORD		string
  SERVER		string
  NAME 			string
  LOCATION_ID 		string
  ORGANIZATION_ID 	string
  ENVIRONMENT_ID	string
  IP			string
  MAC			string
  ARCHITECTURE_ID	string
  DOMAIN_ID		string
  REALM_ID		string
  PUPPET_PROXY_ID	string
  PUPPET_CLASS_IDS	string
  OPERATINGSYSTEM_ID	string
  MEDIUM_ID		string
  PARTITION_TABLE_ID	string
  SUBNET_ID		string
  COMPUTE_RESOURCE_ID	string
  COMPUTE_PROFILE_ID	string
  ROOT_PASS		string
  MODEL_ID		string
  HOSTGROUP_ID		string
  OWNER_ID		string
  IMAGE_ID		string
  PUPPET_CA_PROXY_ID	string
  HOST_PARAMETERS_ATTRIBUTES	string 
  BUILD			string
  ENABLED		string
  PROVISION_METHOD	string
  PROVISION_REPORT_ID	string
  MANAGED		string
  CAPABILITIES		string
  COMMENT		string
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

type hammerArgs struct {
	options    []string
	subcommand string
	args       []string
}

func hammerCLI(h *hammerArgs, bVal string) (output []byte, err error) {
        argStr := ""
	build := bVal	
	
	if strings.ToLower(build) == "true"{
	 argStr += " create "
         for i := 0; i < len(h.args); i++ {
           argStr += h.args[i]
         }
	}else{
	 argStr = "list"	
	 }
 
        //print("/usr/bin/hammer ", h.subcommand, argStr)
	 cmd := exec.Command("/usr/bin/hammer", h.subcommand, argStr)
	 print("/usr/bin/hammer ",h.subcommand, argStr)
	
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(stdout))
	return stdout, err
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
	h := hammerArgs{subcommand: "host "}
        vm := machine{
          NAME: d.Get("name").(string),
        }
        iface := network_if{}
	cattr := compute_attributes{}
	vol   := volume{}	
 	

/* populate vm struct instance */
        if v, ok := d.GetOk("username"); ok {
          vm.USERNAME = v.(string)
        }
        if v, ok := d.GetOk("password"); ok {
          vm.PASSWORD = v.(string)
        }
        if v, ok := d.GetOk("server"); ok {
          vm.SERVER = v.(string)
        }
        if v, ok := d.GetOk("location-id"); ok {
          vm.LOCATION_ID = v.(string)
        }
        if v,ok := d.GetOk("organization-id"); ok{
          vm.ORGANIZATION_ID = v.(string)
        }
        if v,ok := d.GetOk("environment-id"); ok{
          vm.ENVIRONMENT_ID = v.(string)
        }
        if v,ok := d.GetOk("ip"); ok{
          vm.IP = v.(string)
        }
        if v,ok := d.GetOk("mac"); ok{
          vm.MAC = v.(string)
        }
        if v,ok := d.GetOk("architecture-id"); ok{
          vm.ARCHITECTURE_ID = v.(string)
        }
        if v,ok := d.GetOk("domain-id"); ok{
          vm.DOMAIN_ID = v.(string)
        }
        if v,ok := d.GetOk("realm-id"); ok{
          vm.REALM_ID = v.(string)
        }
        if v,ok := d.GetOk("puppet-proxy-id"); ok{
          vm.PUPPET_PROXY_ID = v.(string)
        }
        if v,ok := d.GetOk("puppet-class-ids"); ok{
          vm.PUPPET_CLASS_IDS = v.(string)
        }
        if v,ok := d.GetOk("operatingsystem-id"); ok{
          vm.OPERATINGSYSTEM_ID = v.(string)
        }
        if v,ok := d.GetOk("medium-id"); ok{
          vm.MEDIUM_ID = v.(string)
        }
        if v,ok := d.GetOk("partition-table-id"); ok{
          vm.PARTITION_TABLE_ID = v.(string)
        }
        if v,ok := d.GetOk("subnet-id"); ok{
          vm.SUBNET_ID = v.(string)
        }
        if v,ok := d.GetOk("computer-resource-id"); ok{
          vm.COMPUTE_RESOURCE_ID = v.(string)
        }
        if v,ok := d.GetOk("root-pass"); ok{
          vm.ROOT_PASS = v.(string)
        }
        if v,ok := d.GetOk("model-id"); ok{
          vm.MODEL_ID = v.(string)
        }
        if v,ok := d.GetOk("hostgroup-id"); ok{
          vm.HOSTGROUP_ID = v.(string)
        }
        if v,ok := d.GetOk("owner-id"); ok{
          vm.OWNER_ID = v.(string)
        }
        if v,ok := d.GetOk("puppet-ca-proxy-id"); ok{
          vm.PUPPET_CA_PROXY_ID = v.(string)
        }
        if v,ok := d.GetOk("image-id"); ok{
          vm.IMAGE_ID = v.(string)
        }
        if v,ok := d.GetOk("host-parameters-attributes"); ok{
          vm.HOST_PARAMETERS_ATTRIBUTES = v.(string)
        }
        if v,ok := d.GetOk("build"); ok{
          vm.BUILD = v.(string)
        }
        if v,ok := d.GetOk("enabled"); ok{
          vm.ENABLED = v.(string)
        }
        if v,ok := d.GetOk("provision-method"); ok{
          vm.PROVISION_METHOD = v.(string)
        }
        if v,ok := d.GetOk("managed"); ok{
          vm.MANAGED = v.(string)
        }
        if v,ok := d.GetOk("provision-report-id"); ok{
          vm.PROVISION_REPORT_ID = v.(string)
        }
        if v,ok := d.GetOk("comment"); ok{
          vm.COMMENT = v.(string)
        }
        if v,ok := d.GetOk("capabilities"); ok{
          vm.CAPABILITIES = v.(string)
        }
        if v,ok := d.GetOk("compute-profile-id"); ok{
          vm.COMPUTE_PROFILE_ID = v.(string)
        }
        if v,ok := d.GetOk("organization-id"); ok{
          vm.ORGANIZATION_ID = v.(string)
        }


/* populate iface struct instance */
        if v, ok := d.GetOk("interface-attributes.compute_type"); ok {
          iface.COMPUTE_TYPE = v.(string)
	}
        if v, ok := d.GetOk("interface-attributes.compute_network"); ok {
          iface.COMPUTE_NETWORK = v.(string)
	}

/* populate cattr struct instance */
        if v, ok := d.GetOk("compute-attributes.cpus"); ok {
          cattr.CPUS = v.(string)
	}
        if v, ok := d.GetOk("compute-attributes.start"); ok {
          cattr.START = v.(string)
	}
        if v, ok := d.GetOk("compute-attributes.cluster"); ok {
          cattr.CLUSTER = v.(string)
	}
        if v, ok := d.GetOk("compute-attributes.memory_mb"); ok {
          cattr.MEMORY_MB = v.(string)
	}
        if v, ok := d.GetOk("compute-attributes.guest_id"); ok {
          cattr.GUEST_ID = v.(string)
	}

/* populate vol struct instance */
        if v, ok := d.GetOk("storage-attributes.datastore"); ok {
          vol.DATASTORE = v.(string)
	}
        if v, ok := d.GetOk("storage-attributes.size_gb"); ok {
          vol.SIZE_GB = v.(string)
	}
        if v, ok := d.GetOk("storage-attributes.thin"); ok {
          vol.THIN = v.(string)
	}
        if v, ok := d.GetOk("storage-attributes.eager_zero"); ok {
          vol.EAGER_ZERO = v.(string)
	}

/* populate h.args with structs data from vm, iface, and vol */
/* Trying to iterate dynamically */
 	svm := reflect.ValueOf(&vm).Elem()
	typeOfvm := svm.Type()
	for i := 0; i<svm.NumField(); i++{
	  f := svm.Field(i)
	  fName := typeOfvm.Field(i).Name
          switch {
	   case fName == "BUILD":
            continue
	   case fName == "COMMENT":
	    if (f.Interface() != "") {
	     lStr := fmt.Sprintf("--%s \"%v\" ", strings.ToLower(strings.Replace(fName, "_", "-", -1)) , f.Interface())
	     h.args = append(h.args, lStr)
           } 
	   default:
	    if (f.Interface() != "") {
	     lStr := fmt.Sprintf("--%s %v ", strings.ToLower(strings.Replace(fName, "_", "-", -1)) , f.Interface())
	     h.args = append(h.args, lStr)
           } 
          }
	}

 	sif := reflect.ValueOf(&iface).Elem()
	typeOfiface := sif.Type()
        var ifStr string
	ifStr = "" 
	for i := 0; i<sif.NumField(); i++{
	  f := sif.Field(i)
	  fName := typeOfiface.Field(i).Name
	  if (f.Interface() != "") {
	   if ifStr == "" {
	    ifStr += fmt.Sprintf("%s=%v", strings.ToLower(fName) , f.Interface())
	   }else{
	    ifStr += fmt.Sprintf(",%s=%v", strings.ToLower(fName) , f.Interface())
	   }
          } 
	}
	if ifStr != "" {
	  h.args = append(h.args, fmt.Sprintf("--interface=\"%s\" ",ifStr))
	}

 	svol := reflect.ValueOf(&vol).Elem()
	typeOfvol := svol.Type()
        var volStr string
	volStr = ""
	for i := 0; i<svol.NumField(); i++{
	  f := svol.Field(i)
	  fName := typeOfvol.Field(i).Name
	  if (f.Interface() != "") {
	   if volStr == ""{
	   volStr += fmt.Sprintf("%s=%v", strings.ToLower(fName) , f.Interface())
	   }else{
	   volStr += fmt.Sprintf(",%s=%v", strings.ToLower(fName) , f.Interface())
	   	
	   }
          } 
	}

	if volStr != "" {
	  h.args = append(h.args, fmt.Sprintf("--volume=\"%s\" ",volStr))
	}

 	scattr := reflect.ValueOf(&cattr).Elem()
	typeOfcattr := scattr.Type()
        var cattrStr string
	cattrStr = ""
	for i := 0; i<scattr.NumField(); i++{
	  f := scattr.Field(i)
	  fName := typeOfcattr.Field(i).Name
	  if (f.Interface() != "") {
	   if cattrStr == "" {
	    cattrStr += fmt.Sprintf("%s=%v", strings.ToLower(fName) , f.Interface())
	   }else{
	    cattrStr += fmt.Sprintf(",%s=%v", strings.ToLower(fName) , f.Interface()) 
	   }
           } 
	}
	if cattrStr != "" {
	  h.args = append(h.args, fmt.Sprintf("--compute-attributes=\"%s\" ",cattrStr))
	}

	output, err := hammerCLI(&h,vm.BUILD)
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
