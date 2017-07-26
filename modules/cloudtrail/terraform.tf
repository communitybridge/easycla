variable "aws_account_id" {}

variable "bucket_creation" {}

variable "environment" {}

resource "aws_cloudtrail" "audit" {
  name                          = "aws-audit"
  s3_bucket_name                = "${aws_s3_bucket.audit-logs.id}"
  s3_key_prefix                 = "${var.aws_account_id}"
  include_global_service_events = true
  sns_topic_name                = "${aws_sns_topic.audit-action.name}"

}

resource "aws_sns_topic" "audit-action" {
  display_name = "AWS API Actions"
  name = "account-wide-api-actions"
}

resource "aws_sns_topic_policy" "custom" {
  arn = "${aws_sns_topic.audit-action.arn}"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "default",
  "Statement":[
    {
      "Sid": "default",
      "Effect": "Allow",
      "Principal": {"AWS":"*"},
      "Action": [
        "SNS:GetTopicAttributes",
        "SNS:SetTopicAttributes",
        "SNS:AddPermission",
        "SNS:RemovePermission",
        "SNS:DeleteTopic"
      ],
      "Resource": "${aws_sns_topic.audit-action.arn}"
    },
    {
        "Sid": "AWSCloudTrailSNSPolicy20131101",
        "Effect": "Allow",
        "Principal": {
          "Service": "cloudtrail.amazonaws.com"
        },
        "Action": "SNS:Publish",
        "Resource": "${aws_sns_topic.audit-action.arn}"
    }
  ]
}
POLICY
}

resource "aws_s3_bucket" "audit-logs" {
  count         = "${var.bucket_creation}"
  bucket        = "lf-engineering-${var.environment}-cloudtrail-audit"
  force_destroy = true

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AWSCloudTrailAclCheck",
            "Effect": "Allow",
            "Principal": {
              "Service": "cloudtrail.amazonaws.com"
            },
            "Action": "s3:GetBucketAcl",
            "Resource": "arn:aws:s3:::lf-engineering-${var.environment}-cloudtrail-audit"
        },
        {
            "Sid": "AWSCloudTrailWrite",
            "Effect": "Allow",
            "Principal": {
              "Service": "cloudtrail.amazonaws.com"
            },
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::lf-engineering-${var.environment}-cloudtrail-audit/*",
            "Condition": {
                "StringEquals": {
                    "s3:x-amz-acl": "bucket-owner-full-control"
                }
            }
        }
    ]
}
POLICY
}