variable "environment" {
  description = "Environment name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
}

variable "public_subnet_cidrs" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
}

variable "private_subnet_cidrs" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
}

variable "lambda_runtime" {
  description = "Lambda runtime"
  type        = string
}

variable "db_name" {
  default = "kinancedb"
}

variable "db_username" {
  default = "kinanceuser"
}

variable "handler" {
  type = string
}

variable "db_password" {
  description = "Database password (min 8 characters)"
  sensitive   = true
}

variable "cluster_name" {
  default = "kinance-aurora-cluster"
}

variable "skip_final_snapshot" {
  type = bool
  description = "Don't create backup when destroying."
  default = false
}

variable "backup_retention_period" {
  type = number
  description = "How long the backup should be retained for"
  default = 7
}

variable "service_name" {
  type = string
  description = "unique name for the lambda function"
}

variable "api_gateway_name" {
  type = string
}