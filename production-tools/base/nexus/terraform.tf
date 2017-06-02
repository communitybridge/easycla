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

variable "building" {}

data "template_file" "nexus_ecs_task" {
  count    = "${var.building}"
  template = "${file("${path.module}/nexus-ecs-task.json")}"

  vars {
    # DNS Servers for Container Resolution
    DNS_SERVER_1          = "${var.dns_servers[0]}"
    DNS_SERVER_2          = "${var.dns_servers[1]}"
    DNS_SERVER_3          = "${var.dns_servers[2]}"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "nexus" {
  provider = "aws.local"
  count    = "${var.building}"
  family   = "nexus"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.nexus_ecs_task.rendered}"

  volume {
    name      = "nexus-storage"
    host_path = "/mnt/storage/nexus"
  }
}

resource "aws_ecs_service" "nexus" {
  count    = "${var.building}"
  provider                           = "aws.local"
  name                               = "nexus-repository"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.nexus.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "nexus" {
  count    = "${var.building}"
  provider = "aws.local"
  name = "nexus"
  subnets = ["${var.internal_subnets}"]
  security_groups = ["${var.internal_elb_sg}"]
  internal = true

  listener {
    instance_port = 8081
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  listener {
    instance_port = 8081
    instance_protocol = "http"
    lb_port = 443
    lb_protocol = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:8081"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "nexus-cluster"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "nexus" {
  count    = "${var.building}"
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.nexus.id}"
}

resource "aws_route53_record" "consul" {
  count    = "${var.building}"
  zone_id = "Z2MDT77FL23F9B"
  name    = "nexus"
  type    = "A"

  alias {
    name                   = "${aws_elb.nexus.dns_name}"
    zone_id                = "${aws_elb.nexus.zone_id}"
    evaluate_target_health = true
  }
}