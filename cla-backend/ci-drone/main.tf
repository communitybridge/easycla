terraform {
  backend "s3" {
    key    = "ci/terraform.tfstate"
    region = "us-east-1"
  }
}

data "terraform_remote_state" "ci" {
  backend = "s3"

  config {
    bucket     = "${var.bucket}"
    key        = "ci/terraform.tfstate"
    region     = "us-east-1"
    encrypt    = true
    access_key = "${var.access_key}"
    secret_key = "${var.secret_key}"
  }
}

provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "us-east-1"
}

# Key pair is used for ssh credentials, logging into ec2 instance without password.
resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = "${file(var.key_file)}"
}

resource "aws_eip" "droneip" {
  instance = "${aws_instance.droneci.id}"
}

# Security group with ip configuration
resource "aws_security_group" "default" {
  name        = "drone_security_group"
  description = "Used by drone ci"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.safe_network_ip}"]
  }

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTPS access from anywhere
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "droneci" {
  ami             = "ami-43a15f3e"                         # Ubuntu 16.04, in us-east-1
  instance_type   = "m5.large"
  security_groups = ["${aws_security_group.default.name}"]
  key_name        = "${aws_key_pair.deployer.id}"

  connection {
    user = "ubuntu"
  }

  # Copy our setup files across to the instance
  provisioner "file" {
    source      = "setup"
    destination = "~/"
  }

  # Execute the setup files, injecting our setup environment into the script
  provisioner "remote-exec" {
    inline = [
      "DOMAIN='${var.domain}' GITHUB_CLIENT_ID='${var.github_client_id}' GITHUB_CLIENT_SECRET='${var.github_client_secret}' bash ~/setup/ci-install.sh",
    ]
  }
}

resource "aws_route53_zone" "primary" {
  name = "${var.domain}"

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_route53_record" "www" {
  zone_id = "${aws_route53_zone.primary.zone_id}"
  name    = "${var.domain}"
  type    = "A"
  ttl     = "300"
  records = ["${aws_eip.droneip.public_ip}"]
}

# After we have connected our elatic ip and zone id, we generate a Let's Encrypt certificate on the droneci instance.
resource "null_resource" "provision_drone_cert" {
  triggers {
    subdomain_id = "${aws_route53_record.www.id}"
  }

  depends_on = [
    "aws_instance.droneci",
  ]

  connection {
    host = "${aws_eip.droneip.public_ip}"
    user = "ubuntu"
  }

  provisioner "remote-exec" {
    inline = [
      "DOMAIN='${var.domain}' EMAIL='${var.email}' bash ~/setup/cert-install.sh",
    ]
  }
}
