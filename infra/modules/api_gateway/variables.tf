variable "gateway_name" {
  type = string
}

variable "lambda_function_invoke_arn" {
  type = string
  description = "The invoke arn of the lambda function"
}

variable "lambda_function_name" {
  type = string
  description = "Lambda function name to be invoked"
}