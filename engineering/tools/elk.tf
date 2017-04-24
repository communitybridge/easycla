resource "aws_elasticsearch_domain" "elk_es" {
  domain_name           = "lf-engineering-elk"
  elasticsearch_version = "5.1"
}

resource "aws_elasticsearch_domain_policy" "elk_es" {
  domain_name = "${aws_elasticsearch_domain.elk_es.domain_name}"

  access_policies = <<POLICIES
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "es:*",
            "Principal": "*",
            "Effect": "Allow",
            "Condition": {
                "IpAddress": {"aws:SourceIp": "10.40.0.0/16"}
            },
            "Resource": "${aws_elasticsearch_domain.elk_es.arn}"
        }
    ]
}
POLICIES
}
