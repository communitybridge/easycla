//variable "access_key" {
//  description = "Your AWS Access Key"
//}
//
//variable "secret_key" {
//  description = "Your AWS Secret Key"
//}

provider "aws" {
  region = "us-west-2"
  access_key = "AKIAJZSPEM5HOEMPP67Q"
  secret_key = "VtFYzpbv+TC9RGxEHVEDAneRMGAfVUaW+GswaruV"
}

terraform {
  backend "s3" {
    bucket = "lfe-terraform-states"
    access_key = "AKIAJZSPEM5HOEMPP67Q"
    secret_key = "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"
    region = "us-west-2"
    key = "shared-prod-tools/terraform.tfstate"
  }
}

module "ebs_bckup" {
  source = "github.com/kgorskowski/terraform/modules//tf_ebs_bckup"
  EC2_INSTANCE_TAG = "EBS-Backup"
  RETENTION_DAYS   = 30
  regions          = ["us-west-2"]
  cron_expression  = "22 1 * * ? *"
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/files/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "shared-production-tools"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
  }
}

module "vpc" "engineering_vpc" {
  source             = "../modules/vpc"
  name               = "Shared Production Tools"
  cidr               = "10.50.0.0/16"
  internal_subnets   = ["10.50.0.0/19" ,"10.50.64.0/19", "10.50.128.0/19"]
  external_subnets   = ["10.50.32.0/20", "10.50.96.0/20", "10.50.160.0/20"]
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

module "dns" {
  source = "../modules/dns"
  name   = "prod.engineering.internal"
  vpc_id = "${module.vpc.id}"
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "${module.dns.name}"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost("10.50.0.0/16", 2)}"
}

module "security_groups" {
  source  = "./security_groups"
  cidr    = "10.50.0.0/16"
  vpc_id  = "${module.vpc.id}"
  name    = "engineering"
}

resource "aws_cloudwatch_log_group" "tools" {
  name = "shared-production-infra"
}

module "shared-production-tools-ecs-cluster" {
  source               = "../modules/ecs-cluster"
  environment          = "Production"
  team                 = "Engineering"
  name                 = "shared-production-tools"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  key_name             = "production-shared-tools"
  iam_instance_profile = "arn:aws:iam::433610389961:instance-profile/ecsInstanceRole"
  region               = "us-west-2"
  availability_zones   = "${module.vpc.availability_zones}"
  instance_type        = "t2.micro"
  security_group       = "${module.security_groups.tools-ecs-cluster}"
  instance_ebs_optimized = false
  desired_capacity     = "3"
  min_size             = "3"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
}

// The region in which the infra lives.
output "region" {
  value = "us-west-2"
}

// The internal route53 zone ID.
output "zone_id" {
  value = "${module.dns.zone_id}"
}

// The VPC's CIDR
output "cidr" {
  value = "10.50.0.0/16"
}

// Comma separated list of internal subnet IDs.
output "internal_subnets" {
  value = "${module.vpc.internal_subnets}"
}

// Comma separated list of external subnet IDs.
output "external_subnets" {
  value = "${module.vpc.external_subnets}"
}

// The internal domain name, e.g "stack.local".
output "domain_name" {
  value = "${module.dns.name}"
}

// The VPC availability zones.
output "availability_zones" {
  value = "${module.vpc.availability_zones}"
}

// The VPC security group ID.
output "vpc_security_group" {
  value = "${module.vpc.security_group}"
}

// The VPC ID.
output "vpc_id" {
  value = "${module.vpc.id}"
}

// Comma separated list of internal route table IDs.
output "internal_route_tables" {
  value = "${module.vpc.internal_rtb_id}"
}

// The external route table ID.
output "external_route_tables" {
  value = "${module.vpc.external_rtb_id}"
}

// Internal SSH allows ssh connections from the external ssh security group.
output "sg_internal_ssh" {
  value = "${module.security_groups.internal_ssh}"
}

// External ELB allows traffic from the world.
output "sg_internal_elb" {
  value = "${module.security_groups.internal_elb}"
}

// External ELB allows traffic from the world.
output "sg_vpn" {
  value = "${module.security_groups.vpn}"
}
