package template_test

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	tt "github.com/pthethanh/template"
)

type (
	testCase struct {
		name       string
		template   string
		data       interface{}
		output     string
		verifyFunc func(got string) error
	}
)

func TestIsTrue(t *testing.T) {
	testIt(t, []testCase{
		{
			name:     "istrue: string true",
			template: "{{is_true .}}",
			data:     "ok",
			output:   "true",
		},
		{
			name:     "istrue: string false",
			template: "{{is_true .}}",
			data:     "",
			output:   "false",
		},
		{
			name:     "istrue: number false",
			template: "{{is_true .}}",
			data:     0,
			output:   "false",
		},
		{
			name:     "istrue: number true",
			template: "{{is_true .}}",
			data:     1,
			output:   "true",
		},
		{
			name:     "istrue: array false",
			template: "{{is_true .}}",
			data:     []byte{},
			output:   "false",
		},
		{
			name:     "istrue: array true",
			template: "{{is_true .}}",
			data:     []byte("x"),
			output:   "true",
		},
		{
			name:     "istrue: pipeline true",
			template: "{{.|is_true}}",
			data:     []byte("x"),
			output:   "true",
		},
		{
			name:     "istrue: pipeline false",
			template: "{{.|is_true}}",
			data:     []byte{},
			output:   "false",
		},
		{
			name:     "istrue: empty",
			template: "{{.|is_true}}",
			data:     "",
			output:   "false",
		},
	})
}

func TestDefault(t *testing.T) {
	testIt(t, []testCase{
		{
			name:     "default: use val",
			template: `{{.|default "NOK"}}`,
			data:     "OK",
			output:   "OK",
		},
		{
			name:     "default: use default",
			template: `{{.|default "NOK"}}`,
			data:     "",
			output:   "NOK",
		},
		{
			name:     "default: number use default",
			template: `{{.|default "NOK"}}`,
			data:     0,
			output:   "NOK",
		},
		{
			name:     "default: number use val",
			template: `{{.|default "NOK"}}`,
			data:     1,
			output:   "1",
		},
		{
			name:     "default: array use default",
			template: `{{.|default "NOK"}}`,
			data:     []int{},
			output:   "NOK",
		},
		{
			name:     "default: array use val",
			template: `{{.|default "NOK"}}`,
			data:     []int{1, 2, 3},
			output:   "[1 2 3]",
		},
		{
			name:     "default: map use default",
			template: `{{.|default "NOK"}}`,
			data:     map[string]string{},
			output:   "NOK",
		},
		{
			name:     "default: map use val",
			template: `{{.|default "NOK"}}`,
			data:     map[string]string{"x": "y"},
			output:   `map[x:y]`,
		},
	})
}

func TestYesNo(t *testing.T) {
	testIt(t, []testCase{
		{
			name:     "yesno: string ok",
			template: `{{.|yesno "OK" "NOK"}}`,
			data:     "ok",
			output:   "OK",
		},
		{
			name:     "yesno: string nok",
			template: `{{.|yesno "OK" "NOK"}}`,
			data:     "",
			output:   "NOK",
		},
		{
			name:     "yesno: number ok",
			template: `{{.|yesno "OK" "NOK"}}`,
			data:     1,
			output:   "OK",
		},
		{
			name:     "yesno: number nok",
			template: `{{.|yesno "OK" "NOK"}}`,
			data:     0,
			output:   "NOK",
		},
		{
			name:     "yesno: bool ok",
			template: `{{.|yesno "OK" "NOK"}}`,
			data:     true,
			output:   "OK",
		},
		{
			name:     "yesno: bool nok",
			template: `{{.|yesno "OK" "NOK"}}`,
			data:     false,
			output:   "NOK",
		},
	})
}

func TestCoalesce(t *testing.T) {
	type data struct {
		X interface{}
		Y interface{}
		Z interface{}
	}
	testIt(t, []testCase{
		{
			name:     "string first not empty",
			template: `{{coalesce .X .Y .Z}}`,
			data: data{
				X: "1",
				Y: "2",
				Z: "3",
			},
			output: "1",
		},
		{
			name:     "string first empty",
			template: `{{coalesce .X .Y .Z}}`,
			data: data{
				X: "",
				Y: "2",
				Z: "3",
			},
			output: "2",
		},
		{
			name:     "bool first false",
			template: `{{coalesce .X .Y .Z}}`,
			data: data{
				Y: true,
			},
			output: "true",
		},
		{
			name:     "bool first true",
			template: `{{coalesce .X .Y .Z}}`,
			data: data{
				X: true,
			},
			output: "true",
		},
	})
}

func TestEnv(t *testing.T) {
	envVal := "hello"
	os.Setenv("TEST_NAME", envVal)
	tmpl := template.Must(template.New("").Funcs(tt.FuncMap()).Parse(`{{env "TEST_NAME"}}`))
	buff := bytes.Buffer{}
	if err := tmpl.Execute(&buff, nil); err != nil {
		t.Error(err)
	}
	if buff.String() != envVal {
		t.Errorf("got result=%v, want result=%v", buff.String(), envVal)
	}
}

func TestContains(t *testing.T) {
	arr := []int{1, 2}
	x := 5
	testIt(t, []testCase{
		{
			name:     "string true",
			template: `{{contains . "x"}}`,
			data:     "hellox",
			output:   "true",
		},
		{
			name:     "string false",
			template: `{{contains . "x"}}`,
			data:     "hello",
			output:   "false",
		},
		{
			name:     "slice true",
			template: `{{contains . "x"}}`,
			data:     []string{"y", "x"},
			output:   "true",
		},
		{
			name:     "slice false",
			template: `{{contains . "z"}}`,
			data:     []string{"y", "x"},
			output:   "false",
		},
		{
			name:     "map true",
			template: `{{contains . 1}}`,
			data:     map[int]int{0: 0, 1: 1},
			output:   "true",
		},
		{
			name:     "map false",
			template: `{{contains . 2}}`,
			data:     map[int]int{0: 0, 1: 1},
			output:   "false",
		},
		{
			name:     "map multiple one not in map",
			template: `{{contains . 0 1 2}}`,
			data:     map[int]int{0: 0, 1: 1},
			output:   "false",
		},
		{
			name:     "map multiple all exists in map",
			template: `{{contains . 0 1 2}}`,
			data:     map[int]int{0: 0, 1: 1, 2: 2},
			output:   "true",
		},
		{
			name:     "invalid type - false",
			template: `{{contains . 1}}`,
			data:     1,
			output:   "false",
		},
		{
			name:     "pointer array",
			template: `{{contains . 1}}`,
			data:     &arr,
			output:   "true",
		},
		{
			name:     "contains any: map multiple all exists in map",
			template: `{{contains_any . 5 6 2}}`,
			data:     map[int]int{0: 0, 1: 1, 2: 2},
			output:   "true",
		},
		{
			name:     "contains any string",
			template: `{{contains_any . "x" "y"}}`,
			data:     "my name is jack",
			output:   "true",
		},
		{
			name:     "contains any false",
			template: `{{contains_any . "x" "y"}}`,
			data:     "mi name is jack",
			output:   "false",
		},
		{
			name:     "contains any slice of pointer, pointer val",
			template: `{{contains_any (index . "list") (index . "val")}}`,
			data: map[string]interface{}{
				"list": []*int{&x},
				"val":  &x,
			},
			output: "true",
		},
		{
			name:     "contains any pointer slice, pointer val",
			template: `{{contains_any (index . "list") (index . "val")}}`,
			data: map[string]interface{}{
				"list": &arr,
				"val":  &x,
			},
			output: "false",
		},
		{
			name:     "contains any pointer slice, normal inline val",
			template: `{{contains_any . 5}}`,
			data:     []*int{&x},
			output:   "true",
		},
		{
			name:     "contains any normal slice, pointer val",
			template: `{{contains_any (index . "list") (index . "val")}}`,
			data: map[string]interface{}{
				"list": []int{1, 2, 3, 4, 5},
				"val":  &x,
			},
			output: "true",
		},
		{
			name:     "contains nil nil",
			template: `{{contains (index . "list") (index . "val")}}`,
			data: map[string]interface{}{
				"list": nil,
				"val":  nil,
			},
			output: "false",
		},
		{
			name:     "contains val nil",
			template: `{{contains (index . "list") (index . "val")}}`,
			data: map[string]interface{}{
				"list": []interface{}{1, nil},
				"val":  nil,
			},
			output: "true",
		},
	})
}

func TestUUID(t *testing.T) {
	testIt(t, []testCase{
		{
			name:     "uuid",
			template: "{{uuid}}",
			verifyFunc: func(got string) error {
				if _, err := uuid.Parse(got); err != nil {
					return fmt.Errorf("got result=%s, want result is an UUID", got)
				}
				return nil
			},
		},
	})
}

func TestRepeat(t *testing.T) {
	x := 5
	testIt(t, []testCase{
		{
			name:     "repeat string",
			template: `{{.|repeat 3}}`,
			data:     "x",
			output:   "xxx",
		},
		{
			name:     "repeat int",
			template: `{{.|repeat 3}}`,
			data:     3,
			output:   "333",
		},
		{
			name:     "repeat bool",
			template: `{{.|repeat 3}}`,
			data:     true,
			output:   "truetruetrue",
		},
		{
			name:     "repeat make sure result is string",
			template: `{{eq (.|repeat 3) "111"}}`,
			data:     1,
			output:   "true",
		},
		{
			name:     "repeat pointer",
			template: `{{.|repeat 3}}`,
			data:     &x,
			output:   "555",
		},
	})
}

func TestJoin(t *testing.T) {
	x := 5
	s := []interface{}{"1", 2, 3.0, 4.1, &x, true}
	testIt(t, []testCase{
		{
			name:     "join multiple types - map",
			template: `{{join "," 1 "2" 3 .}}`,
			data: map[string]int{
				"x": 1,
				"y": 2,
			},
			verifyFunc: func(s string) error {
				if !strings.HasPrefix(s, "1,2,3,") || len(s) != 9 {
					return fmt.Errorf("got result=%v, want result=%v or result=%v", s, "1,2,3,1,2", "1,2,3,2,1")
				}
				return nil
			},
		},
		{
			name:     "join multiple types - slice",
			template: `{{join "," 1 "2" 3 .}}`,
			data:     []int{1, 2},
			output:   "1,2,3,1,2",
		},
		{
			name:     "join multiple types - slice interface pointer",
			template: `{{.|join ","}}`,
			data:     &s,
			output:   "1,2,3,4.1,5,true",
		},
		{
			name:     "join multiple types - pointer",
			template: `{{join "," 1 2 3 4 .}}`,
			data:     &x,
			output:   "1,2,3,4,5",
		},
	})
}

func testIt(t *testing.T, cases []testCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tmpl := template.Must(template.New("").Funcs(tt.FuncMap()).Parse(c.template))
			buff := bytes.Buffer{}
			if err := tmpl.Execute(&buff, c.data); err != nil {
				t.Error(err)
			}
			if c.verifyFunc != nil {
				if err := c.verifyFunc(buff.String()); err != nil {
					t.Error(err)
				}
				return
			}
			if strings.Compare(buff.String(), c.output) != 0 {
				t.Errorf("got result=%s, want result=%s", buff.String(), c.output)
			}
		})
	}
}
