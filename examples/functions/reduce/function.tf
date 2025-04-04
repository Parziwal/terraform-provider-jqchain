terraform {
  required_providers {
    jqchain = {
      source = "parziwal/jqchain"
    }
  }
}

locals {
  initial = jsonencode({
    users = [
      { name = "alice", active = true },
      { name = "bob", active = false },
      { name = "carol", active = true },
    ]
  })
  reducers = [
    <<EOT
    {
      users: [.context.users[] | select(.active == true)]
    }
    EOT
    ,
    <<EOT
    {
      users: [.context.users[] | . + { upper_case_name: (.name | ascii_upcase) }]
    }
    EOT
  ]
}

output "result" {
  value = provider::jqchain::reduce(local.initial, local.reducers)
}