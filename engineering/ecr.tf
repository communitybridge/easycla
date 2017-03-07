resource "aws_ecr_repository" "pypi" {
  name = "tools/pypi"
}

resource "aws_ecr_repository" "mongodb" {
  name = "tools/mongodb"
}

resource "aws_ecr_repository" "mailhog" {
  name = "mailhog"
}

resource "aws_ecr_repository" "java" {
  name = "workspace/java"
}

resource "aws_ecr_repository" "node" {
  name = "workspace/node"
}

resource "aws_ecr_repository" "php" {
  name = "workspace/php"
}

resource "aws_ecr_repository" "openjdk" {
  name = "openjdk"
}

resource "aws_ecr_repository" "mailman" {
  name = "mailman"
}

resource "aws_ecr_repository" "preprod" {
  name = "preprod"
}

resource "aws_ecr_repository" "base" {
  name = "base"
}

resource "aws_ecr_repository" "elasticsearch" {
  name = "elasticsearch"
}

resource "aws_ecr_repository" "nginx" {
  name = "nginx"
}

resource "aws_ecr_repository" "redis" {
  name = "redis"
}

resource "aws_ecr_repository" "mariadb" {
  name = "mariadb"
}

resource "aws_ecr_repository" "php-fpm" {
  name = "php-fpm"
}