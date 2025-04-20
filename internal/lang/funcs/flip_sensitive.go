// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/opentofu/opentofu/internal/lang/marks"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// a one-argument function that will return the same value it
// got in the argument, but if the input is sensitive, it will output a
// non-sensitive value and vice-versa
var FlipSensitiveFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowNull:        true,
			AllowMarked:      true,
			AllowDynamicType: true,
		},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		// This function only affects the value's marks, so the result
		// type is always the same as the argument type.
		return args[0].Type(), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		val := args[0]

		if val.IsNull() {
			return val, nil
		}

		if val.HasMark(marks.Sensitive) {
			v, m := val.Unmark()
			delete(m, marks.Sensitive)
			return v.WithMarks(m), nil
		}

		return val.Mark(marks.Sensitive), nil
	},
})

func FlipSensitive(v cty.Value) (cty.Value, error) {
	return FlipSensitiveFunc.Call([]cty.Value{v})
}
