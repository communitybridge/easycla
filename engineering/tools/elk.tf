resource "aws_ecs_task_definition" "elk" {
  family = "elk"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  volume {
    name = "tools-storage"
    host_path = "/mnt/storage/elk"
  }

  container_definitions = <<EOF
[
  {
    "essential": true,
    "image": "sebp/elk:521",
    "memoryReservation": 1024,
    "name": "elk",
    "portMappings": [{
          "hostPort": 5601,
          "containerPort": 5601
        },
        {
          "hostPort": 9200,
          "containerPort": 9200
        },
        {
          "hostPort": 5044,
          "containerPort": 5044
        }],
    "mountPoints": [{
            "sourceVolume": "tools-storage",
            "containerPath": "/var/lib/elasticsearch"
          }],
    "ulimits": [{
      "name": "nofile",
      "softLimit": 65536,
      "hardLimit": 65536
    }],
    "environment": [
      { "name" : "ES_JAVA_OPTS", "value" : "-Xmx2g -Xms2g" }
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "engineering-tools",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "elk"
      }
    }
  }
]
EOF

}

resource "aws_ecs_service" "elk" {
  name                               = "elk"
  cluster                            = "${module.engineering-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.elk.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }

}

# Security Group for EFS
resource "aws_security_group" "sg_external_elk_elb" {
  name        = "engineering-tools-elk-elb"
  description = "Allows External Access to Engineering ELK Instance"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  ingress {
    from_port       = 9200
    to_port         = 9200
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 5044
    to_port         = 5044
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 443
    to_port         = 443
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
    Name        = "engineering-tools-elk-elb"
  }
}

# Create a new load balancer
resource "aws_elb" "elk" {
  name = "engineering-tools-elk"
  subnets = ["${data.terraform_remote_state.engineering.internal_subnets}"]
  internal = true
  security_groups = ["${aws_security_group.sg_external_elk_elb.id}"]

  listener {
    instance_port = 5601
    instance_protocol = "http"
    lb_port = 443
    lb_protocol = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
  }

  listener {
    instance_port = 5601
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  listener {
    instance_port = 9200
    instance_protocol = "tcp"
    lb_port = 9200
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 5044
    instance_protocol = "tcp"
    lb_port = 5044
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:9200"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400

  tags {
    Name = "engineering-tools-elk"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "elk" {
  autoscaling_group_name = "${module.engineering-tools-ecs-cluster.asg_name}"
  elb                    = "${aws_elb.elk.id}"
}

resource "aws_route53_record" "elk" {
  zone_id = "Z2MDT77FL23F9B"
  name = "elk.engineering.tux.rocks."
  type = "A"

  alias {
    name = "${aws_elb.elk.dns_name}"
    zone_id = "${aws_elb.elk.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_security_group_rule" "elk" {
  type = "ingress"
  from_port = -1
  to_port = -1
  protocol = "-1"
  source_security_group_id = "${aws_security_group.sg_external_elk_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}