terraform {
  backend "s3" {
    bucket         = "one-million-rps-tfstate-jyotishmoy-20261006"
    key            = "envs/dev/terraform.tfstate"
    region         = "ap-south-1"
    dynamodb_table = "one-million-rps-terraform-locks"
    encrypt        = true
  }
}