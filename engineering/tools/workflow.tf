resource "aws_ecs_task_definition" "workflow_task" {
  family = "workflow"

  lifecycle {
    create_before_destroy = true
  }

  container_definitions = <<EOF
[
  {
    "essential": true,
    "image": "433610389961.dkr.ecr.us-west-2.amazonaws.com/tools/workflow:latest",
    "memoryReservation": 256,
    "name": "workflow",
    "portMappings": [
        {
          "containerPort": 5555
        }
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "engineering-tools",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "workflow"
      }
    }
  }
]
EOF

}

resource "aws_ecs_service" "workflow_service" {
  name                               = "workflow"
  cluster                            = "${module.engineering-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.workflow_task.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"
  iam_role                           = "arn:aws:iam::433610389961:role/ecsServiceRole"

  lifecycle {
    create_before_destroy = true
  }

  load_balancer {
    target_group_arn = "${aws_alb_target_group.workflow.arn}"
    container_name = "workflow"
    container_port = 5555
  }
}

# Security Group for EFS
resource "aws_security_group" "sg_external_workflow_elb" {
  name        = "engineering-tools-workflow-elb"
  description = "Allows External Access to Engineering Workflow API"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  ingress {
    from_port       = 5555
    to_port         = 5555
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
    Name        = "engineering-tools-workflow-elb"
  }
}

# Create a new load balancer
resource "aws_alb" "workflow" {
  name            = "engineering-tools-workflow"
  security_groups = ["${aws_security_group.sg_external_workflow_elb.id}"]
  subnets = ["${data.terraform_remote_state.engineering.internal_subnets}"]
  internal = true

  tags {
    Environment = "Production",
    Team        = "Engineering",
    Name        = "Workflow API"
  }
}

resource "aws_alb_target_group" "workflow" {
  name     = "engineering-tools-workflow"
  port     = 5555
  protocol = "HTTP"
  vpc_id   = "${data.terraform_remote_state.engineering.vpc_id}"
}

resource "aws_alb_listener" "workflow_80" {
  load_balancer_arn = "${aws_alb.workflow.arn}"
  port = "80"
  protocol = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.workflow.arn}"
    type = "forward"
  }
}

resource "aws_alb_listener" "workflow_443" {
  load_balancer_arn = "${aws_alb.workflow.arn}"
  port = "443"
  protocol = "HTTPS"
  ssl_policy = "ELBSecurityPolicy-2015-05"
  certificate_arn = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"

  default_action {
    target_group_arn = "${aws_alb_target_group.workflow.arn}"
    type = "forward"
  }
}

resource "aws_route53_record" "workflow" {
  zone_id = "Z2MDT77FL23F9B"
  name = "workflow.engineering.tux.rocks."
  type = "A"

  alias {
    name = "${aws_alb.workflow.dns_name}"
    zone_id = "${aws_alb.workflow.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_security_group_rule" "workflow" {
  type = "ingress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  source_security_group_id = "${aws_security_group.sg_external_workflow_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}