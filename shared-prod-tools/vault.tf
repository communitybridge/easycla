//data "template_file" "vault_ecs_task" {
//  template = "${file("${path.module}/files/vault-ecs-task.json")}"
//
//  vars {
//    CONSUL_MASTER_IP = "consul.prod.engineering.internal"
//    VPC_DNS_SERVER   = "10.50.0.2"
//  }
//}
//
//resource "aws_ecs_task_definition" "vault" {
//  family                = "vault"
//  container_definitions = "${data.template_file.vault_ecs_task.rendered}"
//  network_mode          = "host"
//}
//
//resource "aws_ecs_service" "vault" {
//  name                               = "vault"
//  cluster                            = "${module.shared-production-tools-ecs-cluster.name}"
//  task_definition                    = "${aws_ecs_task_definition.vault.arn}"
//  desired_count                      = "3"
//  deployment_minimum_healthy_percent = "100"
//  deployment_maximum_percent         = "200"
//
//  lifecycle {
//    create_before_destroy = true
//  }
//}
//
//# Create a new load balancer
//resource "aws_elb" "vault" {
//  name = "vault"
//  subnets = ["${module.vpc.internal_subnets}"]
//  security_groups = ["${module.security_groups.internal_elb}"]
//  internal = true
//
//  listener {
//    instance_port = 8200
//    instance_protocol = "http"
//    lb_port = 8200
//    lb_protocol = "http"
//  }
//
//  health_check {
//    healthy_threshold = 2
//    unhealthy_threshold = 2
//    timeout = 3
//    target = "TCP:8200"
//    interval = 30
//  }
//
//  cross_zone_load_balancing = true
//  idle_timeout = 400
//  connection_draining = true
//  connection_draining_timeout = 400
//
//  tags {
//    Name = "vault"
//  }
//}
//
//# Create a new load balancer attachment
//resource "aws_autoscaling_attachment" "vault" {
//  autoscaling_group_name = "${module.shared-production-tools-ecs-cluster.asg_name}"
//  elb                    = "${aws_elb.vault.id}"
//}
//
//resource "aws_route53_record" "vault" {
//  zone_id = "${module.dns.zone_id}"
//  name = "vault.prod.engineering.internal."
//  type = "A"
//
//  alias {
//    name = "${aws_elb.vault.dns_name}"
//    zone_id = "${aws_elb.vault.zone_id}"
//    evaluate_target_health = true
//  }
//}
