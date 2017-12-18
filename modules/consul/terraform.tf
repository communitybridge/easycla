variable "ami" {}
variable "dns_zone_id" {}
variable "vpc_id" {}
variable "ec2type" {}
variable "subnet-a" {}
variable "subnet-b" {}
variable "subnet-c" {}

resource "aws_route53_record" "consul" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "consul.e.tux.rocks"
  type     = "NS"
  ttl      = 300
  records = [ "${aws_route53_record.consul-a.name}.", "${aws_route53_record.consul-b.name}.", "${aws_route53_record.consul-c.name}." ]
}

resource "aws_instance" "consul-a" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.consul-server.id}"]
  subnet_id              = "${var.subnet-a}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"
  iam_instance_profile = "${aws_iam_instance_profile.consul-server-profile.name}"
  root_block_device {
    volume_type = "gp2"
    volume_size = "20"
  }
  tags {
    Name = "consul-a.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
    consul = "usw2"
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
  iam_instance_profile = "${aws_iam_instance_profile.consul-server-profile.name}"
  root_block_device {
    volume_type = "gp2"
    volume_size = "20"
  }

  tags {
    Name = "consul-b.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
    consul = "usw2"
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

resource "aws_instance" "consul-c" {
  provider               = "aws.local"
  ami                    = "${var.ami}" 
  vpc_security_group_ids = ["${aws_security_group.consul-server.id}"]
  subnet_id              = "${var.subnet-c}"
  instance_type          = "${var.ec2type}"
  key_name               = "dan"
  iam_instance_profile = "${aws_iam_instance_profile.consul-server-profile.name}"
  root_block_device {
    volume_type = "gp2"
    volume_size = "20"
  }

  tags {
    Name = "consul-c.e.tux.rocks"
    Owner = "dparsons@linuxfoundation.org"
    consul = "usw2"
  }
}

resource "aws_route53_record" "consul-c" {
  provider = "aws.local"
  zone_id  = "${var.dns_zone_id}"
  name     = "consul-c.e.tux.rocks"
  type     = "A"
  ttl      = 300
  records = [ "${aws_instance.consul-c.private_ip}" ]
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

resource "aws_iam_role" "consul-server-role" {
  provider = "aws.local"
  name = "consul-server-role"
  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "consul-autojoin-policy" {
  provider = "aws.local"
  name = "consul-autojoin-policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:DescribeInstances",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "consul-autojoin-attachment" {
  provider = "aws.local"
  name = "consul-autojoin-attachment"
  roles = ["${aws_iam_role.consul-server-role.name}"]
  policy_arn = "${aws_iam_policy.consul-autojoin-policy.arn}"
}

resource "aws_iam_instance_profile" "consul-server-profile" {
  provider = "aws.local"
  name = "consul-server-profile"
  role = "${aws_iam_role.consul-server-role.name}"
}

output "consul-server-profile" {
#  value = "${aws_iam_instance_profile.consul-server-profile.name}"
  value = "yo"
}
