package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/parziwal/terraform-provider-jqchain/internal/core"
)

var (
	_ datasource.DataSource = &reduceDataSource{}
)

func NewReduceDataSource() datasource.DataSource {
	return &reduceDataSource{}
}

type reduceDataSource struct{}

type reduceDataSourceModel struct{
	Initial     types.String   `tfsdk:"initial"`
	Reducers    []types.String `tfsdk:"reducers"`
		ContextName types.String   `tfsdk:"context_name"`
	Result      types.String `tfsdk:"result"`
}

func (d *reduceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reduce"
}

func (d *reduceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reduces an initial JSON value by applying a sequence of JQ expressions. Each expression has access to a context variable (default: `context`) and produces a result, which is passed to the next expression.",
		Attributes: map[string]schema.Attribute{
			"initial": schema.StringAttribute{
				Description: "A JSON-formatted string representing the initial input value.",
				Required:    true,
			},
			"reducers": schema.ListAttribute{
				Description: "A list of JQ expressions to evaluate in order. Each expression receives the previous result via the context variable.",
				ElementType: types.StringType,
				Required:    true,
			},
			"context_name": schema.StringAttribute{
				Description: "The name of the variable to inject into each JQ expression. Defaults to `context` if not specified.",
				Optional:    true,
			},
			"result": schema.StringAttribute{
				Description: "The final result after all reducers have been applied. Returned as a JSON string.",
				Computed:    true,
			},
		},
	}
}

func (d *reduceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var inputData reduceDataSourceModel

	// Read config
	diags := req.Config.Get(ctx, &inputData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Evaluate reducers
	result, err := core.EvaluateJQReducers(core.ReduceModel{
		Initial: inputData.Initial,
		Reducers: inputData.Reducers,
		ContextName: inputData.ContextName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Evaluation failed", err.Error())
		return
	}

	// Save to state
	inputData.Result = result
	diags = resp.State.Set(ctx, &inputData)
	resp.Diagnostics.Append(diags...)
}
