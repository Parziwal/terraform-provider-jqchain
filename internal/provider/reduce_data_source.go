package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/itchyny/gojq"
)

var (
	_ datasource.DataSource = &reduceDataSource{}
)

func NewReduceDataSource() datasource.DataSource {
	return &reduceDataSource{}
}

type reduceDataSource struct{}

type reduceInputModel struct {
	Initial types.String   `tfsdk:"initial"`
	ContextName types.String `tfsdk:"context_name"`
	Reducers []types.String `tfsdk:"reducers"`
	Result  types.String   `tfsdk:"result"`
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
				Required: true,
			},
			"context_name": schema.StringAttribute{
				Description: "The name of the variable to inject into each JQ expression. Defaults to `context` if not specified.",
				Optional: true,
			},
			"reducers": schema.ListAttribute{
				Description: "A list of JQ expressions to evaluate in order. Each expression receives the previous result via the context variable.",
				ElementType: types.StringType,
				Required:    true,
			},
			"result": schema.StringAttribute{
				Description: "The final result after all reducers have been applied. Returned as a JSON string.",
				Computed: true,
			},
		},
	}
}

func (d *reduceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {	
	var inputData reduceInputModel

	// Read config
	diags := req.Config.Get(ctx, &inputData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Unmarshal the initial JSON string into a Go interface{}
	var current interface{}
	if err := json.Unmarshal([]byte(inputData.Initial.ValueString()), &current); err != nil {
		resp.Diagnostics.AddError("Invalid initial JSON", fmt.Sprintf("Failed to parse initial value: %s", err))
		return
	}

	// Set context name
	contextName := "context"
	if !inputData.ContextName.IsNull() && inputData.ContextName.ValueString() != "" {
		contextName = inputData.ContextName.ValueString()
	}

	// Evaluate each reducer expression in sequence
	for _, jqExpr := range inputData.Reducers {
		expr := jqExpr.ValueString()
		query, err := gojq.Parse(expr)
		if err != nil {
			resp.Diagnostics.AddError("Invalid JQ expression", fmt.Sprintf("Could not parse jq expression '%s': %s", expr, err))
			return
		}

		input := map[string]interface{}{
			contextName: current,
		}

		iter := query.Run(input)
		result, ok := iter.Next()
		if !ok {
			resp.Diagnostics.AddError("Empty result", fmt.Sprintf("JQ expression '%s' returned no value", expr))
			return
		}
		if err, ok := result.(error); ok {
			resp.Diagnostics.AddError("JQ evaluation error", fmt.Sprintf("Error evaluating '%s': %s", expr, err))
			return
		}
		current = result
	}

	// Convert final result to JSON string
	final, err := json.Marshal(current)
	if err != nil {
		resp.Diagnostics.AddError("Marshal error", fmt.Sprintf("Could not serialize result: %s", err))
		return
	}

	inputData.Result = types.StringValue(string(final))

	// Save to state
	diags = resp.State.Set(ctx, &inputData)
	resp.Diagnostics.Append(diags...)
}
