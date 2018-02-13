variable "rest_api_id" {
  description = "API Gateway ID"
}

variable "parent_id" {
  description = "API Gateway Parent ID"
}

variable "path_part" {
  description = "API Gateway path"
}

variable "methods" {
  description = "API Gateway endoint methods"

  type    = "list"
  default = ["ANY"]
}

variable "integration_uri" {
  description = "API Gateway integration uri"
}

variable "stages" {
  description = "API Gateway deployment stages"

  type    = "list"
  default = ["prod"]
}

resource "aws_api_gateway_resource" "main" {
  count = "${length(var.path_part) > 0 ? 1 : 0}"

  path_part   = "${var.path_part}"
  parent_id   = "${var.parent_id}"
  rest_api_id = "${var.rest_api_id}"
}

resource "aws_api_gateway_method" "main" {
  count = "${length(var.methods)}"

  rest_api_id   = "${var.rest_api_id}"
  resource_id   = "${aws_api_gateway_resource.main.id}"
  http_method   = "${element(var.methods, count.index)}"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "main" {
  count = "${length(var.methods)}"

  rest_api_id = "${var.rest_api_id}"
  resource_id = "${aws_api_gateway_resource.main.id}"
  http_method = "${aws_api_gateway_method.main.*.http_method[count.index]}"

  // http_method = "${element(var.methods, count.index)}"

  type                    = "AWS_PROXY"
  uri                     = "${var.integration_uri}"
  integration_http_method = "POST"
}

resource "aws_api_gateway_deployment" "hook" {
  count = "${length(var.stages)}"

  depends_on = ["aws_api_gateway_integration.main"]

  rest_api_id = "${var.rest_api_id}"
  stage_name  = "${element(var.stages, count.index)}"
}
