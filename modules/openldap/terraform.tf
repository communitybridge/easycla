variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "ecs_asg_name" {
  description = "The name of the ECS AutoScaling Group"
}

variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "internal_elb_sg" {
  description = "Security Group for the internal ELB"
}

variable "region" {}

variable "vpc_id" {}

variable "dns_servers" {
  type = "list"
}

variable "ldap_org" {}

variable "ldap_domain" {}

variable "ldap_admin_password" {}

data "template_file" "openldap_ecs_task" {
  template = "${file("${path.module}/openldap-ecs-task.json")}"

  vars {
    # OpenLDAP
    LDAP_ORGANISATION     = "${var.ldap_org}"
    LDAP_DOMAIN           = "${var.ldap_domain}"
    LDAP_ADMIN_PASSWORD   = "${var.ldap_admin_password}"

    # DNS Servers for Container Resolution
    DNS_SERVER_1          = "${var.dns_servers[0]}"
    DNS_SERVER_2          = "${var.dns_servers[1]}"
    DNS_SERVER_3          = "${var.dns_servers[2]}"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "openldap" {
  provider = "aws.local"
  family   = "openldap"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.openldap_ecs_task.rendered}"
}

resource "aws_ecs_service" "openldap" {
  provider                           = "aws.local"
  name                               = "openldap"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.openldap.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "openldap" {
  provider = "aws.local"
  name = "openldap"
  subnets = ["${var.internal_subnets}"]
  security_groups = ["${var.internal_elb_sg}"]
  internal = true

  listener {
    instance_port = 389
    instance_protocol = "tcp"
    lb_port = 389
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 636
    instance_protocol = "tcp"
    lb_port = 636
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:389"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "openldap"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "openldap" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.openldap.id}"
}

resource "aws_route53_record" "openldap" {
  provider    = "aws.local"
  zone_id = "Z2MDT77FL23F9B"
  name    = "sandbox-ldap"
  type    = "A"

  alias {
    name                   = "${aws_elb.openldap.dns_name}"
    zone_id                = "${aws_elb.openldap.zone_id}"
    evaluate_target_health = true
  }
}