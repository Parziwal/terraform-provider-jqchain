terraform {
  required_providers {
    jqchain = {
      source = "parziwal/jqchain"
    }
  }
}

provider "jqchain" {
}

data "jqchain_call" "test" {}

output "result" {
  value = data.jqchain_call.test
}
