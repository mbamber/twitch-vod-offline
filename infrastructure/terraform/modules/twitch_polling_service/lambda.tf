resource "aws_lambda_function" "main" {
  description = format("Lambda function for %s", var.service_name)

  s3_bucket = aws_s3_bucket.lambda_versions.bucket
  s3_key = "versions/latest/source.zip"

  function_name = var.service_name
  handler = var.lambda_entrypoint

  role = aws_iam_role.lambda_execution_role.arn

  runtime = "go1.x"
  timeout = 10
  publish = true

  tags = {
      Name = var.service_name
      Terraform = "true"
  }
}
