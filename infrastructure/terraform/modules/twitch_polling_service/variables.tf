variable "service_name" {
  type        = "string"
  description = "The name of the service. This must be unique in the AWS account"
}

variable "lambda_entrypoint" {
  type        = "string"
  description = "The entrypoint into the lambda function. (Default: lambdaHandler)"
  default     = "lambdaHandler"
}

locals {
    service_underscores = replace(var.service_name, "-", "_")
    service_hyphens = replace(var.service_name, "_", "-")
}
