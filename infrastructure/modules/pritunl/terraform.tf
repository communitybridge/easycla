variable "external_subnets" {
  description = "External VPC Subnets"
  type = "list"
}

variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "sg" {
  description = "Security Group for the VPN Server"
}

variable "elb_sg" {
  description = "Security Group for the VPN Server"
}

variable "mongodb_sg" {
  description = "Security Group for the VPN Server"
}

variable "dns_name" {}

variable "dns_zone_id" {}

data "aws_ami" "pritunl" {
  provider    = "aws.local"
  most_recent = true

  filter {
    name = "name"
    values = ["pritunl-2*"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "mongodb" {
  provider               = "aws.local"
  ami                    = "ami-c60d68a6"
  instance_type          = "t2.small"
  subnet_id              = "${element(var.internal_subnets, 1)}"
  key_name               = "infrastructure"
  vpc_security_group_ids = ["${var.mongodb_sg}"]
  monitoring             = true

  tags {
    Name        = "(VPN) MongoDB Storage"
  }
}

resource "aws_launch_configuration" "main" {
  provider    = "aws.local"
  name_prefix = "pritunl-cluster"

  image_id                    = "${data.aws_ami.pritunl.id}"
  instance_type               = "c4.large"
  key_name                    = "infrastructure"
  security_groups             = ["${var.sg}"]
  associate_public_ip_address = true

  # root
  root_block_device {
    volume_type = "gp2"
    volume_size = "30"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "main" {
  provider             = "aws.local"
  name                 = "pritunl-cluster"

  availability_zones   = ["us-west-2a", "us-west-2b", "us-west-2c"]
  vpc_zone_identifier  = ["${var.external_subnets}"]
  launch_configuration = "${aws_launch_configuration.main.id}"
  min_size             = "2"
  max_size             = "2"
  desired_capacity     = "2"
  termination_policies = ["OldestLaunchConfiguration", "Default"]

  tag {
    key                 = "Name"
    value               = "(VPN) Pritunl Node"
    propagate_at_launch = true
  }

  load_balancers = ["${aws_elb.pritunl-cluster.name}"]

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "pritunl-cluster" {
  provider           = "aws.local"
  name               = "pritunl-cluster"
  subnets            = ["${var.external_subnets}"]
  security_groups    = ["${var.elb_sg}"]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  listener {
    instance_port      = 9700
    instance_protocol  = "http"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:643009352547:certificate/bd4bc0c6-ca71-4264-a3f8-273dcd541ba1"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:9700/ping"
    interval            = 30
  }

  cross_zone_load_balancing   = true
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400

  tags {
    Name = "pritunl-cluster"
  }
}

resource "aws_route53_record" "pritunl-cluster" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "galaxy.${var.dns_name}"
  type     = "A"

  alias {
    name                   = "${aws_elb.pritunl-cluster.dns_name}"
    zone_id                = "${aws_elb.pritunl-cluster.zone_id}"
    evaluate_target_health = true
  }
}
