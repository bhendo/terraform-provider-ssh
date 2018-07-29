package main

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH hostname or IP address e.g. (bastion-host.local)",
			},
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username for logging into the host (e.g. ec2-user)",
			},
			"key_file_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to a private key file (e.g. ~/.ssh/id_rsa)",
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Connection timeout in seconds (default = 10)",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ssh_tunnel": dataSourceSSHTunnel(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	user := d.Get("user").(string)
	timeout := d.Get("timeout").(int)
	keyFilePath := d.Get("key_file_path").(string)

	config := &Config{
		Host:        host,
		User:        user,
		KeyFilePath: keyFilePath,
		Timeout:     time.Duration(timeout),
	}
	configs = append(configs, config)
	return config, nil
}
