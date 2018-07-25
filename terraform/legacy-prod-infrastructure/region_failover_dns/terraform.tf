variable "west_elb_name" {}

variable "west_elb_zoneid" {}

variable "west_elb_dnsname" {}

variable "east_elb_name" {}

variable "east_elb_zoneid" {}

variable "east_elb_dnsname" {}

variable "dns_zone" {}

variable "dns_name" {}

resource "aws_cloudwatch_metric_alarm" "alarm_west" {
  provider            = "aws.western"
  alarm_name          = "${var.dns_name}-r53-east-health"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "UnHealthyHostCount"
  namespace           = "AWS/ELB"
  period              = "60"
  statistic           = "Average"
  threshold           = "0"
  alarm_description   = "This ELB should never have > 0 unhealthy hosts"

  dimensions {
    LoadBalancerName  = "${var.west_elb_name}"
  }
}

resource "aws_cloudwatch_metric_alarm" "alarm_east" {
  provider            = "aws.eastern"
  alarm_name          = "${var.dns_name}-r53-east-health"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "UnHealthyHostCount"
  namespace           = "AWS/ELB"
  period              = "60"
  statistic           = "Average"
  threshold           = "0"
  alarm_description   = "This ELB should never have > 0 unhealthy hosts"

  dimensions {
    LoadBalancerName  = "${var.east_elb_name}"
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
    Name = "${var.dns_name}-west-healthcheck"
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
    Name = "${var.dns_name}-east-healthcheck"
  }
}

resource "aws_route53_record" "service_west" {
  provider  = "aws.western"
  zone_id   = "${var.dns_zone}"
  name      = "${var.dns_name}"
  type      = "A"

  alias {
    name                   = "${var.west_elb_dnsname}"
    zone_id                = "${var.west_elb_zoneid}"
    evaluate_target_health = true
  }

  failover_routing_policy {
    type = "PRIMARY"
  }

  set_identifier  = "${var.dns_name}-primary"
  health_check_id = "${aws_route53_health_check.healthcheck_west.id}"
}

resource "aws_route53_record" "service_east" {
  provider  = "aws.eastern"
  zone_id   = "${var.dns_zone}"
  name      = "${var.dns_name}"
  type      = "A"

  alias {
    name                   = "${var.east_elb_dnsname}"
    zone_id                = "${var.east_elb_zoneid}"
    evaluate_target_health = true
  }

  failover_routing_policy {
    type = "SECONDARY"
  }

  set_identifier  = "${var.dns_name}-secondary"
  health_check_id = "${aws_route53_health_check.healthcheck_east.id}"
}