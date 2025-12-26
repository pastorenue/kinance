environment = "test"
aws_region  = "us-east-1"

vpc_cidr             = "10.0.0.0/16"
availability_zones   = ["us-east-1a", "us-east-1b"]
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]

lambda_runtime      = "provided.al2"
skip_final_snapshot = true
service_name = "kinance-lambda"

handler = "../../bin/api"
api_gateway_name = "kinance-api-gateway"