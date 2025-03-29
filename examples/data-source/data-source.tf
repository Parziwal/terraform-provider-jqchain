terraform {
  required_providers {
    foldfunc = {
      source = "parziwal/foldfunc"
    }
  }
}

provider "foldfunc" {
}

data "foldfunc_call" "test" {}

output "result" {
  value = data.foldfunc_call.test
}
