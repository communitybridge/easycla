variable "external_subnets" {
  description = "External VPC Subnets"
  type = "list"
}

variable "vpn_sg" {
  description = "Security Group for the VPN Server"
}

variable "project" {
  description = "Project Name"
}

variable "key_pair" {
  description = "The key-pair to use for the EC2 Instance"
}

variable "pritunl_link" {
  description = "The URL from Pritunl for this Link"
}

variable "vpc_cidr" {
  description = "The CIDR of the VPC in which you are deploying to"
}

data "template_file" "installation-bash" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    PRITUNL_LINK = "${var.pritunl_link}"
    DNS_SERVER   = "${replace(var.vpc_cidr, ".0/24", ".2")}"
  }
}

resource "aws_instance" "pritunl" {
  provider               = "aws.local"
  ami                    = "ami-74871414"
  source_dest_check      = false
  instance_type          = "t2.small"
  subnet_id              = "${element(var.external_subnets, 0)}"
  key_name               = "${var.key_pair}"
  iam_instance_profile   = "PritunlLink"
  vpc_security_group_ids = ["${var.vpn_sg}"]
  user_data              = "${data.template_file.installation-bash.rendered}"
  monitoring             = true

  tags {
    Name        = "${var.project} - IPsec Link"
    Team        = "Engineering"
    Environment = "Production"
  }
}

resource "aws_eip" "pritunl" {
  provider = "aws.local"
  instance = "${aws_instance.pritunl.id}"
  vpc      = true
}