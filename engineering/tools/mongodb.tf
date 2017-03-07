# MongoDB Server (for lfweb-cli)
resource "aws_ecs_task_definition" "task" {
  family = "cli_mongodb"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  volume {
    name = "tools-storage"
    host_path = "/mnt/storage/cli/mongodb"
  }

  container_definitions = <<EOF
[
  {
    "essential": true,
    "image": "433610389961.dkr.ecr.us-west-2.amazonaws.com/tools/mongodb:latest",
    "memoryReservation": 512,
    "name": "mongodb",
    "portMappings": [{
          "containerPort": 59215,
          "hostPort": 59215
        }],
    "mountPoints": [{
            "sourceVolume": "tools-storage",
            "containerPath": "/mnt/storage/cli/mongodb"
          }],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "engineering-tools",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "mongo"
      }
    }
  }
]
EOF


}

resource "aws_ecs_service" "main" {
  name                               = "cli-mongodb"
  cluster                            = "${module.engineering-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.task.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Security Group for EFS
resource "aws_security_group" "sg_external_mongo_elb" {
  name        = "engineering-tools-mongo-elb"
  description = "Allows External Access to Engineering MongoDB Instance"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  ingress {
    from_port       = 59215
    to_port         = 59215
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  egress {
    from_port       = 59215
    to_port         = 59215
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Name        = "engineering-tools-mongo-elb"
  }
}

# Create a new load balancer
resource "aws_elb" "mongodb" {
  name = "engineering-tools-mongodb"
  subnets = ["${data.terraform_remote_state.engineering.internal_subnets}"]
  security_groups = ["${aws_security_group.sg_external_mongo_elb.id}"]
  internal = true

  listener {
    instance_port = 59215
    instance_protocol = "ssl"
    lb_port = 59215
    lb_protocol = "ssl"
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:59215"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "engineering-tools-mongodb"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "mongodb" {
  autoscaling_group_name = "${module.engineering-tools-ecs-cluster.asg_name}"
  elb                    = "${aws_elb.mongodb.id}"
}

resource "aws_route53_record" "mongodb" {
  zone_id = "Z2MDT77FL23F9B"
  name = "mongodb.engineering.tux.rocks."
  type = "A"

  alias {
    name = "${aws_elb.mongodb.dns_name}"
    zone_id = "${aws_elb.mongodb.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_security_group_rule" "mongodb" {
  type = "ingress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  source_security_group_id = "${aws_security_group.sg_external_mongo_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}