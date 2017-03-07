//resource "aws_s3_bucket" "b" {
//  bucket = "lf-engineering-pypi-repository"
//  acl = "private"
//
//  tags {
//    Name = "Pypi Repository"
//    Environment = "Production"
//  }
//}

resource "aws_ecs_task_definition" "pypi" {
  family = "pypi"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = <<EOF
[
  {
    "essential": true,
    "image": "433610389961.dkr.ecr.us-west-2.amazonaws.com/tools/pypi:latest",
    "memoryReservation": 512,
    "name": "pypi",
    "portMappings": [{
          "containerPort": 8080
        }],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "engineering-tools",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "pypi"
      }
    }
  }
]
EOF

}

resource "aws_ecs_service" "pypi" {
  name                               = "pypi-server"
  cluster                            = "${module.engineering-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.pypi.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"
  iam_role                           = "arn:aws:iam::433610389961:role/ecsServiceRole"

  lifecycle {
    create_before_destroy = true
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.pypi.arn}"
    container_name = "pypi"
    container_port = 8080
  }
}

# Security Group for EFS
resource "aws_security_group" "sg_external_pypi_elb" {
  name        = "engineering-tools-pypi-elb"
  description = "Allows External Access to Engineering Pypi Instance"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "engineering-tools-pypi-elb"
  }
}

# Create a new load balancer
resource "aws_alb" "pypi" {
  name            = "engineering-tools-pypi"
  internal        = true
  security_groups = ["${aws_security_group.sg_external_pypi_elb.id}"]
  subnets = ["${data.terraform_remote_state.engineering.internal_subnets}"]

  tags {
    Environment = "Production",
    Team        = "Engineering",
    Name        = "PyPi Server"
  }
}

resource "aws_alb_target_group" "pypi" {
  name     = "engineering-tools-pypi"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${data.terraform_remote_state.engineering.vpc_id}"
}

resource "aws_alb_listener" "pypi_80" {
  load_balancer_arn = "${aws_alb.pypi.arn}"
  port = "80"
  protocol = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.pypi.arn}"
    type = "forward"
  }
}

resource "aws_alb_listener" "pypi_443" {
  load_balancer_arn = "${aws_alb.pypi.arn}"
  port = "443"
  protocol = "HTTPS"
  ssl_policy = "ELBSecurityPolicy-2015-05"
  certificate_arn = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"

  default_action {
    target_group_arn = "${aws_alb_target_group.pypi.arn}"
    type = "forward"
  }
}

resource "aws_route53_record" "pypi" {
  zone_id = "Z2MDT77FL23F9B"
  name = "pypi.engineering.tux.rocks."
  type = "A"

  alias {
    name = "${aws_alb.pypi.dns_name}"
    zone_id = "${aws_alb.pypi.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_security_group_rule" "pypi" {
  type = "ingress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  source_security_group_id = "${aws_security_group.sg_external_pypi_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}