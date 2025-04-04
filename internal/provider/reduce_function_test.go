package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFunctionReduce_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionReduceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("result", `"HELLO WORLD!"`),
				),
			},
		},
	})
}

func TestAccFunctionReduce_addition(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionReduceConfig_addition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("result", `{"arrays":{"merge":["a","b","c","d","e"]},"sum":{"result":21}}`),
				),
			},
		},
	})
}

const testAccFunctionReduceConfig_basic = `
locals {
  	initial = "\"hello\""
	reducers = [
		".context + \" world\"",
		".context | ascii_upcase",
		".context + \"!\""
	]
}

output "result" {
  value = provider::jqchain::reduce(local.initial, local.reducers)
}`

const testAccFunctionReduceConfig_addition = `
locals {
  	initial = <<EOT
	{
		"arrays": {
			"array1": ["a", "b"],
			"array2": ["c", "d", "e"]
		},
		"sum": {
			"value1": 10,
			"value2": 11
		}
	}
	EOT

	reducers = [
		".context + { \"arrays\": { \"merge\": .context.arrays.array1 + .context.arrays.array2 }}",
		".context + { \"sum\": { \"result\": .context.sum.value1 + .context.sum.value2 }}"
	]
}

output "result" {
  value = provider::jqchain::reduce(local.initial, local.reducers)
}`
