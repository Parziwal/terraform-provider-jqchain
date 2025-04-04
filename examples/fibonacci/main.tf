terraform {
  required_providers {
    jqchain = {
      source = "parziwal/jqchain"
    }
  }
}

data "jqchain_reduce" "fibonacci" {
  initial = jsonencode({
    fibonacci = [1, 1]
  })
  reducers = [
    for count in range(40) :
    <<EOT
    {
      fibonacci: .context.fibonacci + [.context.fibonacci | .[length -1] + .[length -2]]
    }
    EOT
  ]
}

output "result" {
  value = jsondecode(data.jqchain_reduce.fibonacci.result)
}
