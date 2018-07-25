variable "access_key" {
  description = "Your AWS Access Key"
}

variable "secret_key" {
  description = "Your AWS Secret Key"
}

variable "region" {
  default = "us-west-2"
}

variable "external_dns_name" {
  default = "e.tux.rocks"
}

variable "external_dns_zoneid" {
  default = "Z1CCIH5LOWNGEU"
}

variable "cidr" {
  default = "10.45.0.0/16"
}

terraform {
  backend "consul" {
    address = "consul.eng.linuxfoundation.org:443"
    scheme  = "https"
    path    = "terraform/infrastructure2.0"
    access_token = "99e1dd84-a0dc-2bca-5cab-89291a1db801"
  }
}

provider "aws" {
  region     = "${var.region}"
  alias      = "local"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

resource "aws_cloudwatch_log_group" "infrastructure" {
  provider = "aws.local"
  name = "infra"
}

module "vpc" {
  source             = "../modules/vpc"
  name               = "Infrastructure"
  cidr               = "${var.cidr}"
  internal_subnets   = ["10.45.0.0/19", "10.45.32.0/19", "10.45.64.0/19"]
  external_subnets   = ["10.45.96.0/19", "10.45.128.0/19", "10.45.160.0/19"]
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

module "dhcp" {
  source  = "../modules/dhcp"
  name    = "aws.eng.tux.rocks"
  vpc_id  = "${module.vpc.id}"
  servers = "${cidrhost("${var.cidr}", 2)}"
}

module "dns" {
  source = "../modules/dns"
  name   = "aws.eng.tux.rocks"
  vpc_id = "${module.vpc.id}"
}

module "security_groups" "sgs" {
  source = "security_groups/"
  vpc_id = "${module.vpc.id}"
}

module "s3_buckets" {
  source = "s3_buckets"
}

module "danvpn" {
  source                 = "../modules/danvpn"
  ami                    = "ami-32d8124a" // Amazon Linux AMI 2017.09.1 (HVM), SSD Volume Type
  subnet                 = "${module.vpc.external_subnets[0]}"
  name                   = "danvpn.e.tux.rocks"
  dns_zone_id            = "${var.external_dns_zoneid}"
  vpc_id                 = "${module.vpc.id}"
}

//module "salt" {
//  source                 = "../modules/salt"
//  ami                    = "ami-0c2aba6c" // CentOS Linux 7 x86_64 hvm ebs us-west-2
//  subnet                 = "${module.vpc.external_subnets[0]}"
//  name                   = "salt.e.tux.rocks"
//  dns_zone_id            = "${var.external_dns_zoneid}"
//  vpc_id                 = "${module.vpc.id}"
//  ec2type                = "m4.large"
//}
//
//module "consul" {
//  source                 = "../modules/consul"
//  ami                    = "ami-0c2aba6c" // CentOS Linux 7 x86_64 hvm ebs us-west-2
//  dns_zone_id            = "${var.external_dns_zoneid}"
//  vpc_id                 = "${module.vpc.id}"
//  ec2type                = "t2.medium"
//  subnet-a               = "${module.vpc.external_subnets[0]}"
//  subnet-b               = "${module.vpc.external_subnets[1]}"
//  subnet-c               = "${module.vpc.external_subnets[2]}"
//}
//
//module "vault" {
//  source                 = "../modules/vault"
//  ami                    = "ami-0c2aba6c" // CentOS Linux 7 x86_64 hvm ebs us-west-2
//  dns_zone_id            = "${var.external_dns_zoneid}"
//  vpc_id                 = "${module.vpc.id}"
//  ec2type                = "c4.large"
//  subnet-a               = "${module.vpc.external_subnets[0]}"
//  subnet-b               = "${module.vpc.external_subnets[1]}"
//  subnet-c               = "${module.vpc.external_subnets[2]}"
//}

module "ghe" {
  source = "../modules/ghe"

  vpc_id                 = "${module.vpc.id}"
  replica_count          = "1"
  ghe_sg                 = "${module.security_groups.ghe}"
  elb_sg                 = "${module.security_groups.ghe-elb}"
  it_elb_sg              = "${module.security_groups.it-ghe-elb}"
  internal_subnets       = "${module.vpc.internal_subnets}"
  external_subnets       = "${module.vpc.external_subnets}"
}


// IAM Profiles, Roles & Policies for ECS
module "ecs-iam-profile" {
  source = "./iam-role"

  name = "infra"
  environment = "production"
}

module "efs-consul" {
  source = "../modules/efs"
  name   = "consul"
  display_name = "Consul Storage"
  vpc_id  = "${module.vpc.id}"
  vpc_cidr = "${var.cidr}"
  subnet_ids = "${module.vpc.internal_subnets}"
  security_group = "${module.security_groups.infra-ecs-cluster}"
}

module "efs-vault" {
  source = "../modules/efs"
  name   = "vault"
  display_name = "Vault Storage"
  vpc_id  = "${module.vpc.id}"
  vpc_cidr = "${var.cidr}"
  subnet_ids = "${module.vpc.internal_subnets}"
  security_group = "${module.security_groups.vault-ecs-cluster}"
}

module "efs-nexus" {
  source = "../modules/efs"
  name   = "nexus"
  display_name = "Nexus Storage"
  vpc_id  = "${module.vpc.id}"
  vpc_cidr = "${var.cidr}"
  subnet_ids = "${module.vpc.internal_subnets}"
  security_group = "${module.security_groups.infra-ecs-cluster}"
}

module "efs-mongodb" {
  source = "../modules/efs"
  name   = "mongodb"
  display_name = "Pritunl MongoDB Storage"
  vpc_id  = "${module.vpc.id}"
  vpc_cidr = "${var.cidr}"
  subnet_ids = "${module.vpc.internal_subnets}"
  security_group = "${module.security_groups.infra-ecs-cluster}"
}

data "template_file" "ecs_cloud_config" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    ecs_cluster_name = "infra"
    newrelic_key     = "bc34e4b264df582c2db0b453bd43ee438043757c"
    aws_region       = "us-west-2"
    consul_efs       = "${module.efs-consul.id}"
    nexus_efs        = "${module.efs-nexus.id}"
    mongodb_efs      = "${module.efs-mongodb.id}"
  }
}

module "infrastructure-cluster" {
  source               = "../modules/ecs-cluster"
  environment          = "Production"
  team                 = "Engineering"
  name                 = "infra"
  vpc_id               = "${module.vpc.id}"
  subnet_ids           = "${module.vpc.internal_subnets}"
  key_name             = "engineering-production"
  iam_instance_profile = "${module.ecs-iam-profile.ecsInstanceProfile}"
  region               = "us-west-2"
  availability_zones   = "${module.vpc.availability_zones}"
  instance_type        = "m5.large"
  security_group       = "${module.security_groups.infra-ecs-cluster}"
  instance_ebs_optimized = false
  desired_capacity     = "3"
  min_size             = "3"
  cloud_config_content = "${data.template_file.ecs_cloud_config.rendered}"
  root_volume_size     = 100
  docker_volume_size   = 100
}

module "logstash" {
  source                 = "../modules/logstash"

  ecs_cluster_name       = "${module.infrastructure-cluster.name}"
  ecs_asg_name           = "${module.infrastructure-cluster.asg_name}"
  vpc_id                 = "${module.vpc.id}"
  region                 = "${var.region}"
  internal_subnets       = "${module.vpc.internal_subnets}"
  route53_zone_id        = "Z1WGT54F777KX0"
  internal_elb_sg        = "${module.security_groups.logstash-elb}"
}

module "consul-server" {
  source               = "../modules/consul-server"
  ecs_cluster_name     = "infra"
  ecs_asg_name         = "${module.infrastructure-cluster.asg_name}"
  internal_subnets     = "${module.vpc.internal_subnets}"
  internal_elb_sg      = "${module.security_groups.consul-elb}"
  region               = "us-west-2"
  data_path            = "/mnt/storage/consul"
  route53_zone_id      = "Z1WGT54F777KX0"
}

module "vault-server" {
  source               = "./vault-server"
  internal_subnets     = "${module.vpc.internal_subnets}"
  internal_elb_sg      = "${module.security_groups.vault-elb}"
  region               = "us-west-2"
  data_path            = "/mnt/storage/vault"
  route53_zone_id      = "Z1WGT54F777KX0"
  efs_id               = "${module.efs-vault.id}"
  ecs_sg               = "${module.security_groups.vault-ecs-cluster}"
  iam_role             = "${module.ecs-iam-profile.ecsInstanceProfile}"
}

//module "vault-ui" {
//  source               = "../modules/vault-ui"
//  ecs_cluster_name     = "infra"
//  ecs_asg_name         = "${module.infrastructure-cluster.asg_name}"
//  internal_subnets     = "${module.vpc.internal_subnets}"
//  internal_elb_sg      = "${module.security_groups.vault-elb}"
//  region               = "us-west-2"
//  route53_zone_id      = "Z1WGT54F777KX0"
//  vpc_id               = "${module.vpc.id}"
//  ecs_role             = "${module.ecs-iam-profile.service_role_arn}"
//}

module "nexus-server" {
  source               = "../modules/nexus"
  ecs_cluster_name     = "infra"
  ecs_asg_name         = "${module.infrastructure-cluster.asg_name}"
  internal_subnets     = "${module.vpc.internal_subnets}"
  internal_elb_sg      = "${module.security_groups.nexus-elb}"
  region               = "us-west-2"
  data_path            = "/mnt/storage/nexus"
  route53_zone_id      = "Z1WGT54F777KX0"
}

module "pritunl-mongodb" {
  source               = "../modules/mongodb"
  ecs_cluster_name     = "infra"
  ecs_asg_name         = "${module.infrastructure-cluster.asg_name}"
  internal_subnets     = "${module.vpc.internal_subnets}"
  internal_elb_sg      = "${module.security_groups.mongodb-elb}"
  region               = "us-west-2"
  data_path            = "/mnt/storage/mongodb"
  route53_zone_id      = "Z1WGT54F777KX0"
  vpc_id               = "${module.vpc.id}"
}

module "pritunl" {
  source            = "../modules/pritunl"
  external_subnets  = "${module.vpc.external_subnets}"
  internal_subnets  = "${module.vpc.internal_subnets}"
  sg                = "${module.security_groups.pritunl_node}"
  elb_sg            = "${module.security_groups.pritunl_elb}"
  dns_zone_id       = "Z1WGT54F777KX0"
}

# Setup Peering with Staging
module "staging_peering" {
  source = "../modules/peering_setup"

  raw_route_tables_id = "${module.vpc.raw_route_tables_id}"
  external_rtb_id = "${module.vpc.external_rtb_id}"
  project_cidr = "10.37.0.0/16"
  peering_id = "pcx-c882b1a1"
}

//# Setup Peering with Production
//module "production_peering" {
//  source = "../modules/peering_setup"
//
//  raw_route_tables_id = "${module.vpc.raw_route_tables_id}"
//  external_rtb_id = "${module.vpc.external_rtb_id}"
//  project_cidr = "10.45.0.0/16"
//  peering_id = "pcx-c882b1a1"
//}

# Setup Peering with CI
module "ci_peering" {
  source = "../modules/peering_setup"

  raw_route_tables_id = "${module.vpc.raw_route_tables_id}"
  external_rtb_id = "${module.vpc.external_rtb_id}"
  project_cidr = "10.32.2.0/24"
  peering_id = "pcx-1efecc77"
}

# Setup Peering with Sandboxes
module "sandboxes_peering" {
  source = "../modules/peering_setup"

  raw_route_tables_id = "${module.vpc.raw_route_tables_id}"
  external_rtb_id = "${module.vpc.external_rtb_id}"
  project_cidr = "10.32.1.0/24"
  peering_id = "pcx-44e2d02d"
}

# Setup Peering with Sandboxes
module "it_ldap_peering" {
  source = "../modules/peering_setup"

  raw_route_tables_id = "${module.vpc.raw_route_tables_id}"
  external_rtb_id = "${module.vpc.external_rtb_id}"
  project_cidr = "10.30.112.0/21"
  peering_id = "pcx-3be9db52"
}

# Setup Peering with Certification
module "cert_peering" {
  source = "../modules/peering_setup"

  raw_route_tables_id = "${module.vpc.raw_route_tables_id}"
  external_rtb_id = "${module.vpc.external_rtb_id}"
  project_cidr = "10.150.0.0/26"
  peering_id = "pcx-ab132fc2"
}

# IAM Account for Vault AWS Integration
resource "aws_iam_user" "vault-iam-user" {
  provider = "aws.local"
  name = "vault-iam-user"
  path = "/system/"
}

resource "aws_iam_user_policy" "vault-iam-user" {
  provider = "aws.local"
  name = "iam-full-permissions"
  user = "${aws_iam_user.vault-iam-user.name}"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "iam:*",
            "Resource": "*"
        }
    ]
}
EOF
}


output "account_number" {
  value = "643009352547"
}

output "vpc_id" {
  value = "${module.vpc.id}"
}

output "cidr" {
  value = "${var.cidr}"
}
