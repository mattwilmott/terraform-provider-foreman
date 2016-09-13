# terraform-foreman
Foreman provider for Terraform - !Please note that this is by no means production ready!

** Note, I am a sysops guy and have never touched golang in my life. I started this project to see if I could provision machines in foreman with terraform.

## What is Terraform
Terraform is an orchestration tool that can be used to deploy and manage the lifecycle of cloud resources, virtual machines, physical machines, DNS records. Effectively, if you can code it you can provision it.
You can find more information at https://terraform.io/

## What is Foreman
Foreman (or more accurately TheForeman) is a Red Hat sponsored open source tool used to manage infrastructure in a Private Cloud. Foreman can provision for instance Docker instances, Virtual Machines in OpenStack or VMWare and it can even deploy to bare metal. Foreman ties into PuppetLabs Puppet infrastructure and provides ENC data regarding the servers it manages. Check it out at http://theforeman.org/

## Usage


Add the following to your ~/.terraformrc
```
providers {
     foreman = "/path/to/bin/terraform-provider-foreman"
}
```

Sample terraform config. Save this as a example.tf or similar
```

resource "foreman_server" "myVM" {
  username = "username"
  password = "password"
  url = "https://foreman.domain.com/api"
  name = "hostname"
  environment_id = "environment_id"
  ip = "10.0.0.2"
  mac	= "ff:ff:ff:ff:ff:ff"
  architecture_id = 1
  domain_id = 1
  realm_id = 1
  puppet_proxy_id	= 1
  puppetclass_ids	= [1,2,3]
  operatingsystem_id = 1
  medium_id = "1" # yeah, it's a string. idk, don't ask
  ptable_id	= 1
  subnet_id	= 1
  compute_resource_id	= 25
  root_pass	= "Superawesomehash"
  model_id = 1
  hostgroup_id = 1
  owner_id = 1
  owner_type = "User" # must be either User or Usergroup
  puppet_ca_proxy_id = 1
  image_id = 1
  build	= true
  enabled	= true
  provision_method = "build"
  managed	= true
  progress_report_id = "progress_report_id"
  comment = "Build purpose"
  capabilities = "Something"
  compute_profile_id = 1
	host_parameters_attributes {
    roles = "server_role"
    puppet = "true"
    chef = "false"
  }
  interfaces_attributes{
    mac = "ff:ff:ff:ff:ff:ff"
    ip = "ip"
    name = "name"
    subnet_id = 1
    domain_id = 1
    identifier = "identifier"
    managed = true
    primary = true
    provision = true
    username = "username" # only for bmc
    password = "password" # only for bmc
    provider = "provider" # only accepted IPMI
    virtual = false
    tag = "tag"
    attached_to = "something"
    mode = "mode" # with validations
    attached_devices 	=	[]string
    bond_options = "bond opts"
  }
  compute_attributes {
    cpus = "2"
  	start = "1"
  	cluster = "clustername"
  	memory_mb = "2048"
  	guest_id = "guest_id"
  }
	volumes_attributes{
    name = "name"
  	size_gb	= 16
  	_delete	= "false"
  	datastore	"Datastore_name"
  }
}
```

**TODO**


if you have already planned the terraform resources you can taint them essentially marking them to be rebuilt

Now you can write your own main.tf similar to the example. Reference the terraform documentation at https://terraform.io/intro/getting-started/build.html

Once a plan has been created you are ready to apply the plan and actually deploy.

**This provider is an interface between terraform and the Foreman API. Terraform providers are typically fully CRUD designed but so far I have only made the create functionality. I will work on the others as necessary.**

```
# In the directory of your my_custom_terraform.tf file
terraform apply
```
This will interrogate the provider and make the changes. If a server isn't built for instance it will call the create method and the provider will instantiate the resource.

## How to Build the Source
In order to build/install the source, navigate to the checked out directory and ensure your $GOPATH is defined.

Execute
```
go get
go build
```
