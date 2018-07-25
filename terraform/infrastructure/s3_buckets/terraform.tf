resource "aws_s3_bucket" "kops" {
  provider = "aws.local"
  bucket = "lf-engineering-kops-state-store"
  acl    = "private"

  versioning {
    enabled = true
  }

  tags {
    Name        = "Kops State Store"
  }
}

output "bucket_name" {
  value = "${aws_s3_bucket.kops.id}"
}