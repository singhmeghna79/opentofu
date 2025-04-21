package funcs

import (
	"testing"

	"github.com/opentofu/opentofu/internal/lang/marks"
	"github.com/zclconf/go-cty/cty"
)

func TestFlipSensitiveFunc(t *testing.T) {
	tests := []struct {
		name            string
		input           cty.Value
		expectSensitive bool
	}{
		{
			name:            "non-sensitive string becomes sensitive",
			input:           cty.StringVal("hello"),
			expectSensitive: true,
		},
		{
			name:            "sensitive string becomes non-sensitive",
			input:           cty.StringVal("secret").Mark(marks.Sensitive),
			expectSensitive: false,
		},
		{
			name:            "null value remains null and not sensitive",
			input:           cty.NullVal(cty.String),
			expectSensitive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FlipSensitiveFunc.Call([]cty.Value{tt.input})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			hasSensitive := result.HasMark(marks.Sensitive)

			if hasSensitive != tt.expectSensitive {
				t.Errorf("Expected sensitive=%v, got %v", tt.expectSensitive, hasSensitive)
			}

			t.Logf("Input: %v → Output: %v (sensitive: %v)", tt.input.GoString(), result.GoString(), hasSensitive)
		})
	}
}
