variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "newrelic_key" {
  default = "bc34e4b264df582c2db0b453bd43ee438043757c"
}

terraform {
  backend "consul" {
    address = "consul.service.production.consul:8500"
    path    = "terraform/infrastructure"
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

# Creating our many required S3 Buckets
module "s3_buckets" {
  source             = "./s3"
  access_key         = "${var.access_key}"
  secret_key         = "${var.secret_key}"
}

# Creating Cloudtrail audit logs on all AWS Events
module "cloudtrail_audit" {
  source          = "../../modules/cloudtrail"

  aws_account_id  = "643009352547"
  bucket_creation = true
  environment     = "prod"
}

# First VPC
module "vpc_west" {
  source             = "./base"
  access_key         = "${var.access_key}"
  secret_key         = "${var.secret_key}"
  region             = "us-west-2"
  region_identitier  = "west"
  dns_server         = "10.40.0.2"

  name               = "Infrastructure"
  cidr               = "10.40.0.0/21"
  internal_subnets   = ["10.40.0.0/24",   "10.40.1.0/24",   "10.40.2.0/24"]
  external_subnets   = ["10.40.4.0/24",   "10.40.5.0/24",   "10.40.6.0/24"]
  availability_zones = ["us-west-2a",     "us-west-2b",     "us-west-2c"]

  newrelic_key       = "${var.newrelic_key}"
  key_name           = "infrastructure"
  lb_bucket          = "${module.s3_buckets.production_logs}"

  # Consul
  consul_encryption_key = "9F2n4KWdxSj2Z4MMVqbHqg=="
}

# Second VPC
//module "vpc_east" {
//  source             = "./base"
//  access_key         = "${var.access_key}"
//  secret_key         = "${var.secret_key}"
//  region             = "us-east-2"
//  region_identitier  = "east"
//  dns_server         = "10.32.1.2"
//  r53_zone_id        = "${aws_route53_zone.prod.zone_id}"
//
//  name               = "Eastern Production Tools"
//  cidr               = "10.32.1.0/24"
//  internal_subnets   = ["10.32.1.128/27", "10.32.1.160/27", "10.32.1.192/27"]
//  external_subnets   = ["10.32.1.0/27",   "10.32.1.32/27",  "10.32.1.64/27"]
//  availability_zones = ["us-east-2a",     "us-east-2b",     "us-east-2c"]
//
//  newrelic_key       = "${var.newrelic_key}"
//  key_name           = "eastern-production-tools"
//
//  # Consul
//  consul_encryption_key = "9F2n4KWdxSj2Z4MMVqbHqg=="
//
//  # Peering for GHE
//  ghe_peering       = false
//
//  # Nexus Repository
//  nexus             = false
//}


# Peering Request with CINCO Production
module "project_cinco" {
  source = "./project_peering"

  raw_route_tables_id = "${module.vpc_west.raw_route_tables_id}"
  external_rtb_id = "${module.vpc_west.external_rtb_id}"
  project_cidr = "10.32.3.0/24"
  peering_id = "pcx-91f961f8"
}

# Peering Request with PMC Production
module "project_pmc" {
  source = "./project_peering"

  raw_route_tables_id = "${module.vpc_west.raw_route_tables_id}"
  external_rtb_id = "${module.vpc_west.external_rtb_id}"
  project_cidr = "10.32.4.0/24"
  peering_id = "pcx-c4f961ad"
}

output "account_number" {
  value = "643009352547"
}

output "west_vpc_id" {
  value = "${module.vpc_west.vpc_id}"
}

output "west_dns_servers" {
  value = "${module.vpc_west.dns_servers}"
}

output "west_cidr" {
  value = "${module.vpc_west.cidr}"
}

//output "east_vpc_id" {
//  value = "${module.vpc_east.vpc_id}"
//}
//
//output "east_dns_servers" {
//  value = "${module.vpc_east.bind_dns_servers}"
//}
//
//output "east_cidr" {
//  value = "${module.vpc_east.cidr}"
//}

output "newrelic_key" {
  value = "${var.newrelic_key}"
}