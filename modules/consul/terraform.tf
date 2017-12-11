variable "ami" {}
variable "dns_zone_id" {}
variable "vpc_id" {}
variable "ec2type" {}
variable "subnet-a" {}
variable "subnet-b" {}

resource "aws_instance" "consul-a" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.consul-server.id}"]
  subnet_id              = "${var.subnet-a}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"

  root_block_device {
    volume_type = "gp2"
    volume_size = "20"
  }

  tags {
    Name = "consul-a.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
  }
}

resource "aws_route53_record" "consul-a" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "consul-a.e.tux.rocks"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.consul-a.private_ip}" ]
}

resource "aws_instance" "consul-b" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.consul-server.id}"]
  subnet_id              = "${var.subnet-b}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"

  root_block_device {
    volume_type = "gp2"
    volume_size = "20"
  }

  tags {
    Name = "consul-b.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
  }
}

resource "aws_route53_record" "consul-b" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "consul-b.e.tux.rocks"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.consul-b.private_ip}" ]
}

resource "aws_security_group" "consul-server" {
  provider    = "aws.local"
  name        = "consul-server"
  description = "consul-server"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["10.0.0.0/8"] 
  }

  ingress {
    from_port = 53
    protocol = "tcp"
    to_port = 53
    cidr_blocks = ["10.0.0.0/8"] 
  }

  ingress {
    from_port = 53
    protocol = "udp"
    to_port = 53
    cidr_blocks = ["10.0.0.0/8"] 
  }

  ingress { 
    from_port = 8300
    to_port = 8302
    protocol = "tcp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress { 
    from_port = 8300
    to_port = 8302
    protocol = "udp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress { 
    from_port = 8400
    to_port = 8400
    protocol = "tcp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress { 
    from_port = 80
    to_port = 80
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
    Name        = "consul-server"
  }
}