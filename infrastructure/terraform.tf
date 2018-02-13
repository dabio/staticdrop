terraform {
  backend "s3" {
    encrypt = true
    bucket  = "tf-infra"
    key     = "staticdrop.tfstate"
    region  = "eu-central-1"
  }
}

variable "apex_function_hook" {}
variable "apex_function_role" {}

variable "apex_function_arns" {
  type = "map"
}

variable "apex_function_names" {
  type = "map"
}

variable "dropbox_app_key" {}
variable "dropbox_app_secret" {}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

resource "aws_api_gateway_rest_api" "main" {
  name        = "staticdrop"
  description = "Managed by Terraform"

  binary_media_types = [
    "*/*",
  ]
}

module "api_hook" {
  source = "./api"

  // is actually: "/hook"
  path_part   = "hook"
  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
  parent_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
  methods     = ["GET", "POST"]
  stages      = ["prod"]

  integration_uri = "${module.lambda_hook.invoke_arn}"
}

module "lambda_hook" {
  source  = "./lambda"
  name    = "${var.apex_function_names["hook"]}"
  role    = "${var.apex_function_role}"
  handler = "main"
  runtime = "go1.x"
  timeout = 5

  environment {
    APEX_FUNCTION_NAME   = "hook"
    LAMBDA_FUNCTION_NAME = "${var.apex_function_names["hook"]}"
    DROPBOX_APP_KEY      = "${var.dropbox_app_key}"
    DROPBOX_APP_SECRET   = "${var.dropbox_app_secret}"
  }

  permission {
    statement_id = "AllowExecutionFromAPIGateway"
    principal    = "apigateway.amazonaws.com"
    source_arn   = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.main.id}/*"
  }
}
