variable "ami" {}
variable "subnet" {}
variable "name" {}
variable "dns_zone_id" {}
variable "vpc_id" {}
variable "ec2type" {}

resource "aws_instance" "salt" {
  provider               = "aws.local"
  ami                    = "${var.ami}" // Amazon Linux 2017.09.1 (HVM) SSD Volume Type
  vpc_security_group_ids = ["${aws_security_group.salt-sg.id}"]
  subnet_id              = "${var.subnet}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "${var.name}"
    Owner = "dparsons@linuxfoundation.org"
  }
}

resource "aws_route53_record" "salt" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "${var.name}"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.salt.private_ip}" ]
}

resource "aws_security_group" "salt-sg" {
  provider    = "aws.local"
  name        = "salt-sg"
  description = "salt master"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["10.0.0.0/8"] // dan's home IP. can be removed later
  }

  ingress { // ports needed for salt minions to talk to salt master
    from_port = 4505
    to_port = 4506
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
    Name        = "salt master"
  }
}