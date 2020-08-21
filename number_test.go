package template_test

import "testing"

func TestNumber(t *testing.T) {
	testIt(t, []testCase{
		{
			name:     "mul",
			template: `{{mul 1 2 3}}`,
			output:   "6",
		},
		{
			name:     "mul 0",
			template: `{{mul 1 2 3 0}}`,
			output:   "0",
		},
		{
			name:     "div",
			template: `{{div 1 2 2}}`,
			output:   "0.25",
		},
		{
			name:     "div 1",
			template: `{{div 2 1 1 1}}`,
			output:   "2",
		},
		{
			name:     "add",
			template: `{{add 1 2 3 0}}`,
			output:   "6",
		},
		{
			name:     "sum",
			template: `{{sum 1 2 3 -1 1 0}}`,
			output:   "6",
		},
		{
			name:     "sub",
			template: `{{sub 1 2 3}}`,
			output:   "-4",
		},
		{
			name:     "pow",
			template: `{{pow 2 2 2.0}}`,
			output:   "16",
		},
		{
			name:     "add float check using string to avoid float64 problem",
			template: `{{eq ((add 2.1 2.1 2.1)|printf "%.2f") "6.30"}}`,
			output:   "true",
		},
	})
}
