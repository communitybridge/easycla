/**
 * Creates basic security groups to be used by instances and ELBs.
 */

variable "vpc_id" {
  description = "The VPC ID"
}

resource "aws_security_group" "pritunl-elb" {
  provider    = "aws.local"
  name        = "pritunl-elb"
  description = "ELB in front of Pritunl Nodes"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl Load Balancer"
  }
}

resource "aws_security_group" "pritunl-node" {
  provider    = "aws.local"
  name        = "pritunl-node"
  description = "Pritunl Server Nodes"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    security_groups = ["${aws_security_group.pritunl-elb.id}"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  # Engineering VPN Server (Split Tunnel)
  ingress {
    from_port = 17813
    protocol = "udp"
    to_port = 17813
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl Server Nodes"
  }
}

resource "aws_security_group" "pritunl-mongodb" {
  provider    = "aws.local"
  name        = "pritunl-mongodb"
  description = "Pritunl MongoDB Storage"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 27017
    to_port   = 27017
    protocol  = "tcp"
    security_groups = ["${aws_security_group.pritunl-node.id}"]
  }

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Pritunl MongoDB Storage"
  }
}

resource "aws_security_group" "ghe" {
  provider    = "aws.local"
  name        = "github-enterprise"
  description = "GitHub Enterprise Appliance"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    security_groups = [
      "${aws_security_group.ghe_elb.id}"
    ]
  }

  ingress {
    from_port = 25
    to_port = 25
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    security_groups = [
      "${aws_security_group.ghe_elb.id}"
    ]
  }

  ingress {
    from_port = 122
    to_port = 122
    protocol = "tcp"
    cidr_blocks = [
      "10.32.0.0/12"
    ]
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }

//  ingress {
//    from_port = 122
//    to_port = 122
//    protocol = "tcp"
//    security_groups = [
//      "sg-e5b2d098"
//    ]
//  }

  ingress {
    from_port = 123
    to_port = 123
    protocol = "udp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 161
    to_port = 161
    protocol = "udp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 8443
    to_port = 8443
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    security_groups = [
      "${aws_security_group.ghe_elb.id}"
    ]
  }

  ingress {
    from_port = 1194
    to_port = 1194
    protocol = "udp"
    self = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "GitHub Enterprise Appliance"
  }
}


resource "aws_security_group" "ghe_elb" {
  provider    = "aws.local"
  name        = "github-enterprise-elb"
  description = "GitHub Enterprise Appliance"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["10.32.0.0/12"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    security_groups = [
      "sg-e5b2d098",
      "sg-0d2f4376",
      "sg-493b5732"
    ]
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    security_groups = [
      "sg-e5b2d098",
      "sg-0d2f4376",
      "sg-493b5732"
    ]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    security_groups = [
      "sg-e5b2d098",
      "sg-0d2f4376",
      "sg-493b5732"
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "GitHub Enterprise Appliance"
  }
}

resource "aws_security_group" "it-ghe-elb" {
  provider    = "aws.local"
  name        = "it-github-enterprise"
  description = "IT GHE Bridge"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "198.145.29.72/32",
      "198.145.29.65/32",
      "140.211.169.2/32",
      "140.211.169.30/32",
      "54.201.117.121/32" // for salt.e.tux.rocks
    ]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [
      "198.145.29.72/32",
      "198.145.29.65/32",
      "140.211.169.2/32",
      "140.211.169.30/32"
    ]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [
      "198.145.29.72/32",
      "198.145.29.65/32",
      "140.211.169.2/32",
      "140.211.169.30/32",
      "54.201.117.121/32" // for salt.e.tux.rocks
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "IT GHE Bridge"
  }
}

resource "aws_security_group" "infra-ecs-cluster" {
  provider    = "aws.local"
  name        = "infra-ecs"
  description = "Infra ECS Cluster"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    self            = true
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_groups = [
      "${aws_security_group.consul-elb.id}"
    ]
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_groups = [
      "${aws_security_group.nexus-elb.id}"
    ]
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_groups = [
      "${aws_security_group.mongo-elb.id}"
    ]
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_groups = [
      "${aws_security_group.logstash-elb.id}"
    ]
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_groups = [
      "${aws_security_group.vault-ecs-cluster.id}"
    ]
  }

  // Giving Access to Consul from other ECS Cluster
  ingress {
    from_port = 8300
    to_port = 8300
    protocol = "tcp"
    security_groups = [
      "sg-de5b9fa1" // Staging ECS Cluster
    ]
  }

  // Giving Access to Consul from other ECS Cluster
  ingress {
    from_port = 8301
    to_port = 8301
    protocol = "tcp"
    security_groups = [
      "sg-de5b9fa1" // Staging ECS Cluster
    ]
  }

  // Giving Access to Logstash from other ECS Cluster
  ingress {
    from_port = 12201
    to_port = 12201
    protocol = "tcp"
    security_groups = [
      "sg-de5b9fa1" // Staging ECS Cluster
    ]
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["10.45.114.105/32"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Infra ECS Cluster"
  }
}

resource "aws_security_group" "vault-ecs-cluster" {
  provider    = "aws.local"
  name        = "vault-ecs"
  description = "Vault ECS Cluster"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    self            = true
  }

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_groups = [
      "${aws_security_group.vault-elb.id}"
    ]
  }

  // Authorizing VPN Users to ping Vault nodes directly (required).
  ingress {
    from_port = 8200
    to_port = 8200
    protocol = "tcp"
    security_groups = [
      "${aws_security_group.pritunl-node.id}"
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Vault ECS Cluster"
  }
}

resource "aws_security_group" "consul-elb" {
  provider    = "aws.local"
  name        = "consul-elb"
  description = "Consul ELB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Consul ELB"
  }
}

resource "aws_security_group" "nexus-elb" {
  provider    = "aws.local"
  name        = "nexus-elb"
  description = "Nexus ELB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Nexus ELB"
  }
}

resource "aws_security_group" "vault-elb" {
  provider    = "aws.local"
  name        = "vault-elb"
  description = "Vault ELB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Vault ELB"
  }
}

resource "aws_security_group" "logstash-elb" {
  provider    = "aws.local"
  name        = "logstash-elb"
  description = "Logstash ELB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "Logstash ELB"
  }
}

resource "aws_security_group" "mongo-elb" {
  provider    = "aws.local"
  name        = "mongodb-elb"
  description = "MongoDB ELB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["10.32.0.0/12"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "MongoDB ELB"
  }
}

output "mongodb-elb" {
  value = "${aws_security_group.mongo-elb.id}"
}

output "logstash-elb" {
  value = "${aws_security_group.logstash-elb.id}"
}

output "consul-elb" {
  value = "${aws_security_group.consul-elb.id}"
}

output "vault-elb" {
  value = "${aws_security_group.vault-elb.id}"
}

output "nexus-elb" {
  value = "${aws_security_group.nexus-elb.id}"
}

output "infra-ecs-cluster" {
  value = "${aws_security_group.infra-ecs-cluster.id}"
}

output "vault-ecs-cluster" {
  value = "${aws_security_group.vault-ecs-cluster.id}"
}

output "pritunl_elb" {
  value = "${aws_security_group.pritunl-elb.id}"
}

output "pritunl_node" {
  value = "${aws_security_group.pritunl-node.id}"
}

output "pritunl_mongodb" {
  value = "${aws_security_group.pritunl-mongodb.id}"
}

output "ghe-elb" {
  value = "${aws_security_group.ghe_elb.id}"
}

output "it-ghe-elb" {
  value = "${aws_security_group.it-ghe-elb.id}"
}

output "ghe" {
  value = "${aws_security_group.ghe.id}"
}