module "lambda" {
  source = "../../modules/lambda"

  service_name    = var.service_name
  environment     = var.environment
  vpc_id          = module.networking.vpc_id
  private_subnets = module.networking.private_subnet_ids
  lambda_runtime  = var.lambda_runtime
  handler         = var.handler
}
