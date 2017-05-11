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

data "template_file" "logstash_ecs_task" {
  template = "${file("${path.module}/logstash-ecs-task.json")}"

  vars {
    # Those are coming from the logstash-user on AWS, needs Read/Write on above S3 Bucket
    AWS_ACCESS_KEY_ID     = "AKIAIQEMUX46ZEUE4TTQ"
    AWS_SECRET_ACCESS_KEY = "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"
    S3_LOG_BUCKET         = "lf-engineering-production-logs"
    ES_USERNAME           = "elastic"
    ES_PASSWORD           = "zShz2Q3vamuhtIhZ3FkV55Py"
    ES_HOST_URL           = "https://c99ea3c2e8e8a90dd2fed3e9564a4c1e.us-west-2.aws.found.io:9243"
    PIPELINE_NAME         = "production_tools"

    # Tags for Registrator
    TAG_REGION            = "${var.region}"
    TAG_VPC_ID            = "${var.vpc_id}"
  }
}

resource "aws_ecs_task_definition" "logstash" {
  provider = "aws.local"
  family = "logstash"

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
