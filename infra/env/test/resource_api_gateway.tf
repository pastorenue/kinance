module "api_gateway" {
  source = "../../modules/api_gateway"
  
  gateway_name = var.api_gateway_name
  lambda_function_invoke_arn = local.lambda_function_invoke_arn
  lambda_function_name = local.lambda_function_name
}