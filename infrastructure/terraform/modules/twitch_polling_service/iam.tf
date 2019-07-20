data "aws_caller_identity" "current" {}

resource "aws_iam_role" "lambda_execution_role" {
  name        = format("%s_lambda_execution_role", local.service_underscores)
  description = format("Role used by lambda for execution of the %s service", var.service_name)

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "lambda_execution_permissions" {
  name = format("%s_lambda_execution_role", local.service_underscores)
  role = aws_iam_role.lambda_execution_role.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowCloudwatchLogs",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:logs:eu-west-1:${data.aws_caller_identity.current.account_id}:log-group:${local.service_hyphens}",
                "arn:aws:logs:eu-west-1:${data.aws_caller_identity.current.account_id}:log-group:${local.service_hyphens}:execution-logs:*"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role" "codebuild_role" {
    name        = format("%s_codebuild_role", local.service_underscores)
    description = format("Role used by codebuild for building the %s service", var.service_name)

    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "codebuild_permissions" {
  name = format("%s_codebuild_role", local.service_underscores)
  role = aws_iam_role.codebuild_role.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowCloudwatchLogs",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:logs:eu-west-1:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${local.service_hyphens}",
                "arn:aws:logs:eu-west-1:${data.aws_caller_identity.current.account_id}:log-group:/aws/codebuild/${local.service_hyphens}:*:*"
            ]
        },
        {
            "Sid": "AllowArtifactUpload",
            "Action": [
                "s3:PutObject"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:s3:::${aws_s3_bucket.lambda_versions.bucket}/versions/latest/*"
            ]
        }
    ]
}
EOF
}
