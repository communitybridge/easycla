variable "access_key" {
  description   = "Your AWS Access Key"
}

variable "secret_key" {
  description   = "Your AWS Secret Key"
}

# S3 Bucket for Backups
resource "aws_s3_bucket" "backups" {
  provider      = "aws.western"
  bucket        = "lf-engineering-prod-backups"
  acl           = "private"

  tags {
    Name        = "Backups"
    Environment = "Production"
    Team        = "Engineering"
  }
}

# S3 Bucket for Storing Secrets
resource "aws_s3_bucket" "secrets" {
  provider      = "aws.western"
  bucket        = "lf-engineering-prod-secrets"
  acl           = "private"

  tags {
    Name        = "Secrets"
    Environment = "Production"
    Team        = "Engineering"
  }

  versioning {
    enabled = true
  }
}

# S3 Bucket for Production Logstash Logs
resource "aws_s3_bucket" "production-logs" {
  provider      = "aws.western"
  bucket        = "lf-engineering-prod-logs"
  acl           = "private"

  tags {
    Name        = "Production Logs"
    Environment = "Production"
    Team        = "Engineering"
  }
}

resource "aws_s3_bucket_policy" "b" {
  bucket = "${aws_s3_bucket.production-logs.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Stmt1429136633762",
      "Action": "s3:PutObject",
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::lf-engineering-prod-logs/*",
      "Principal": {
        "AWS": "arn:aws:iam::797873946194:root"
      }
    }
  ]
}
EOF
}

output "backups_bucket" {
  value         = "${aws_s3_bucket.backups.bucket}"
}

output "production_logs" {
  value         = "${aws_s3_bucket.production-logs.bucket}"
}
