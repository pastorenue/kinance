# Output the API URL
output "api_url" {
  description = "API Gateway URL"
  value       = "http://localhost:4566/restapis/${aws_api_gateway_rest_api.main.id}/prod/_user_request_"
}
