package template_test

import "testing"

func TestEqualAny(t *testing.T) {
	testIt(t, []testCase{
		{
			name:     "string",
			template: `{{eq_any "1" "1" "2" "3" "4"}}`,
			output:   "true",
		},
		{
			name:     "string false",
			template: `{{eq_any "1" "2" "3" "4"}}`,
			output:   "false",
		},
		{
			name:     "number",
			template: `{{eq_any 1.2 1 2.9 2.0 1.2}}`,
			output:   "true",
		},
		{
			name:     "number false",
			template: `{{eq_any 7.2 1 2.9 2.0 1.2}}`,
			output:   "false",
		},
	})
}
