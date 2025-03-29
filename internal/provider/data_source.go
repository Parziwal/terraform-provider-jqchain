package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &jqchainDataSource{}
)

// NewJqchainDataSource is a helper function to simplify the provider implementation.
func NewJqchainDataSource() datasource.DataSource {
	return &jqchainDataSource{}
}

// jqchainDataSource is the data source implementation.
type jqchainDataSource struct{}

type testModel struct {
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *jqchainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_call"
}

// Schema defines the schema for the data source.
func (d *jqchainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *jqchainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state = testModel{
		Name: types.StringValue("test"),
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
