package provider

import (
  "context"

  "github.com/hashicorp/terraform-plugin-framework/datasource"
  "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
  "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
  _ datasource.DataSource = &foldfuncDataSource{}
)

// NewFoldfuncDataSource is a helper function to simplify the provider implementation.
func NewFoldfuncDataSource() datasource.DataSource {
  return &foldfuncDataSource{}
}

// foldfuncDataSource is the data source implementation.
type foldfuncDataSource struct{}

type testModel struct {
    Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *foldfuncDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_call"
}

// Schema defines the schema for the data source.
func (d *foldfuncDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema {
		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
            	Computed: true,
            },
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *foldfuncDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state = testModel {
		Name: types.StringValue("test"),
	}
	
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
