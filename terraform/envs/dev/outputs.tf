output "vpc_id" {
  description = "ID of the dev VPC."
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets."
  value       = aws_subnet.public[*].id
}

output "public_route_table_id" {
  description = "ID of the public route table."
  value       = aws_route_table.public.id
}

output "load_balancer_security_group_id" {
  description = "Security group ID for the load balancer."
  value       = aws_security_group.load_balancer.id
}

output "api_security_group_id" {
  description = "Security group ID for the API."
  value       = aws_security_group.api.id
}

output "postgres_security_group_id" {
  description = "Security group ID for Postgres."
  value       = aws_security_group.postgres.id
}

output "valkey_security_group_id" {
  description = "Security group ID for Valkey."
  value       = aws_security_group.valkey.id
}

output "api_ecr_repository_name" {
  description = "Name of the API ECR repository."
  value       = aws_ecr_repository.api.name
}

output "api_ecr_repository_url" {
  description = "URL of the API ECR repository."
  value       = aws_ecr_repository.api.repository_url
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster."
  value       = aws_ecs_cluster.main.name
}

output "ecs_cluster_arn" {
  description = "ARN of the ECS cluster."
  value       = aws_ecs_cluster.main.arn
}

output "api_log_group_name" {
  description = "CloudWatch log group for the API."
  value       = aws_cloudwatch_log_group.api.name
}

output "ecs_task_execution_role_arn" {
  description = "IAM role ARN used by ECS tasks to pull images and write logs."
  value       = aws_iam_role.ecs_task_execution.arn
}

output "api_task_definition_arn" {
  description = "ARN of the API ECS task definition."
  value       = aws_ecs_task_definition.api.arn
}

output "ecs_instance_role_arn" {
  description = "IAM role ARN for ECS EC2 instances."
  value       = aws_iam_role.ecs_instance.arn
}

output "ecs_capacity_provider_name" {
  description = "ECS EC2 capacity provider name."
  value       = aws_ecs_capacity_provider.ec2.name
}

output "ecs_autoscaling_group_name" {
  description = "Auto Scaling Group name for ECS EC2 capacity."
  value       = aws_autoscaling_group.ecs.name
}

output "api_service_name" {
  description = "Name of the ECS API service."
  value       = aws_ecs_service.api.name
}

output "api_load_balancer_dns_name" {
  description = "DNS name of the API application load balancer."
  value       = aws_lb.api.dns_name
}

output "api_url" {
  description = "HTTP URL for the API through the application load balancer."
  value       = "http://${aws_lb.api.dns_name}"
}

output "api_target_group_arn" {
  description = "ARN of the API load balancer target group."
  value       = aws_lb_target_group.api.arn
}

output "load_generator_instance_id" {
  description = "EC2 instance ID for the benchmark load generator."
  value       = aws_instance.load_generator.id
}

output "load_generator_public_ip" {
  description = "Public IP address for the benchmark load generator."
  value       = aws_instance.load_generator.public_ip
}
