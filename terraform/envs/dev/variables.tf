variable "aws_region" {
  description = "AWS region for the dev environment."
  type        = string
  default     = "ap-south-1"
}

variable "project_name" {
  description = "Project name used in resource names."
  type        = string
  default     = "one-million-rps"
}

variable "environment" {
  description = "Environment name."
  type        = string
  default     = "dev"
}

variable "vpc_cidr" {
  description = "CIDR range for the dev VPC."
  type        = string
  default     = "10.20.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR ranges for public subnets."
  type        = list(string)
  default     = ["10.20.1.0/24", "10.20.2.0/24"]
}

variable "api_image_tag" {
  description = "Docker image tag for the API."
  type        = string
  default     = "latest"
}

variable "api_container_port" {
  description = "Port exposed by the API container."
  type        = number
  default     = 8080
}

variable "api_cpu" {
  description = "CPU units for the API task."
  type        = number
  default     = 256
}

variable "api_memory" {
  description = "Memory in MiB for the API task."
  type        = number
  default     = 512
}

variable "api_desired_count" {
  description = "Number of API ECS tasks to run."
  type        = number
  default     = 2
}


variable "ecs_instance_type" {
  description = "EC2 instance type for ECS capacity."
  type        = string
  default     = "t3.micro"
}

variable "ecs_desired_capacity" {
  description = "Desired number of ECS EC2 instances."
  type        = number
  default     = 2
}

variable "ecs_min_size" {
  description = "Minimum number of ECS EC2 instances."
  type        = number
  default     = 1
}

variable "ecs_max_size" {
  description = "Maximum number of ECS EC2 instances."
  type        = number
  default     = 2
}

variable "load_generator_instance_type" {
  description = "EC2 instance type for the benchmark load generator."
  type        = string
  default     = "t3.micro"
}
