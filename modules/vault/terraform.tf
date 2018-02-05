variable "ami" {}
variable "dns_zone_id" {}
variable "vpc_id" {}
variable "ec2type" {}
variable "subnet-a" {}
variable "subnet-b" {}
variable "subnet-c" {}

resource "aws_instance" "vault-a" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.vault-server.id}"]
  subnet_id              = "${var.subnet-a}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"
  iam_instance_profile = "consul-server-profile"
  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }
  tags {
    Name = "vault-a.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
    vault = "usw2"
  }
}

resource "aws_route53_record" "vault-a" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "vault-a.e.tux.rocks"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.vault-a.private_ip}" ]
}

resource "aws_instance" "vault-b" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.vault-server.id}"]
  subnet_id              = "${var.subnet-b}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"
  iam_instance_profile = "consul-server-profile"
  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "vault-b.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
    vault = "usw2"
  }
}

resource "aws_route53_record" "vault-b" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "vault-b.e.tux.rocks"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.vault-b.private_ip}" ]
}

resource "aws_instance" "vault-c" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.vault-server.id}"]
  subnet_id              = "${var.subnet-c}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"
  iam_instance_profile = "consul-server-profile"
  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "vault-c.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
    vault = "usw2"
  }
}

resource "aws_route53_record" "vault-c" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "vault-c.e.tux.rocks"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.vault-c.private_ip}" ]
}

resource "aws_security_group" "vault-server" {
  provider    = "aws.local"
  name        = "vault-server"
  description = "vault-server"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["10.0.0.0/8"] 
  }

  ingress { 
    from_port = 8200
    to_port = 8201
    protocol = "tcp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    from_port = 8
    to_port = 0
    protocol = "icmp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "vault-server"
  }
}

resource "aws_elb" "vault" {
  provider = "aws.local"
  name = "vault"
  internal = "true"
  subnets = [ "${var.subnet-a}", "${var.subnet-b}", "${var.subnet-c}" ]
  security_groups = [ "${aws_security_group.vault-server.id}" ]

  listener {
    instance_port = 8200
    instance_protocol = "tcp"
    lb_port = 8200
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    interval = 10
    target = "TCP:8200"
  }

  instances = [ "${aws_instance.vault-a.id}", "${aws_instance.vault-b.id}", "${aws_instance.vault-c.id}" ]
  cross_zone_load_balancing = true
  tags {
    Name = "vault.e.tux.rocks"
  }
}

resource "aws_route53_record" "vault" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "vault.e.tux.rocks"
  type     = "A"
  alias {
    name = "${aws_elb.vault.dns_name}"
    zone_id = "${aws_elb.vault.zone_id}"
    evaluate_target_health = "true"
  }
}
