variable "access_key" {
  description   = "Your AWS Access Key"
}

variable "secret_key" {
  description   = "Your AWS Secret Key"
}

# S3 Bucket for Backups
resource "aws_s3_bucket" "backups" {
  provider      = "aws.western"
  bucket        = "lf-engineering-backups"
  acl           = "private"

  tags {
    Name        = "Backups"
    Environment = "Production"
    Team        = "Engineering"
  }
}

# S3 Bucket for Backups
resource "aws_s3_bucket" "production-logs" {
  provider      = "aws.western"
  bucket        = "lf-engineering-production-logs"
  acl           = "private"

  tags {
    Name        = "Production Logs"
    Environment = "Production"
    Team        = "Engineering"
  }
}

# S3 Bucket for Backups
resource "aws_s3_bucket" "sandbox-logs" {
  provider      = "aws.western"
  bucket        = "lf-engineering-sandbox-logs"
  acl           = "private"

  tags {
    Name        = "Sandbox Logs"
    Environment = "Sandbox"
    Team        = "Engineering"
  }
}

resource "aws_s3_bucket" "engineering-database-backups" {
  provider      = "aws.western"
  bucket        = "engineering-database-backups"
  acl           = "private"

  tags {
    Name        = "Engineering Database Backups"
    Environment = "Production"
  }
}

output "backups_bucket" {
  value         = "${aws_s3_bucket.backups.bucket}"
}

output "production_logs" {
  value         = "${aws_s3_bucket.production-logs.bucket}"
}

output "sandbox_logs" {
  value         = "${aws_s3_bucket.sandbox-logs.bucket}"
}

output "engineering_database_backups" {
  value         = "${aws_s3_bucket.engineering-database-backups.bucket}"
}