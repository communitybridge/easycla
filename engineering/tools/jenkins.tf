resource "aws_ecs_task_definition" "jenkins_master_task" {
  family = "jenkins"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  volume {
    name = "tools-storage"
    host_path = "/mnt/storage/ci/jenkins"
  }

  container_definitions = <<EOF
[
  {
    "essential": true,
    "image": "jenkinsci/jenkins",
    "memoryReservation": 1024,
    "name": "jenkins-master",
    "portMappings": [{
          "hostPort": 8080,
          "containerPort": 8080
        },
        {
          "hostPort": 50000,
          "containerPort": 50000
        }],
    "mountPoints": [{
            "sourceVolume": "tools-storage",
            "containerPath": "/var/jenkins_home"
          }],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "engineering-tools",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "jenkins"
      }
    }
  }
]
EOF

}

resource "aws_ecs_service" "jenkins_master_service" {
  name                               = "jenkins-master"
  cluster                            = "${module.engineering-tools-ecs-cluster.name}"
  task_definition                    = "${aws_ecs_task_definition.jenkins_master_task.arn}"
  desired_count                      = "1"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Security Group for EFS
resource "aws_security_group" "sg_external_jenkins_elb" {
  name        = "engineering-tools-jenkins-elb"
  description = "Allows External Access to Engineering Jenkins Instance"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

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
    Name        = "engineering-tools-jenkins-elb"
  }
}

# Create a new load balancer
resource "aws_elb" "engineering_tools_jenkins" {
  name = "engineering-tools-jenkins"
  subnets = ["${data.terraform_remote_state.engineering.internal_subnets}"]
  security_groups = ["${aws_security_group.sg_external_jenkins_elb.id}"]
  internal = true

  listener {
    instance_port = 8080
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  listener {
    instance_port = 8080
    instance_protocol = "http"
    lb_port = 443
    lb_protocol = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:433610389961:certificate/bb946be4-a4f4-4f91-a786-60eddbd055b6"
  }

  listener {
    instance_port = 50000
    instance_protocol = "tcp"
    lb_port = 50000
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:50000"
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
resource "aws_autoscaling_attachment" "jenkins_to_asg" {
  autoscaling_group_name = "${module.engineering-tools-ecs-cluster.asg_name}"
  elb                    = "${aws_elb.engineering_tools_jenkins.id}"
}

resource "aws_route53_record" "jenkins" {
  zone_id = "Z2MDT77FL23F9B"
  name = "jenkins.engineering.tux.rocks."
  type = "A"

  alias {
    name = "${aws_elb.engineering_tools_jenkins.dns_name}"
    zone_id = "${aws_elb.engineering_tools_jenkins.zone_id}"
    evaluate_target_health = true
  }
}

# Security Group for EFS
resource "aws_security_group" "jenkins-slave" {
  name        = "jenkins-slave"
  description = "Jenkins Slaves"
  vpc_id      = "${data.terraform_remote_state.engineering.vpc_id}"

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = ["${module.engineering-tools-ecs-cluster.security_group_id}"]
  }

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = ["0.0.0.0/0"]
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
    Name        = "jenkins-slave"
  }
}

resource "aws_security_group_rule" "jenkins" {
  type = "ingress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  source_security_group_id = "${aws_security_group.sg_external_jenkins_elb.id}"

  security_group_id = "${aws_security_group.tools.id}"
}
