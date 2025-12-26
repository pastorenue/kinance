locals {
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  
  lambda_function_invoke_arn = module.lambda.lambda_invoke_arn
  lambda_function_name = module.lambda.function_name
}