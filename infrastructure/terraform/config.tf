provider "aws" {
  region = "eu-west-1"

  assume_role {
    role_arn = "arn:aws:iam::324219055318:role/administrator"
  }
}

terraform {
  backend "s3" {
    bucket   = "mbamber-terraform-remote-state"
    key      = "tvo/terraform.tfstate"
    region   = "eu-west-1"
    role_arn = "arn:aws:iam::324219055318:role/administrator"
  }

  required_version = ">=0.12"
}
