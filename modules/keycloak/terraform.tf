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

# Create a new load balancer
resource "aws_elb" "keycloak" {
  provider = "aws.local"
  name = "keycloak-${var.env}"
  subnets = ["${var.subnets}"]
  security_groups = ["${var.internal_elb_sg}"]

  listener {
    instance_port = 8576
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  listener {
    instance_port = 8576
    instance_protocol = "http"
    lb_port = 443
    lb_protocol = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
  }

  listener {
    instance_port = 7600
    instance_protocol = "TCP"
    lb_port = 7600
    lb_protocol = "TCP"
  }

  listener {
    instance_port = 57600
    instance_protocol = "TCP"
    lb_port = 57600
    lb_protocol = "TCP"
  }

  listener {
    instance_port = 57601
    instance_protocol = "TCP"
    lb_port = 57601
    lb_protocol = "TCP"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:8576"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "keycloak-${var.env}"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "logstash" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.keycloak.id}"
}

resource "aws_route53_record" "keycloak" {
  provider  = "aws.local"
  zone_id   = "Z2MDT77FL23F9B"
  name      = "keycloak"
  type      = "A"

  alias {
    name                   = "${aws_elb.keycloak.dns_name}"
    zone_id                = "${aws_elb.keycloak.zone_id}"
    evaluate_target_health = true
  }
}