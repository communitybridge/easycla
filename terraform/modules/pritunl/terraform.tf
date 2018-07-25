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

resource "aws_instance" "pritunl-1" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.pritunl.id}"
  vpc_security_group_ids = ["${var.sg}"]
  subnet_id              = "${var.external_subnets[0]}"
  instance_type          = "c5.large"
  key_name               = "engineering-production"

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "Pritunl VPN Node"
  }
}

resource "aws_instance" "pritunl-2" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.pritunl.id}"
  vpc_security_group_ids = ["${var.sg}"]
  subnet_id              = "${var.external_subnets[1]}"
  instance_type          = "c5.large"
  key_name               = "engineering-production"

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "Pritunl VPN Node"
  }
}

resource "aws_eip" "pritunl1" {
  provider    = "aws.local"
  vpc = true
  instance                  = "${aws_instance.pritunl-1.id}"
}

resource "aws_eip" "pritunl2" {
  provider    = "aws.local"
  vpc = true
  instance                  = "${aws_instance.pritunl-2.id}"
}

# Create a new load balancer
resource "aws_elb" "pritunl-node-cluster" {
  provider           = "aws.local"
  name               = "pritunl-node-cluster"
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
    ssl_certificate_id = "arn:aws:acm:us-west-2:643009352547:certificate/4938ed7c-e270-4597-84b2-6374db6149f4"
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
    Name = "pritunl-node-cluster"
  }
}

resource "aws_elb_attachment" "pritunl1" {
  provider    = "aws.local"
  elb      = "${aws_elb.pritunl-node-cluster.id}"
  instance = "${aws_instance.pritunl-1.id}"
}

resource "aws_elb_attachment" "pritunl2" {
  provider    = "aws.local"
  elb      = "${aws_elb.pritunl-node-cluster.id}"
  instance = "${aws_instance.pritunl-2.id}"
}

resource "aws_route53_record" "pritunl-node-cluster" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "vpn"
  type     = "A"

  alias {
    name                   = "${aws_elb.pritunl-node-cluster.dns_name}"
    zone_id                = "${aws_elb.pritunl-node-cluster.zone_id}"
    evaluate_target_health = true
  }
}
