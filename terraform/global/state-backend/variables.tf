variable "aws_region" {
  description = "The AWS region to deploy resources in"
  type        = string
  default     = "ap-south-1"
}

variable "state_bucket_name" {
  description = "The name of the S3 bucket to store Terraform state"
  type        = string
}

variable "lock_table_name" {
  description = "DynamoDB table name for Terraform state locking."
  type        = string
  default     = "one-million-rps-terraform-locks"
}