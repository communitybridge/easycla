variable "sg_engineering_sandboxes" {}

variable "vpc_id" {}

variable "internal_subnets" {
  type = "list"
}

variable "external_subnets" {
  type = "list"
}

variable "region" {}

variable "availability_zones" {
  type = "list"
}

variable "redis_sg" {}

variable "internal_elb_sg" {}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "engineering-sandboxes"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
    aws_region       = "${var.region}"
  }
}

module "engineering-sandboxes-ecs-cluster" {
  source               = "../../modules/ecs-cluster"
  environment          = "Preprod"
  team                 = "Engineering"
  name                 = "engineering-sandboxes"
  vpc_id               = "${var.vpc_id}"
  subnet_ids           = "${var.internal_subnets}"
  key_name             = "engineering-sandboxes"
  iam_instance_profile = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
  region               = "${var.region}"
  availability_zones   = "${var.availability_zones}"
  instance_type        = "t2.medium"
  security_group       = "${var.sg_engineering_sandboxes}"
  instance_ebs_optimized = false
  desired_capacity     = "3"
  min_size             = "3"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
}

module "rds-cluster" {
  source               = "../../modules/rds-cluster"
  master_username      = "lfengineering"
  name                 = "engineering-sandboxes"
  master_password      = "buanCAWwwAGxUyoU2Fai"
  availability_zones   = "${var.availability_zones}"
  vpc_id               = "${var.vpc_id}"
  subnet_ids           = "${var.internal_subnets}"
  environment          = "Preprod"
  team                 = "Engineering"
  security_groups      = ["${module.engineering-sandboxes-ecs-cluster.security_group_id}"]
  engine               = "mariadb"
  engine_version       = "10.1.19"
  parameter_group_name = "engineering"
  instance_type        = "db.t2.medium"
}

module "redis-cluster" {
  source               = "../../modules/redis-cluster"
  name                 = "sandboxes"
  vpc_id               = "${var.vpc_id}"
  subnet_ids           = "${var.internal_subnets}"
  environment          = "Preprod"
  team                 = "Engineering"
  security_groups      = ["${var.redis_sg}"]
  instance_type        = "cache.t2.small"
}

# LDAP used for Keycloak sandbox
module "open-ldap" {
  source                 = "../../modules/openldap"

  ecs_cluster_name       = "${module.engineering-sandboxes-ecs-cluster.name}"
  ecs_asg_name           = "${module.engineering-sandboxes-ecs-cluster.asg_name}"
  internal_subnets       = "${var.internal_subnets}"
  internal_elb_sg        = "${module.engineering-sandboxes-ecs-cluster.security_group_id}"
  dns_servers            = ["10.32.0.140", "10.32.0.180", "10.32.0.220"]

  vpc_id                 = "${var.vpc_id}"
  region                 = "${var.region}"

  ldap_org               = "linuxfoundation"
  ldap_domain            = "linuxfoundation.org"
  ldap_admin_password    = "ZPw4RRzxLikVdN"
}

module "keycloak" {
  source                 = "../../modules/keycloak"

  ecs_cluster_name       = "${module.engineering-sandboxes-ecs-cluster.name}"
  ecs_asg_name           = "${module.engineering-sandboxes-ecs-cluster.asg_name}"
  subnets                = "${var.external_subnets}"
  internal_elb_sg        = "${var.internal_elb_sg}"
  dns_servers            = ["10.32.0.140", "10.32.0.180", "10.32.0.220"]

  vpc_id                 = "${var.vpc_id}"
  region                 = "${var.region}"
  env                    = "sandbox"

  mysql_db               = "keycloak_sandbox"
  mysql_host             = "engineering-sandboxes.cnfn2tun3mjw.us-west-2.rds.amazonaws.com"
  mysql_pass             = "buanCAWwwAGxUyoU2Fai"
  mysql_port             = "3306"
  mysql_user             = "lfengineering"
}
