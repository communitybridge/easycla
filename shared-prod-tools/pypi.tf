resource "aws_s3_bucket" "pypi_repo" {
  bucket = "lf-engineering-python-repository"
  acl = "private"

  tags {
    Name = "Pypi Repository"
    Environment = "Production"
  }
}

data "template_file" "pypi_ecs_task" {
  template = "${file("${path.module}/files/pypi-ecs-task.json")}"

  vars {
    # Those are coming from the pypi-user on AWS, needs Read/Write on above S3 Bucket & FullDynamoDB
    AWS_ACCESS_KEY_ID = "AKIAIB6UWB7QG5QQYPWQ"
    AWS_SECRET_ACCESS_KEY = "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"
  }
}

resource "aws_ecs_task_definition" "pypi" {
  family = "pypi-repo"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.pypi_ecs_task.rendered}"
}

resource "aws_ecs_service" "pypi" {
  name                               = "pypi-repository"
  cluster                            = "${module.shared-production-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.pypi.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "pypi" {
  name = "pypi"
  subnets = ["${module.vpc.internal_subnets}"]
  security_groups = ["${module.security_groups.internal_elb}"]
  internal = true

  listener {
    instance_port = 8080
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:8080"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "pypi"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "pypi" {
  autoscaling_group_name = "${module.shared-production-tools-ecs-cluster.asg_name}"
  elb                    = "${aws_elb.pypi.id}"
}

resource "aws_route53_record" "pypi" {
  zone_id = "${module.dns.zone_id}"
  name = "pypi.prod.engineering.internal."
  type = "A"

  alias {
    name = "${aws_elb.pypi.dns_name}"
    zone_id = "${aws_elb.pypi.zone_id}"
    evaluate_target_health = true
  }
}

module "redis-cluster" "pypi-cache" {
  source          = "../modules/redis-cluster"
  name            = "pypi-server-cache"
  version         = "3.2.4"
  instance_type   = "cache.t2.micro"
  instance_count  = "1"
  environment     = "Production"
  team            = "Engineering"
  security_groups = ["${module.shared-production-tools-ecs-cluster.security_group_id}"]
  subnet_ids      = "${module.vpc.internal_subnets}"
  vpc_id          = "${module.vpc.id}"
}