variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "ecs_asg_name" {
  description = "The name of the ECS AutoScaling Group"
}

variable "subnets" {
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

variable "env" {}

variable "mysql_host" {}

variable "mysql_port" {}

variable "mysql_user" {}

variable "mysql_pass" {}

variable "mysql_db" {}

data "template_file" "keycloak_ecs_task" {
  template = "${file("${path.module}/keycloak-ecs-task.json")}"

  vars {
    # Keycloak env type
    KEYCLOAK_ENV          = "${var.env}"
    MYSQL_HOST            = "${var.mysql_host}"
    MYSQL_PORT            = "${var.mysql_port}"
    MYSQL_USER            = "${var.mysql_user}"
    MYSQL_PASS            = "${var.mysql_pass}"
    MYSQL_DATABASE        = "${var.mysql_db}"

    KEYCLOAK_USER         = "lfengineering"
    KEYCLOAK_PASSWORD     = "YFUUg8omAVjyEB"

    # DNS Servers for Container Resolution
    DNS_SERVER_1          = "${var.dns_servers[0]}"
    DNS_SERVER_2          = "${var.dns_servers[1]}"
    DNS_SERVER_3          = "${var.dns_servers[2]}"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "keycloak" {
  provider = "aws.local"
  family   = "keycloak-${var.env}"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.keycloak_ecs_task.rendered}"
}

resource "aws_ecs_service" "keycloak" {
  provider                           = "aws.local"
  name                               = "keycloak-${var.env}"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.keycloak.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"
}
