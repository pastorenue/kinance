terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  # LocalStack configuration for testing
  skip_credentials_validation = local.skip_credentials_validation
  skip_metadata_api_check     = local.skip_metadata_api_check
  skip_requesting_account_id  = local.skip_requesting_account_id

  endpoints {
    s3         = "http://localhost:4566"
    ec2        = "http://localhost:4566"
    apigateway = "http://localhost:4566"
  }
}
