package core

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReduceModel struct {
	Initial     types.String   `tfsdk:"initial"`
	Reducers    []types.String `tfsdk:"reducers"`
	ContextName types.String   `tfsdk:"context_name"`
}

func EvaluateJQReducers(reduceModel ReduceModel) (types.String, error) {
	// Parse the initial JSON into a Go object
	var currentContext interface{}
	if err := json.Unmarshal([]byte(reduceModel.Initial.ValueString()), &currentContext); err != nil {
		return types.StringNull(), fmt.Errorf("Failed to parse initial JSON: %w", err)
	}
	
	// Set context name
	contextName := "context"
	if !reduceModel.ContextName.IsNull() && reduceModel.ContextName.ValueString() != "" {
		contextName = reduceModel.ContextName.ValueString()
	}

	// Evaluate each reducer expression in sequence
	for _, reducer := range reduceModel.Reducers {
		expression := reducer.ValueString()
		query, err := gojq.Parse(expression)
		if err != nil {
			return types.StringNull(), fmt.Errorf("Invalid JQ expression '%s': %w", expression, err)
		}

		input := map[string]interface{}{
			contextName: currentContext,
		}

		iter := query.Run(input)
		result, ok := iter.Next()
		if !ok {
			return types.StringNull(), fmt.Errorf("JQ expression returned no value: %s", expression)
		}
		if err, ok := result.(error); ok {
			return types.StringNull(), fmt.Errorf("JQ evaluation error in '%s': %w", expression, err)
		}

		currentContext = result
	}

	// Marshal final result back to JSON string
	final, err := json.Marshal(currentContext)
	if err != nil {
		return types.StringNull(), fmt.Errorf("Failed to serialize result: %w", err)
	}

	return types.StringValue(string(final)), nil
}