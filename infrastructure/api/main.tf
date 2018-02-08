variable "name" {
  description = "API Gateway name"
}

variable "integration_uri" {
  description = "API Gateway integration uri"
}

variable "stages" {
  description = "API Gateway deployment stages"

  type    = "list"
  default = ["prod"]
}

output "id" {
  value = "${aws_api_gateway_rest_api.main.id}"
}

resource "aws_api_gateway_rest_api" "main" {
  name        = "${var.name}"
  description = "Managed by Terraform"

  binary_media_types = [
    "*/*",
  ]
}

resource "aws_api_gateway_method" "main" {
  rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
  resource_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "main" {
  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
  resource_id = "${aws_api_gateway_rest_api.main.root_resource_id}"
  http_method = "${aws_api_gateway_method.main.http_method}"

  type = "AWS_PROXY"
  uri  = "${var.integration_uri}"

  integration_http_method = "POST"
}

resource "aws_api_gateway_resource" "path" {
  path_part = "{proxy+}"

  parent_id   = "${aws_api_gateway_rest_api.main.root_resource_id}"
  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
}

resource "aws_api_gateway_method" "path" {
  rest_api_id   = "${aws_api_gateway_rest_api.main.id}"
  resource_id   = "${aws_api_gateway_resource.path.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "path" {
  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
  resource_id = "${aws_api_gateway_resource.path.id}"
  http_method = "${aws_api_gateway_method.path.http_method}"

  type = "AWS_PROXY"
  uri  = "${var.integration_uri}"

  integration_http_method = "POST"
}

resource "aws_api_gateway_deployment" "hook" {
  count = "${length(var.stages)}"

  depends_on = [
    "aws_api_gateway_integration.main",
    "aws_api_gateway_integration.path",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.main.id}"
  stage_name  = "${element(var.stages, count.index)}"
}
