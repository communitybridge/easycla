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

variable "consul_endpoint" {
  description = "Endpoint to hit in order to connect to Consul"
}

variable "dns_servers" {
  type = "list"
}

variable "region" {}

data "template_file" "vault_ecs_task" {
  template = "${file("${path.module}/vault-ecs-task.json")}"

  vars {
    CONSUL_MASTER_IP = "${var.consul_endpoint}"

    # DNS Servers for Container Resolution
    DNS_SERVER_1     = "${var.dns_servers[0]}"
    DNS_SERVER_2     = "${var.dns_servers[1]}"
    DNS_SERVER_3     = "${var.dns_servers[2]}"

    # Used for Docker Tags
    AWS_REGION      = "${var.region}"
  }
}

resource "aws_ecs_task_definition" "vault" {
  provider              = "aws.local"
  family                = "vault"
  container_definitions = "${data.template_file.vault_ecs_task.rendered}"
  network_mode          = "host"
}

resource "aws_ecs_service" "vault" {
  provider                           = "aws.local"
  name                               = "vault"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.vault.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}
