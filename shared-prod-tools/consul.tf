data "aws_ami" "amazon-linux" {
  most_recent = true

  filter {
    name = "name"
    values = ["amzn-ami-*-x86_64-gp2"]
  }
  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }
  filter {
    name = "owner-alias"
    values = ["amazon"]
  }
}


data "template_file" "consul_slaves_ecs_task" {
  template = "${file("${path.module}/files/consul-ecs-task.json")}"

  vars {
    EC2_NAME_TAG            = "${module.shared-production-tools-ecs-cluster.name}"
    CONSUL_ENCRYPTION_KEY   = "kNjj4yLEimIH6GvasrhSIg=="
    CONSUL_CLUSTER_MIN_SIZE = "3"
    CONSUL_DATACENTER       = "AWS"
  }
}

resource "aws_ecs_task_definition" "consul-masters" {
  family                = "consul-server"
  container_definitions = "${data.template_file.consul_slaves_ecs_task.rendered}"
  network_mode          = "host"
}

resource "aws_ecs_service" "consul" {
  name                               = "consul"
  cluster                            = "${module.shared-production-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.consul-masters.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "consul-master" {
  name = "consul-master"
  subnets = ["${module.vpc.internal_subnets}"]
  security_groups = ["${module.security_groups.internal_elb}"]
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
    lb_port = 8600
    lb_protocol = "tcp"
  }

  health_check {
    target = "HTTP:8500/v1/status/leader"
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
  autoscaling_group_name = "${module.shared-production-tools-ecs-cluster.asg_name}"
  elb                    = "${aws_elb.consul-master.id}"
}

resource "aws_route53_record" "consul" {
  zone_id = "${module.dns.zone_id}"
  name = "consul.prod.engineering.internal."
  type = "A"

  alias {
    name = "${aws_elb.consul-master.dns_name}"
    zone_id = "${aws_elb.consul-master.zone_id}"
    evaluate_target_health = true
  }
}
