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

data "aws_region" "current" {
  current = true
}

data "aws_caller_identity" "current" {}

resource "aws_api_gateway_rest_api" "hook" {
  name        = "staticdrop-hook"
  description = "Managed by Terraform"

  binary_media_types = [
    "*/*",
  ]
}

resource "aws_api_gateway_method" "hookroot" {
  rest_api_id   = "${aws_api_gateway_rest_api.hook.id}"
  resource_id   = "${aws_api_gateway_rest_api.hook.root_resource_id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "hookroot" {
  rest_api_id = "${aws_api_gateway_rest_api.hook.id}"
  resource_id = "${aws_api_gateway_rest_api.hook.root_resource_id}"
  http_method = "${aws_api_gateway_method.hookroot.http_method}"

  type = "AWS_PROXY"
  uri  = "${module.lambda_hook.invoke_arn}"

  integration_http_method = "POST"
}

resource "aws_api_gateway_resource" "hookpath" {
  path_part = "{proxy+}"

  parent_id   = "${aws_api_gateway_rest_api.hook.root_resource_id}"
  rest_api_id = "${aws_api_gateway_rest_api.hook.id}"
}

resource "aws_api_gateway_method" "hookpath" {
  rest_api_id   = "${aws_api_gateway_rest_api.hook.id}"
  resource_id   = "${aws_api_gateway_resource.hookpath.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "hookpath" {
  rest_api_id = "${aws_api_gateway_rest_api.hook.id}"
  resource_id = "${aws_api_gateway_resource.hookpath.id}"
  http_method = "${aws_api_gateway_method.hookpath.http_method}"

  type = "AWS_PROXY"
  uri  = "${module.lambda_hook.invoke_arn}"

  integration_http_method = "POST"
}

resource "aws_api_gateway_deployment" "hook" {
  depends_on = [
    "aws_api_gateway_integration.hookroot",
    "aws_api_gateway_integration.hookpath",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.hook.id}"
  stage_name  = "prod"
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
    source_arn   = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.hook.id}/*"
  }
}
