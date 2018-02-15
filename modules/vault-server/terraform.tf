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

variable "route53_zone_id" {}

variable "data_path" {}

data "template_file" "vault_ecs_task" {
  template = "${file("${path.module}/ecs-task-def.json")}"

  vars {
    # Used for Docker Tags
    AWS_REGION      = "${var.region}"
  }
}

resource "aws_ecs_task_definition" "vault" {
  provider              = "aws.local"
  family                = "vault"
  container_definitions = "${data.template_file.vault_ecs_task.rendered}"
  network_mode          = "host"

  volume {
    name      = "vault-data-storage"
    host_path = "${var.data_path}/data"
  }

  volume {
    name      = "vault-logs-storage"
    host_path = "${var.data_path}/logs"
  }
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


# Create a new load balancer
resource "aws_elb" "vault" {
  provider           = "aws.local"
  name               = "vault-server"
  security_groups    = ["${var.internal_elb_sg}"]
  subnets            = ["${var.internal_subnets}"]
  internal           = true

  listener {
    instance_port     = 8200
    instance_protocol = "tcp"
    lb_port           = 8200
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 8201
    instance_protocol = "tcp"
    lb_port           = 8201
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "TCP:8200"
    interval            = 30
  }

  cross_zone_load_balancing   = true
  idle_timeout                = 400

  tags {
    Name = "vault"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "vault" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.vault.id}"
}

resource "aws_route53_record" "vault" {
  provider = "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "vault"
  type    = "A"

  alias {
    name                   = "${aws_elb.vault.dns_name}"
    zone_id                = "${aws_elb.vault.zone_id}"
    evaluate_target_health = true
  }
}
