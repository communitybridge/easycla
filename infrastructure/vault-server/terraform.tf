variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "internal_elb_sg" {
  description = "Security Group for the internal ELB"
}

variable "ecs_sg" {}

variable "region" {}

variable "route53_zone_id" {}

variable "data_path" {}

variable "efs_id" {}

variable "iam_role" {}

data "template_file" "vault_ecs_task" {
  template = "${file("${path.module}/ecs-task-def.json")}"

  vars {
    # Used for Docker Tags
    AWS_REGION      = "${var.region}"
  }
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "vault"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
    aws_region       = "us-west-2"
    vault_efs      = "${var.efs_id}"
  }
}

data "aws_ami" "amazon-linux-ecs-optimized" {
  provider    = "aws.local"
  most_recent = true

  filter {
    name = "name"
    values = ["amzn-ami-*-amazon-ecs-optimized"]
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

resource "aws_ecs_cluster" "vault" {
  provider = "aws.local"
  name = "vault"
}

resource "aws_instance" "vault-a" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.amazon-linux-ecs-optimized.id}"
  vpc_security_group_ids = ["${var.ecs_sg}"]
  subnet_id              = "${var.internal_subnets[0]}"
  instance_type          = "m4.large"
  key_name               = "engineering-production"
  user_data              = "${data.template_file.ecs_cloud_config.rendered}"
  iam_instance_profile   = "${var.iam_role}"

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "vault-a"
  }
}

resource "aws_instance" "vault-b" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.amazon-linux-ecs-optimized.id}"
  vpc_security_group_ids = ["${var.ecs_sg}"]
  subnet_id              = "${var.internal_subnets[1]}"
  instance_type          = "m4.large"
  key_name               = "engineering-production"
  user_data              = "${data.template_file.ecs_cloud_config.rendered}"
  iam_instance_profile   = "${var.iam_role}"

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags {
    Name = "vault-b"
  }
}

resource "aws_ecs_task_definition" "vault" {
  provider              = "aws.local"
  family                = "vault"
  container_definitions = "${data.template_file.vault_ecs_task.rendered}"
  network_mode          = "host"

  volume {
    name      = "vault-logs-storage"
    host_path = "${var.data_path}/logs"
  }
}

resource "aws_ecs_service" "vault" {
  provider                           = "aws.local"
  name                               = "vault-cluster"
  cluster                            = "vault"
  task_definition                    = "${aws_ecs_task_definition.vault.arn}"
  desired_count                      = "2"
  deployment_minimum_healthy_percent = "100"
  deployment_maximum_percent         = "200"

  lifecycle {
    create_before_destroy = true
  }
}

# Consul Agent
module "consul" {
  source           = "../../modules/consul-agent"

  # Consul
  encryption_key   = "yYQCMdOCtX73dmYeEQ/NYA=="
  datacenter       = "us-west-2"
  endpoint         = "consul.eng.linuxfoundation.org"
  ecs_cluster_name = "vault"
}


# Create a new load balancer
resource "aws_elb" "vault" {
  provider           = "aws.local"
  name               = "vault-server"
  security_groups    = ["${var.internal_elb_sg}"]
  subnets            = ["${var.internal_subnets}"]
  internal           = true

  listener {
    instance_port      = 8200
    instance_protocol  = "https"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "arn:aws:acm:us-west-2:643009352547:certificate/4938ed7c-e270-4597-84b2-6374db6149f4"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTPS:8200/v1/sys/health"
    interval            = 30
  }

  cross_zone_load_balancing   = true
  idle_timeout                = 400

  tags {
    Name = "vault"
  }
}

# Create a new load balancer attachment
resource "aws_elb_attachment" "vault-a" {
  provider = "aws.local"
  elb      = "${aws_elb.vault.id}"
  instance = "${aws_instance.vault-a.id}"
}

resource "aws_elb_attachment" "vault-b" {
  provider = "aws.local"
  elb      = "${aws_elb.vault.id}"
  instance = "${aws_instance.vault-b.id}"
}

resource "aws_route53_record" "vault" {
  provider = "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "vault"
  type    = "A"

  alias {
    name                   = "${aws_elb.vault.dns_name}"
    zone_id                = "${aws_elb.vault.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "vault-a" {
  provider = "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "vault-a"
  type    = "A"
  ttl     = 300
  records = ["${aws_instance.vault-a.private_ip}"]
}

resource "aws_route53_record" "vault-b" {
  provider = "aws.local"
  zone_id = "${var.route53_zone_id}"
  name    = "vault-b"
  type    = "A"
  ttl     = 300
  records = ["${aws_instance.vault-b.private_ip}"]
}
