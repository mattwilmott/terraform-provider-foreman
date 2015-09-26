package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSCreate,
		Read:   resourceDNSRead,
		Update: resourceDNSUpdate,
		Delete: resourceDNSDelete,

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address_csv": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceDNSRead(d, m)
	if err != nil {
		return err
	}

	d.SetId(d.Get("host").(string))
	return nil
}

func resourceDNSRead(d *schema.ResourceData, m interface{}) error {
	ips, err := getIPs(d.Get("host").(string))
	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return fmt.Errorf("No IP addresses found for %s", d.Get("host"))
	}

	d.Set("ip_address", ips[0])
	d.Set("ip_address_csv", strings.Join(ips, ","))

	return nil
}

func resourceDNSUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDNSDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getIPs(host string) ([]string, error) {
	ips := []string{}

	adds, err := net.LookupIP(host)
	if err != nil {
		return ips, err
	}

	for _, v := range adds {
		ips = append(ips, v.String())
	}

	return ips, nil
}
