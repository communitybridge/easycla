variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "dns_servers" {
  type = "list"
}

variable "region" {
  description = "The AWS Region in which we are providing the service in"
  default = "us-west-2"
}

variable "encryption_key" {
  description = "The Consul Encryption Key"
}

variable "datacenter" {
  description = "The Consul Datacenter we are connecting to"
}

variable "endpoint" {}

data "template_file" "consul_ecs_task" {
  template = "${file("${path.module}/consul-ecs-task.json")}"

  vars {
    # DNS Servers for Container Resolution
    DNS_SERVER_1  = "${var.dns_servers[0]}"
    DNS_SERVER_2  = "${var.dns_servers[1]}"
    DNS_SERVER_3  = "${var.dns_servers[2]}"

    # Consul Info
    ENCRYPTION    = "${var.encryption_key}"
    DATACENTER    = "${var.datacenter}"
    CONSUL_ENDPOINT = "${var.endpoint}"

    # For Logging
    AWS_REGION    = "${var.region}"
  }
}

resource "aws_ecs_task_definition" "consul" {
  provider              = "aws.local"
  family                = "consul-agent"
  container_definitions = "${data.template_file.consul_ecs_task.rendered}"
  network_mode          = "host"
}

resource "aws_ecs_service" "registrator" {
  provider                           = "aws.local"
  name                               = "consul-agent"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.consul.arn}"
  desired_count                      = "12"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  placement_strategy {
    type = "spread"
    field = "instanceId"
  }

  placement_constraints {
    type = "distinctInstance"
  }
}
