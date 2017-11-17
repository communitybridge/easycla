variable "ami" {}
variable "subnet" {}
variable "name" {}
variable "dns_zone_id" {}
variable "vpc_id" {}

resource "aws_instance" "danvpn" {
  provider               = "aws.local"
  ami                    = "${var.ami}" // Amazon Linux 2017.09.1 (HVM) SSD Volume Type
  vpc_security_group_ids = ["${aws_security_group.danvpn-sg.id}"]
  subnet_id              = "${var.subnet}"
  instance_type          = "t2.small"
  key_name               = "dan"

  ebs_block_device {
    device_name = "/dev/xvdf"
    volume_size = 50
    encrypted   = true
  }

  tags {
    Name = "${var.name}"
    Owner = "dparsons@linuxfoundation.org"
  }
}

resource "aws_eip" "danvpn" {
  provider = "aws.local"
  instance = "${aws_instance.danvpn.id}"
  vpc = true
}

resource "aws_route53_record" "danvpn" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "${var.name}"
  type     = "A"
  ttl      = 300
  records = [ "${aws_eip.danvpn.public_ip}" ]
}

resource "aws_security_group" "danvpn-sg" {
  provider    = "aws.local"
  name        = "danvpn-sg"
  description = "Dans OpenVPN Server"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["67.160.133.122/32"] // dan's home IP. can be removed later
  }

  ingress {
    from_port = 2112
    to_port = 2112
    protocol = "tcp"
    cidr_blocks = ["67.160.133.122/32"]
  }

  ingress {
    from_port = 2112
    to_port = 2112
    protocol = "udp"
    cidr_blocks = ["67.160.133.122/32"]
  }

  ingress {
    from_port = 8
    to_port = 0
    protocol = "icmp"
    cidr_blocks = ["67.160.133.122/32"]
  }

   egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Dans OpenVPN Server"
  }
}