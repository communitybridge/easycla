variable "vpc_id" {}

variable "ghe_sg" {}

variable "elb_sg" {}

variable "replica_count" {}

variable "internal_subnets" {
  type = "list"
}

variable "external_subnets" {
  type = "list"
}

data "aws_ami" "github-enterprise-ami" {
  provider         = "aws.local"
  most_recent      = true

  filter {
    name   = "owner-id"
    values = ["895557238572"]
  }

  filter {
    name   = "name"
    values = ["GitHub Enterprise 2.11.2"]
  }
}

resource "aws_instance" "ghe-master" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.github-enterprise-ami.id}"
  vpc_security_group_ids = ["${var.ghe_sg}"]
  subnet_id              = "${var.internal_subnets[2]}"
  instance_type          = "r3.large"

  ebs_block_device {
    device_name = "/dev/xvdf"
    volume_size = 250
    encrypted   = true
  }

  tags {
    Name = "[M] GitHub Enterprise"
  }
}

resource "aws_instance" "ghe-replicas" {
  provider               = "aws.local"
  count                  = "${var.replica_count}"
  ami                    = "${data.aws_ami.github-enterprise-ami.id}"
  vpc_security_group_ids = ["${var.ghe_sg}"]
  subnet_id              = "${var.internal_subnets[count.index]}"
  instance_type          = "r3.large"

  ebs_block_device {
    device_name = "/dev/xvdf"
    volume_size = 250
    encrypted   = true
  }

  tags {
    Name = "[R] GitHub Enterprise"
  }
}

# Create a new load balancer
resource "aws_elb" "ghe" {
  name               = "github-enterprise"
  security_groups    = ["${var.elb_sg}"]
  subnets            = ["${var.internal_subnets}"]
  internal           = true

  listener {
    instance_port     = 22
    instance_protocol = "tcp"
    lb_port           = 22
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  listener {
    instance_port      = 443
    instance_protocol  = "https"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/43931226-dabf-4ad4-88e7-069f07e28edd"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTPS:443/status"
    interval            = 30
  }

  instances                   = ["${aws_instance.ghe-master.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 400

  tags {
    Name = "github-enterprise"
  }
}

resource "aws_route53_record" "www" {
  provider= "aws.local"
  zone_id = "Z2MEXX9ZCWUHX6"
  name    = "www"
  type    = "A"

  alias {
    name                   = "${aws_elb.ghe.dns_name}"
    zone_id                = "${aws_elb.ghe.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "dot" {
  provider= "aws.local"
  zone_id = "Z2MEXX9ZCWUHX6"
  name    = "."
  type    = "A"

  alias {
    name                   = "${aws_elb.ghe.dns_name}"
    zone_id                = "${aws_elb.ghe.zone_id}"
    evaluate_target_health = true
  }
}
