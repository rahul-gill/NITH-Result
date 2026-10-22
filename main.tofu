terraform {
  required_providers {
    vercel = {
      source  = "vercel/vercel"
      version = "~> 2.0"
    }
  }
}

provider "vercel" {
  # OpenTofu automatically reads your VERCEL_API_TOKEN environment variable
}

resource "vercel_project" "imported_app" {
  name                       = "nith-result-rg"
  framework                  = "sveltekit-1"
  node_version               = "22.x"
  root_directory             = "web-app"
  install_command            = "pnpm install"
  serverless_function_region = "iad1"
  auto_assign_custom_domains = false
  
  vercel_authentication = {
    deployment_type = "none"
  }

  git_repository = {
    type              = "github"
    repo              = "rahul-gill/NITH-Result"
    production_branch = "master"
  }
}

