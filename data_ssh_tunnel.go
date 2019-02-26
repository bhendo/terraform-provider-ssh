package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	localHost   = "127.0.0.1"
	successText = "Authentication succeeded"
)

var (
	errorConnectionTimeout    = fmt.Errorf("connection timeout")
	errorFailedPortAllocation = fmt.Errorf("failed to allocate a local port")
)

func dataSourceSSHTunnel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHTunnelRead,
		Schema: map[string]*schema.Schema{
			"local_host": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The local host address (always 127.0.0.1)",
			},
			"local_port": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The local bind port (e.g. 62774)",
			},
			"remote_host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The remote bind host (e.g. internal-postgres-server)",
			},
			"remote_port": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The remote bind port (e.g. 5432)",
			},
			"command_args": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The command and arguments used to create the tunnel (e.g. ssh -N -L 65534:remote_host:5432 user@ssh_host)",
			},
		},
	}
}

func dataSourceSSHTunnelRead(d *schema.ResourceData, meta interface{}) error {
	remoteHost := d.Get("remote_host").(string)
	remotePort := d.Get("remote_port").(string)
	config := meta.(*Config)

	localPort, err := newPort()
	if err != nil {
		return errorFailedPortAllocation
	}
	d.Set("local_host", localHost)
	d.Set("local_port", localPort)
	d.SetId(localPort)

	args := []string{"-v", "-N", "-oStrictHostKeyChecking=accept-new"}
	args = append(args, "-L", fmt.Sprintf("%s:%s:%s", localPort, remoteHost, remotePort))
	if config.KeyFilePath != "" {
		args = append(args, "-i", config.KeyFilePath)
	}
	args = append(args, fmt.Sprintf("%s@%s", config.User, config.Host))

	d.Set("command_args", fmt.Sprintf("ssh %s", strings.Join(args, " ")))

	cmd := exec.Command("ssh", args...)
	stdErrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	config.AddCmd(cmd)

	go func() {
		cmd.Wait()
	}()

	success := make(chan bool, 1)

	scanner := bufio.NewScanner(stdErrPipe)
	go func() {
		for {
			if scanner.Scan() && strings.Contains(scanner.Text(), successText) {
				success <- true
			}
		}
	}()

	select {
	case <-success:
		return nil
	case <-time.After(time.Second * config.Timeout):
		return errorConnectionTimeout
	}
}

func newPort() (port string, err error) {
	server, err := net.Listen("tcp", ":0")

	if err != nil {
		return "", err
	}

	hostString := server.Addr().String()

	if err := server.Close(); err != nil {
		return "", err
	}
	_, portString, err := net.SplitHostPort(hostString)

	if err != nil {
		return "", err
	}

	return portString, nil
}
