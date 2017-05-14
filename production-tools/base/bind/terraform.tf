variable "internal_subnets" {
  description = "Internal VPC Subnets"
  type = "list"
}

variable "bind_sg" {
  description = "Security Group for the internal ELB"
}

variable "cidr" {}

data "template_file" "bind-installation-bash" {
  template = "${file("${path.module}/cloud-config.sh.tpl")}"

  vars {
    newrelic_key = "bc34e4b264df582c2db0b453bd43ee438043757c"
    CIDR_PREFIX  = "${replace(var.cidr, ".0/24", ".")}"
  }
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

resource "aws_instance" "bind_1" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.amazon-linux-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(var.internal_subnets, 0)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${var.bind_sg}"]
  monitoring             = true
  user_data              = "${data.template_file.bind-installation-bash.rendered}"
  private_ip             = "${replace(var.cidr, ".0/24", ".150")}"

  tags {
    Name        = "Bind - DNS Server #1"
    Team        = "Engineering"
    Environment = "Production"
  }

  lifecycle {
    ignore_changes = ["user_data"]
  }
}

resource "aws_instance" "bind_2" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.amazon-linux-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(var.internal_subnets, 1)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${var.bind_sg}"]
  monitoring             = true
  user_data              = "${data.template_file.bind-installation-bash.rendered}"
  private_ip             = "${replace(var.cidr, ".0/24", ".180")}"

  tags {
    Name        = "Bind - DNS Server #2"
    Team        = "Engineering"
    Environment = "Production"
  }

  lifecycle {
    ignore_changes = ["user_data"]
  }
}

resource "aws_instance" "bind_3" {
  provider               = "aws.local"
  ami                    = "${data.aws_ami.amazon-linux-ami.id}"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(var.internal_subnets, 2)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${var.bind_sg}"]
  monitoring             = true
  user_data              = "${data.template_file.bind-installation-bash.rendered}"
  private_ip             = "${replace(var.cidr, ".0/24", ".200")}"

  tags {
    Name        = "Bind - DNS Server #3"
    Team        = "Engineering"
    Environment = "Production"
  }

  lifecycle {
    ignore_changes = ["user_data"]
  }
}

output "bind_dns_servers" {
  value = ["${replace(var.cidr, ".0/24", ".150")}", "${replace(var.cidr, ".0/24", ".180")}", "${replace(var.cidr, ".0/24", ".200")}"]
}
