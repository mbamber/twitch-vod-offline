resource "aws_codebuild_project" "main" {
  name        = var.service_name
  description = format("Builds the source code for %s", var.service_name)

  badge_enabled = true
  build_timeout = 5

  service_role = aws_iam_role.codebuild_role.arn

  environment {
      compute_type = "BUILD_GENERAL1_SMALL"
      image = "aws/codebuild/standard:2.0"
      type = "LINUX_CONTAINER"
  }

  source {
      type = "GITHUB"
      location = "https://github.com/mbamber/twitch-vod-offline.git"
      report_build_status = true
      auth {
          type = "OAUTH"
      }
  }

  artifacts {
      type = "S3"
      location = aws_s3_bucket.lambda_versions.bucket
      path = "versions/latest"
      name = "source.zip"
      packaging = "ZIP"
  }

  tags = {
      Name = var.service_name
      Terraform = "true"
  }
}

resource "aws_codebuild_webhook" "main" {
  project_name = aws_codebuild_project.main.name
  branch_filter = "master"
}
