# Terraform Provider for GoDaddy Domains

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.14+
* [Go](https://golang.org/doc/install) 1.16.0 or higher

## Installing the provider

Enter the provider directory and run the following command:

```shell
make install
```

## Using the provider

See the [example](./examples/main.tf) directory for an example usage.

## Importing DNS record

To import a DNS record from GoDaddy, the ID should conform to the following syntax:

`<domain>:<name>:<type>`

Example

```terraform
terraform import domains_record.test_record test-domain.com:terraform:CNAME
```
