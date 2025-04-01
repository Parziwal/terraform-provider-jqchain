package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceReduce_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceReduceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jqchain_reduce.test", "result", `"HELLO WORLD!"`),
				),
			},
		},
	})
}

func TestAccDataSourceReduce_addition(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceReduceConfig_addition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jqchain_reduce.test", "result", `{"arrays":{"merge":["a","b","c","d","e"]},"sum":{"result":21}}`),
				),
			},
		},
	})
}

const testAccDataSourceReduceConfig_basic = `
data "jqchain_reduce" "test" {
	initial = "\"hello\""
	reducers = [
		".context + \" world\"",
		".context | ascii_upcase",
		".context + \"!\""
	]
}`

const testAccDataSourceReduceConfig_addition = `
data "jqchain_reduce" "test" {
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
`