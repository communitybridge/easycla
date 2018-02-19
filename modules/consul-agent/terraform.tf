variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "encryption_key" {
  description = "The Consul Encryption Key"
}

variable "datacenter" {
  description = "The Consul Datacenter we are connecting to"
}

variable "endpoint" {
  description = "Consul Server Endpoint the Agent needs to connect to"
}


data "template_file" "consul_ecs_task" {
  template = "${file("${path.module}/consul-ecs-task.json")}"

  vars {
    # Consul Info
    ENCRYPTION    = "${var.encryption_key}"
    DATACENTER    = "${var.datacenter}"
    CONSUL_ENDPOINT = "${var.endpoint}"
  }
}

resource "aws_ecs_task_definition" "consul" {
  provider              = "aws.local"
  family                = "consul-agent"
  container_definitions = "${data.template_file.consul_ecs_task.rendered}"
  network_mode          = "host"
}

resource "aws_ecs_service" "consul" {
  provider                           = "aws.local"
  name                               = "consul-agent"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.consul.arn}"
  desired_count                      = "2"
  deployment_maximum_percent         = "200"

  placement_strategy {
    type = "spread"
    field = "instanceId"
  }

  placement_constraints {
    type = "distinctInstance"
  }
}
