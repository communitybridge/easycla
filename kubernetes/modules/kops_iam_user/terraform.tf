resource "aws_iam_group" "kops" {
  name = "kops"
  path = "/users/"
}

resource "aws_iam_group_policy_attachment" "kops-ec2" {
  group      = "${aws_iam_group.kops.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

resource "aws_iam_group_policy_attachment" "kops-route53" {
  group      = "${aws_iam_group.kops.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonRoute53FullAccess"
}

resource "aws_iam_group_policy_attachment" "kops-s3" {
  group      = "${aws_iam_group.kops.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_group_policy_attachment" "kops-iam" {
  group      = "${aws_iam_group.kops.name}"
  policy_arn = "arn:aws:iam::aws:policy/IAMFullAccess"
}

resource "aws_iam_group_policy_attachment" "kops-vpc" {
  group      = "${aws_iam_group.kops.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonVPCFullAccess"
}

resource "aws_iam_user" "kops" {
  name = "kops"
  path = "/system/"
}

resource "aws_iam_group_membership" "kops" {
  name = "tf-kops-group-membership"

  users = [
    "${aws_iam_user.kops.name}"
  ]

  group = "${aws_iam_group.kops.name}"
}

output "kops_user" {
  value = "${aws_iam_user.kops.id}"
}

output "kops_group" {
  value = "${aws_iam_group.kops.id}"
}