# terraform-provider-ssh

This provider enables SSH port forwarding in Terraform and leverages a standerd ssh client

**Note** An ssh client will need to be installed and properly configured in your path to work

## Inspiration

The idea for this provider came from [https://github.com/stefansundin/terraform-provider-ssh](https://github.com/stefansundin/terraform-provider-ssh)

## Example

```terraform
variable "database_password" {}
variable "database_host {}
variable "database_port {}

provider "ssh" {
  alias         = "ssh1"
  host          = "127.0.0.1"
  user          = "ec2-user"
  key_file_path = "~/.ssh/id_rsa"
  timeout       = 15
}
data "ssh_tunnel" "tunnel" {
  provider    = "ssh.ssh1"
  remote_host = "${var.database_host}"
  remote_port = "${var.database_port}"
}

provider "postgresql" {
  host     = "localhost"
  port     = "${data.ssh_tunnel.test.local_port}"
  database = "mydb"
  username = "myuser"
  password = "${var.password}"
  sslmode  = "require"
}

output "local_address" {
  value = "${data.ssh_tunnel.tunnel.local_host}:${data.ssh_tunnel.tunnel.local_port}"
}
```

### Gotcha

Like the provider by stefansundin this provider does not support applying plan files. However, there is a work around.

Create your own ssh tunnel using the output from `data.ssh_tunnel.tunnel_name.command_args` e.g. `ssh -N -L 65534:database:5432 user@host`