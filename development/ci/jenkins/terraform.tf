variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "vpc_id" {}

variable "sg_jenkins" {}

variable "sg_jenkins_efs" {}

variable "sg_internal_elb" {}

variable "region" {}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    efs_id               = "${aws_efs_file_system.jenkins-home.id}"
    newrelic_license     = "951db34ebed364ea663002571b63db5d3f827758"
    aws_region           = "${var.region}"
  }
}

data "template_cloudinit_config" "userdata" {
  gzip          = false
  base64_encode = true

  part {
    content      = "${data.template_file.ecs_cloud_config.rendered}"
  }
}

# Creating EFS for Tools Storage
resource "aws_efs_file_system" "jenkins-home" {
  provider = "aws.local"
  creation_token = "enginnering-jenkins-home"

  tags {
    Name = "Enginnering - Jenkins Home"
  }
}

resource "aws_efs_mount_target" "efs_mount_1" {
  provider        = "aws.local"
  file_system_id  = "${aws_efs_file_system.jenkins-home.id}"
  subnet_id       = "${var.internal_subnets[0]}"
  security_groups = ["${var.sg_jenkins_efs}"]
}

resource "aws_efs_mount_target" "efs_mount_2" {
  provider        = "aws.local"
  file_system_id  = "${aws_efs_file_system.jenkins-home.id}"
  subnet_id       = "${var.internal_subnets[1]}"
  security_groups = ["${var.sg_jenkins_efs}"]
}

resource "aws_efs_mount_target" "efs_mount_3" {
  provider        = "aws.local"
  file_system_id  = "${aws_efs_file_system.jenkins-home.id}"
  subnet_id       = "${var.internal_subnets[2]}"
  security_groups = ["${var.sg_jenkins_efs}"]
}

data "aws_ami" "amazon-linux-ami" {
  provider    = "aws.local"
  most_recent = true

  filter {
    name = "name"
    values = ["amzn-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "jenkins" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.amazon-linux-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.large"
  subnet_id              = "${element(var.internal_subnets, 0)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${var.sg_jenkins}"]
  monitoring             = true
  user_data              = "${data.template_cloudinit_config.userdata.rendered}"
  private_ip             = "10.32.2.140"

  tags {
    Name        = "Jenkins - Master Instance"
    Team        = "Engineering"
    Environment = "Production"
  }
}

# Create a new load balancer
resource "aws_elb" "jenkins" {
  provider = "aws.local"
  name = "jenkins-master"
  subnets = ["${var.internal_subnets}"]
  security_groups = ["${var.sg_internal_elb}"]
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
    Name = "jenkins-master"
  }
}

# Create a new load balancer attachment
resource "aws_elb_attachment" "jenkins" {
  provider = "aws.local"
  instance = "${aws_instance.jenkins.id}"
  elb      = "${aws_elb.jenkins.id}"
}

resource "aws_route53_record" "jenkins" {
  provider = "aws.local"
  zone_id = "Z2MDT77FL23F9B"
  name    = "jenkins"
  type    = "A"

  alias {
    name                   = "${aws_elb.jenkins.dns_name}"
    zone_id                = "${aws_elb.jenkins.zone_id}"
    evaluate_target_health = true
  }
}

module "iam" {
  source = "./iam-role"

  name = "jenkins"
  environment = "slaves"
}