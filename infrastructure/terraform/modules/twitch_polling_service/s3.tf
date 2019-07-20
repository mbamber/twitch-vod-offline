resource "aws_s3_bucket" "lambda_versions" {
  bucket = format("mbamber-lambda-versions-%s", local.service_hyphens)
  acl    = "private"

  tags = {
    Name       = format("mbamber-lambda-versions_%s", local.service_hyphens)
    Terrraform = "true"
  }
}
