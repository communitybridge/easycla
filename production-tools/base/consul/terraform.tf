variable "vpc_id" {
  description = "The ID of the VPC"
}

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

variable "consul_encryption_key" {
  description = "Encryption key to use with Consul"
}

variable "r53_zone_id" {}

variable "region_identifier" {}

variable "region" {}

data "template_file" "consul_slaves_ecs_task" {
  template = "${file("${path.module}/consul-ecs-task.json")}"

  vars {
    CONSUL_ENCRYPTION_KEY   = "${var.consul_encryption_key}"
    CONSUL_DATACENTER       = "AWS"

    EC2_REGION              = "${var.region}"
    EC2_TAG_KEY             = "Name"
    EC2_TAG_VALUE           = "production-tools"
  }
}

resource "aws_ecs_task_definition" "consul-masters" {
  provider    = "aws.local"
  family                = "consul-server"
  container_definitions = "${data.template_file.consul_slaves_ecs_task.rendered}"
  network_mode          = "host"
}

resource "aws_ecs_service" "consul" {
  provider    = "aws.local"
  name                               = "consul-cluster"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.consul-masters.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "consul-cluster" {
  provider    = "aws.local"
  name = "consul-cluster"
  subnets = ["${var.internal_subnets}"]
  security_groups = ["${var.internal_elb_sg}"]
  internal = true

  listener {
    instance_port = 8300
    instance_protocol = "tcp"
    lb_port = 8300
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8301
    instance_protocol = "tcp"
    lb_port = 8301
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8302
    instance_protocol = "tcp"
    lb_port = 8302
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8400
    instance_protocol = "tcp"
    lb_port = 8400
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 8500
    instance_protocol = "http"
    lb_port = 8500
    lb_protocol = "http"
  }

  listener {
    instance_port = 8600
    instance_protocol = "tcp"
    lb_port = 53
    lb_protocol = "tcp"
  }

  health_check {
    target = "TCP:8300"
    healthy_threshold = 2
    unhealthy_threshold = 2
    interval = 30
    timeout = 5
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "Consul Cluster"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "consul-slaves" {
  provider    = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.consul-cluster.id}"
}

resource "aws_route53_record" "consul" {
  zone_id = "${var.r53_zone_id}"
  name    = "consul.${var.region_identifier}"
  type    = "A"

  alias {
    name                   = "${aws_elb.consul-cluster.dns_name}"
    zone_id                = "${aws_elb.consul-cluster.zone_id}"
    evaluate_target_health = true
  }
}

output "consul_elb_cname" {
  value = "${aws_elb.consul-cluster.dns_name}"
}

output "consul_elb_name" {
  value = "${aws_elb.consul-cluster.name}"
}

output "consul_elb_zoneid" {
  value = "${aws_elb.consul-cluster.zone_id}"
}

output "consul_service_name" {
  value = "${aws_ecs_service.consul.name}"
}
