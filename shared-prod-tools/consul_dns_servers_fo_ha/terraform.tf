variable "west_ec2_machines" {
  type = "list"
}

variable "west_ecs_cluster_name" {}

variable "west_ecs_service_name" {}

variable "east_ec2_machines" {
  type = "list"
}

variable "east_ecs_cluster_name" {}

variable "east_ecs_service_name" {}

variable "dns_zone" {}

resource "aws_cloudwatch_metric_alarm" "alarm_west" {
  provider            = "aws.western"
  alarm_name          = "consul-dns-r53-east-health"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "MemoryUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "3"
  alarm_description   = "This ECS Service should never have > 3 unhealthy hosts"

  dimensions {
    ClusterName  = "${var.west_ecs_cluster_name}"
    ServiceName  = "${var.west_ecs_service_name}"
  }
}

resource "aws_cloudwatch_metric_alarm" "alarm_east" {
  provider            = "aws.eastern"
  alarm_name          = "consul-dns-r53-east-health"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "MemoryUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "SampleCount"
  threshold           = "3"
  alarm_description   = "This ECS Service should never have > 3 unhealthy hosts"

  dimensions {
    ClusterName  = "${var.east_ecs_cluster_name}"
    ServiceName  = "${var.east_ecs_service_name}"
  }
}

// a route53 health check on CW
resource "aws_route53_health_check" "healthcheck_west" {
  provider                        = "aws.western"
  type                            = "CLOUDWATCH_METRIC"
  cloudwatch_alarm_name           = "${aws_cloudwatch_metric_alarm.alarm_west.alarm_name}"
  cloudwatch_alarm_region         = "us-west-2"
  insufficient_data_health_status = "Unhealthy"

  tags {
    Name = "consul-dns-servers-west-healthcheck"
  }
}

// a route53 health check on CW
resource "aws_route53_health_check" "healthcheck_east" {
  provider                        = "aws.eastern"
  type                            = "CLOUDWATCH_METRIC"
  cloudwatch_alarm_name           = "${aws_cloudwatch_metric_alarm.alarm_east.alarm_name}"
  cloudwatch_alarm_region         = "us-east-2"
  insufficient_data_health_status = "Unhealthy"

  tags {
    Name = "consul-dns-servers-east-healthcheck"
  }
}

resource "aws_route53_record" "service_west" {
  provider  = "aws.western"
  zone_id   = "${var.dns_zone}"
  name      = "dns"
  type      = "A"
  records   = "${var.west_ec2_machines}"
  ttl       = 30

  failover_routing_policy {
    type = "PRIMARY"
  }

  set_identifier  = "dns-primary"
  health_check_id = "${aws_route53_health_check.healthcheck_west.id}"
}

resource "aws_route53_record" "service_east" {
  provider  = "aws.eastern"
  zone_id   = "${var.dns_zone}"
  name      = "dns"
  type      = "A"
  records   = "${var.east_ec2_machines}"
  ttl       = 30

  failover_routing_policy {
    type = "SECONDARY"
  }

  set_identifier  = "dns-secondary"
  health_check_id = "${aws_route53_health_check.healthcheck_east.id}"
}