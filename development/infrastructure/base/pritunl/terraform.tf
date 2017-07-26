variable "external_subnets" {
  description = "External VPC Subnets"
  type = "list"
}

variable "vpn_sg" {
  description = "Security Group for the VPN Server"
}

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

resource "aws_instance" "pritunl" {
  count                  = "${length(var.external_subnets)}"
  provider               = "aws.local"
  ami                    = "${data.aws_ami.pritunl.id}"
  source_dest_check      = false
  instance_type          = "t2.small"
  subnet_id              = "${element(var.external_subnets, count.index)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${var.vpn_sg}"]
  monitoring             = true

  tags {
    Name        = "Pritunl - Node #${count.index}"
    Team        = "Engineering"
    Environment = "Production"
  }
}

resource "aws_eip" "pritunl" {
  provider = "aws.local"
  instance = "${element(aws_instance.pritunl.*.id, count.index)}"
  vpc      = true
}

# Create a new load balancer
resource "aws_elb" "pritunl-cluster" {
  name               = "pritunl-cluster"
  subnets            = ["${var.external_subnets}"]
  security_groups    = ["${var.vpn_sg}"]

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
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:9700/ping"
    interval            = 30
  }

  instances                   = ["${aws_instance.pritunl.*.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400

  tags {
    Name = "pritunl-cluster"
  }
}

resource "aws_route53_record" "pritunl-cluster" {
  zone_id = "Z2MDT77FL23F9B"
  name    = "vpn"
  type    = "A"

  alias {
    name                   = "${aws_elb.pritunl-cluster.dns_name}"
    zone_id                = "${aws_elb.pritunl-cluster.zone_id}"
    evaluate_target_health = true
  }
}
