variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "bind_sg" {
  description = "Security Group for the internal ELB"
}

variable "cidr" {}

variable "region" {}

variable "internal_elb_sg" {}

data "template_file" "bind-installation-bash" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"
}

data "aws_ami" "bind-consul-ami" {
  provider    = "aws.local"
  most_recent = true

  filter {
    name = "name"
    values = ["*Consul*"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }

  owners     = ["self"]
}

resource "aws_instance" "bind_1" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.bind-consul-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(var.internal_subnets, 0)}"
  key_name               = "production-shared-tools"
  iam_instance_profile   = "ecsInstanceRole"
  vpc_security_group_ids = ["${var.bind_sg}"]
  monitoring             = true
  user_data              = "${data.template_file.bind-installation-bash.rendered}"
  private_ip             = "${replace(var.cidr, ".0/24", ".140")}"

  tags {
    Name        = "Bind/Consul - Node #1"
    Team        = "Engineering"
    Environment = "Production"
    consul-node = "true"
  }
}

resource "aws_instance" "bind_2" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.bind-consul-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(var.internal_subnets, 1)}"
  key_name               = "production-shared-tools"
  iam_instance_profile   = "ecsInstanceRole"
  vpc_security_group_ids = ["${var.bind_sg}"]
  monitoring             = true
  user_data              = "${data.template_file.bind-installation-bash.rendered}"
  private_ip             = "${replace(var.cidr, ".0/24", ".180")}"

  tags {
    Name        = "Bind/Consul - Node #2"
    Team        = "Engineering"
    Environment = "Production"
    consul-node = "true"
  }
}

resource "aws_instance" "bind_3" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.bind-consul-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(var.internal_subnets, 2)}"
  key_name               = "production-shared-tools"
  iam_instance_profile   = "ecsInstanceRole"
  vpc_security_group_ids = ["${var.bind_sg}"]
  monitoring             = true
  user_data              = "${data.template_file.bind-installation-bash.rendered}"
  private_ip             = "${replace(var.cidr, ".0/24", ".220")}"

  tags {
    Name        = "Bind/Consul - Node #3"
    Team        = "Engineering"
    Environment = "Production"
    consul-node = "true"
  }
}

# Create a new load balancer
resource "aws_elb" "consul" {
  provider    = "aws.local"
  name        = "consul-bind-cluster"
  subnets     = ["${var.internal_subnets}"]
  security_groups = ["${var.internal_elb_sg}"]
  internal = true

  listener {
    instance_port = 8300
    instance_protocol = "tcp"
    lb_port = 8300
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8301
    instance_protocol = "tcp"
    lb_port = 8301
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8302
    instance_protocol = "tcp"
    lb_port = 8302
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8400
    instance_protocol = "tcp"
    lb_port = 8400
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8500
    instance_protocol = "http"
    lb_port = 8500
    lb_protocol = "http"
  }

  listener {
    instance_port = 8600
    instance_protocol = "tcp"
    lb_port = 53
    lb_protocol = "tcp"
  }

  health_check {
    target = "TCP:8300"
    healthy_threshold = 2
    unhealthy_threshold = 2
    interval = 30
    timeout = 5
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "Consul Cluster"
  }
}

resource "aws_route53_record" "consul" {
  zone_id = "Z2LGKA6SQQWJ1F"
  name    = "consul"
  type    = "A"

  alias {
    name                   = "${aws_elb.consul.dns_name}"
    zone_id                = "${aws_elb.consul.zone_id}"
    evaluate_target_health = true
  }
}

# Create a new load balancer attachment
resource "aws_elb_attachment" "consul_elb_1" {
  elb      = "${aws_elb.consul.id}"
  instance = "${aws_instance.bind_1.id}"
}

# Create a new load balancer attachment
resource "aws_elb_attachment" "consul_elb_2" {
  elb      = "${aws_elb.consul.id}"
  instance = "${aws_instance.bind_2.id}"
}

# Create a new load balancer attachment
resource "aws_elb_attachment" "consul_elb_3" {
  elb      = "${aws_elb.consul.id}"
  instance = "${aws_instance.bind_3.id}"
}

output "dns_servers" {
  value = ["${replace(var.cidr, ".0/24", ".140")}", "${replace(var.cidr, ".0/24", ".180")}", "${replace(var.cidr, ".0/24", ".220")}"]
}
