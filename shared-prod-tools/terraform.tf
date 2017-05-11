variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

terraform {
  backend "s3" {
    bucket = "lfe-terraform-states"
    access_key = "AKIAJQ7437CC6PYAZAXQ"
    secret_key = "B3mojX2tskF2bJpMW95kfCQTd2vlgUKSBKq2nJIt"
    region = "us-west-2"
    key = "production-tools/terraform.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

provider "aws" {
  region     = "us-west-2"
  alias      = "western"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

provider "aws" {
  region     = "us-east-2"
  alias      = "eastern"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

# Creating combined zone for prod
resource "aws_route53_zone" "prod" {
  name = "prod.engineering.internal."
  vpc_id = "vpc-e2f2ad85"
  vpc_region = "us-west-2"
}

# Creating our many required S3 Buckets
module "s3_buckets" {
  source             = "./s3"
  access_key         = "${var.access_key}"
  secret_key         = "${var.secret_key}"
}

# First VPC
module "vpc_west" {
  source             = "./base"
  access_key         = "${var.access_key}"
  secret_key         = "${var.secret_key}"
  region             = "us-west-2"
  region_identitier  = "west"
  dns_server         = "10.32.0.2"
  r53_zone_id        = "${aws_route53_zone.prod.zone_id}"

  name               = "Western Production Tools"
  cidr               = "10.32.0.0/24"
  internal_subnets   = ["10.32.0.128/27", "10.32.0.160/27", "10.32.0.192/27"]
  external_subnets   = ["10.32.0.0/27",   "10.32.0.32/27",  "10.32.0.64/27"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]

  newrelic_key       = "bc34e4b264df582c2db0b453bd43ee438043757c"
  key_name           = "production-shared-tools"

  # Pypi Server
  pypi_redis_host    = "pypi-storage.fbnrd8.0001.usw2.cache.amazonaws.com"
  pypi_bucket        = "${module.s3_buckets.pypi_repo_bucket}"

  # Consul
  consul_encryption_key = "9F2n4KWdxSj2Z4MMVqbHqg=="

  # Peering for GHE
  ghe_peering       = true
}

# Second VPC
module "vpc_east" {
  source             = "./base"
  access_key         = "${var.access_key}"
  secret_key         = "${var.secret_key}"
  region             = "us-east-2"
  region_identitier  = "east"
  dns_server         = "10.32.1.2"
  r53_zone_id        = "${aws_route53_zone.prod.zone_id}"

  name               = "Eastern Production Tools"
  cidr               = "10.32.1.0/24"
  internal_subnets   = ["10.32.1.128/27", "10.32.1.160/27", "10.32.1.192/27"]
  external_subnets   = ["10.32.1.0/27",   "10.32.1.32/27",  "10.32.1.64/27"]
  availability_zones = ["us-east-2a",     "us-east-2b",     "us-east-2c"]

  newrelic_key       = "bc34e4b264df582c2db0b453bd43ee438043757c"
  key_name           = "eastern-production-tools"

  # Pypi Server
  pypi_redis_host    = "pypi-storage.fbnrd8.0001.usw2.cache.amazonaws.com"
  pypi_bucket        = "${module.s3_buckets.pypi_repo_bucket}"

  # Consul
  consul_encryption_key = "9F2n4KWdxSj2Z4MMVqbHqg=="

  # Peering for GHE
  ghe_peering       = false
}

# Consul DNS Region-Balancing (FO/HA)
module "consul_dns_failover" {
  source = "./region_failover_dns"

  # General
  dns_name = "consul"
  dns_zone = "${aws_route53_zone.prod.zone_id}"

  # West
  west_elb_dnsname = "${module.vpc_west.consul_elb_cname}"
  west_elb_name = "${module.vpc_west.consul_elb_name}"
  west_elb_zoneid = "${module.vpc_west.consul_elb_zoneid}"

  # East
  east_elb_dnsname = "${module.vpc_east.consul_elb_cname}"
  east_elb_name = "${module.vpc_east.consul_elb_name}"
  east_elb_zoneid = "${module.vpc_east.consul_elb_zoneid}"
}

# PyPi Server DNS Region-Balancing (FO/HA)
module "pypi_dns_failover" {
  source = "./region_failover_dns"

  # General
  dns_name = "pypi"
  dns_zone = "${aws_route53_zone.prod.zone_id}"

  # West
  west_elb_dnsname = "${module.vpc_west.pypi_elb_cname}"
  west_elb_name = "${module.vpc_west.pypi_elb_name}"
  west_elb_zoneid = "${module.vpc_west.pypi_elb_zoneid}"

  # East
  east_elb_dnsname = "${module.vpc_east.pypi_elb_cname}"
  east_elb_name = "${module.vpc_east.pypi_elb_name}"
  east_elb_zoneid = "${module.vpc_east.pypi_elb_zoneid}"
}

# Consul DNS Region-Balancing (FO/HA)
module "consul_dns_servers_failover" {
  source = "./consul_dns_servers_fo_ha"

  # General
  dns_zone = "${aws_route53_zone.prod.zone_id}"

  # West
  west_ec2_machines = ["10.32.0.219", "10.32.0.134", "10.32.0.183"]
  west_ecs_cluster_name = "${module.vpc_west.tools_ecs_name}"
  west_ecs_service_name = "${module.vpc_west.consul_service_name}"

  # East
  east_ec2_machines = ["10.32.1.178", "10.32.1.133", "10.32.1.206"]
  east_ecs_cluster_name = "${module.vpc_east.tools_ecs_name}"
  east_ecs_service_name = "${module.vpc_west.consul_service_name}"
}