terraform {
  required_providers {
    godaddy-domains = {
      version = "1.0.0",
      source = "sigmadigital.io/godaddy/domains"
    }
  }
}

provider "godaddy-domains" {
  host = "https://api.ote-godaddy.com"
  key = "<godaddy-api-key>"
  secret = "<godaddy-api-secret>"
}

data "domains_domain" "test_domain" {
  provider = "godaddy-domains"

  domain = "test-domain.com"
}

resource "domains_record" "test_record" {
  provider = "godaddy-domains"

  domain = data.domains_domain.test_domain.domain
  data = "www.terraform.io"
  name = "terraform"
  type = "CNAME"
}
