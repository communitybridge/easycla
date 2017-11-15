variable "ami" {}
variable "sg" {}
variable "subnet" {}
variable "name" {}

resource "aws_instance" "danvpn" {
  provider               = "aws.local"
  ami                    = "${var.ami}" // Amazon Linux 2017.09.1 (HVM) SSD Volume Type
  vpc_security_group_ids = ["${var.sg}"]
  subnet_id              = "${var.subnet}"
  instance_type          = "t2.small"
  key_name               = "dan"

  ebs_block_device {
    device_name = "/dev/xvdf"
    volume_size = 50
    encrypted   = true
  }

  tags {
    Name = "${var.name}"
    Owner = "dparsons@linuxfoundation.org"
  }
}

resource "aws_eip" "danvpn" {
  provider = "aws.local"
  instance = "${aws_instance.danvpn.id}"
  vpc = true
}