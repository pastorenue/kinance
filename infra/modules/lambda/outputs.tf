output "function_arn" {
  description = "ARN of Lambda function"
  value       = aws_lambda_function.function.arn
}

output "function_name" {
  description = "Names of Lambda function"
  value       =  aws_lambda_function.function.function_name
}

output "function_url" {
  description = "Function URL (if enabled)"
  value       = aws_lambda_function_url.function.function_url
}

output "lambda_role_arn" {
  description = "Lambda execution role ARN"
  value       = aws_iam_role.lambda.arn
}

output "lambda_security_group_id" {
  description = "Lambda security group ID"
  value       = aws_security_group.lambda.id
}

output "lambda_invoke_arn" {
  description = "Lambda invoke arn"
  value = aws_lambda_function.function.invoke_arn
}