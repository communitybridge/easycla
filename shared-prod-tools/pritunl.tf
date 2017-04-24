resource "aws_instance" "pritunl" {
  ami                    = "ami-12f27672"
  source_dest_check      = false
  instance_type          = "t2.micro"
  subnet_id              = "${element(module.vpc.external_subnets, 0)}"
  key_name               = "production-shared-tools"
  vpc_security_group_ids = ["${module.security_groups.vpn}"]
  monitoring             = true

  tags {
    Name        = "Pritunl - OpenVPN Master"
    EBS-Backup  = "true"
  }
}

resource "aws_eip" "pritunl" {
  instance = "${aws_instance.pritunl.id}"
  vpc      = true
}

resource "aws_route53_record" "pritunl" {
  zone_id = "Z2MDT77FL23F9B"
  name    = "openvpn.engineering.tux.rocks."
  type    = "A"
  ttl     = "300"
  records = ["${aws_eip.pritunl.public_ip}"]
}

resource "aws_route53_record" "pritunl_internal" {
  zone_id = "${module.dns.zone_id}"
  name    = "vpn.prod.engineering.internal."
  type    = "A"
  records = ["${aws_instance.pritunl.private_ip}"]
  ttl     = "300"
}
