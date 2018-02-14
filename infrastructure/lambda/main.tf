variable "name" {
  description = "Lambda function name"
}

variable "role" {
  description = "Lambda function role"
}

variable "handler" {
  description = "Lambda function handler"
}

variable "runtime" {
  description = "Lambda function runtime"
}

variable "timeout" {
  description = "Lambda function timeout"
}

variable "environment" {
  type        = "map"
  description = "Lambda function environment variables"
  default     = {}
}

variable "permission" {
  type        = "map"
  description = "Lambda permission configuration: statement_id, principal, source_arn"
  default     = {}
}

output "arn" {
  value = "${aws_lambda_function.main.arn}"
}

output "name" {
  value = "${aws_lambda_function.main.function_name}"
}

output "invoke_arn" {
  value = "${aws_lambda_function.main.invoke_arn}"
}

resource "aws_lambda_function" "main" {
  function_name = "${var.name}"

  role    = "${var.role}"
  handler = "${var.handler}"
  runtime = "${var.runtime}"
  timeout = "${var.timeout}"

  environment {
    variables = "${var.environment}"
  }
}

resource "aws_lambda_permission" "main" {
  // careful with this. the source_arn needs to be created before, otherwise  // this won't work.  // count = "${length(var.permission) > 0 ? 1 : 0}"

  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.main.arn}"
  statement_id  = "${var.permission["statement_id"]}"
  principal     = "${var.permission["principal"]}"
  source_arn    = "${var.permission["source_arn"]}"
}
