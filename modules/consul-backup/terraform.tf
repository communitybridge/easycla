resource "aws_ecs_task_definition" "consulBackup" {
  family                = "consulBackup"
  container_definitions = "${file("consul-backup-ecs-task.json")}"
}
