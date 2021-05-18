terraform {
  required_providers {
    domains = {
      version = "1.0.0",
      source = "sigmadigital.io/godaddy/domains"
    }
  }
}

provider "domains" {
  host = "https://api.ote-godaddy.com"
  key = "<godaddy-api-key>"
  secret = "<godaddy-api-secret>"
}

data "domains_domain" "test_domain" {
  domain = "test-domain.com"
}

resource "domains_record" "terraform_record" {
  domain = data.domains_domain.test_domain.domain
  data = "www.terraform.io"
  name = "terraform"
  type = "CNAME"
}
