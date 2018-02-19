variable "ecs_cluster_name" {
  description = "The name of the ECS Cluster"
}

variable "ecs_asg_name" {
  description = "The name of the ECS AutoScaling Group"
}

variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "internal_elb_sg" {
  description = "Security Group for the internal ELB"
}

variable "region" {}

variable "vpc_id" {}

variable "route53_zone_id" {}

resource "aws_s3_bucket" "logstash_bucket" {
  provider = "aws.local"
  bucket   = "logstash-storage-linuxfoundation-org"
  acl      = "private"

  versioning {
    enabled = true
  }
}

resource "aws_iam_user" "logstash" {
  provider = "aws.local"
  name = "logstash-internal"
  path = "/system/"
}

resource "aws_iam_access_key" "logstash" {
  provider = "aws.local"
  user = "${aws_iam_user.logstash.name}"
}

resource "aws_iam_user_policy" "logstash_ro" {
  provider = "aws.local"
  name = "logstash-bucket-permissions"
  user = "${aws_iam_user.logstash.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": "${aws_s3_bucket.logstash_bucket.arn}/*"
    }
  ]
}
EOF
}

data "template_file" "logstash_ecs_task" {
  template = "${file("${path.module}/ecs-task-def.json")}"

  vars {
    AWS_ACCESS_KEY_ID     = "${aws_iam_access_key.logstash.id}"
    AWS_SECRET_ACCESS_KEY = "${aws_iam_access_key.logstash.secret}"
    S3_LOG_BUCKET         = "${aws_s3_bucket.logstash_bucket.bucket}"
    ES_USERNAME           = "elastic"
    ES_PASSWORD           = "lXUOIVmZyeKkoa2sFNejH7K3"
    ES_HOST_URL           = "https://697e30844ccaa3c63d1f99a6095faff0.us-west-2.aws.found.io:9243"
    PIPELINE_NAME         = "production"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "logstash" {
  provider = "aws.local"
  family = "logstash"
  network_mode = "host"

  lifecycle {
    ignore_changes        = ["image"]
    create_before_destroy = true
  }

  container_definitions = "${data.template_file.logstash_ecs_task.rendered}"
}

resource "aws_ecs_service" "logstash" {
  provider                           = "aws.local"
  name                               = "logstash"
  cluster                            = "${var.ecs_cluster_name}"
  task_definition                    = "${aws_ecs_task_definition.logstash.arn}"
  desired_count                      = "3"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Create a new load balancer
resource "aws_elb" "logstash" {
  provider           = "aws.local"
  name               = "logstash"
  security_groups    = ["${var.internal_elb_sg}"]
  subnets            = ["${var.internal_subnets}"]
  internal           = true

  listener {
    instance_port     = 12201
    instance_protocol = "tcp"
    lb_port           = 12201
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 1514
    instance_protocol = "tcp"
    lb_port           = 1514
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 5044
    instance_protocol = "tcp"
    lb_port           = 5044
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "TCP:12201"
    interval            = 30
  }

  cross_zone_load_balancing   = true
  idle_timeout                = 400

  tags {
    Name = "logstash"
  }
}

# Create a new load balancer attachment
resource "aws_autoscaling_attachment" "logstash" {
  provider = "aws.local"
  autoscaling_group_name = "${var.ecs_asg_name}"
  elb                    = "${aws_elb.logstash.id}"
}

resource "aws_route53_record" "logstash" {
  provider = "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "logstash"
  type    = "A"

  alias {
    name                   = "${aws_elb.logstash.dns_name}"
    zone_id                = "${aws_elb.logstash.zone_id}"
    evaluate_target_health = true
  }
}
