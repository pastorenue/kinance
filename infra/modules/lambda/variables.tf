variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnets" {
  description = "Private subnet IDs for Lambda"
  type        = list(string)
}

variable "lambda_runtime" {
  description = "Lambda runtime"
  type        = string
}

variable "enable_function_urls" {
  description = "Enable Lambda function URLs"
  type        = bool
  default     = false
}

variable "service_name" {
  type = string
  description = "unique name of the lambda function"
}

variable "handler" {
  type = string
}

variable "memory_size" {
  type = number
  default = 512
}

variable "timeout" {
  type = number
  default = 300
}