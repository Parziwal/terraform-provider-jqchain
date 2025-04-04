package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/parziwal/terraform-provider-jqchain/internal/core"
)

var _ function.Function = &reduceFunction{}

func NewReduceFunction() function.Function {
	return &reduceFunction{}
}

type reduceFunction struct{}

func (f *reduceFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "reduce"
}

func (f *reduceFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Description: "Reduces an initial JSON value by applying a sequence of JQ expressions. Each expression has access to a context variable (default: `context`) and produces a result, which is passed to the next expression. Returns the final result as a JSON-formatted string.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "initial",
				Description: "A JSON-formatted string representing the initial input value.",
			},
			function.ListParameter{
				Name:        "reducers",
				ElementType: types.StringType,
				Description: "A list of JQ expressions to evaluate in order. Each expression receives the previous result via the context variable.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *reduceFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var inputData core.ReduceModel

	// Decode the function arguments
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputData.Initial, &inputData.Reducers))
	if resp.Error != nil {
		return
	}

	// Evaluate reducers
	result, err := core.EvaluateJQReducers(inputData)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("Evaluation failed: %s", err.Error()))
		return
	}

	// Return the final value to Terraform
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}
