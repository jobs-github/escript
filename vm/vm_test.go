package vm

import (
	"fmt"
	"testing"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

func parse(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p, err := parser.New(l)
	if nil != err {
		t.Fatal(err)
	}
	r, err := p.ParseProgram()
	if nil != err {
		t.Fatal(err)
	}
	program, ok := r.(*ast.Program)
	if !ok {
		t.Fatal("parse error, type not Program")
	}
	return program
}

func testIntegerObject(want int64, obj object.Object) error {
	result, ok := obj.(*object.Integer)
	if !ok {
		return function.NewError(fmt.Errorf("object is not integer, got=%v", obj))
	}
	if result.Value != want {
		return function.NewError(fmt.Errorf("object has wrong value, got=%v, want: %v", result.Value, want))
	}
	return nil
}

type vmTestCase struct {
	name  string
	input string
	want  interface{}
}

func testExpectedObject(t *testing.T, want interface{}, v object.Object) {
	switch et := want.(type) {
	case int:
		if err := testIntegerObject(int64(et), v); nil != err {
			t.Errorf("testIntegerObject failed, err: %v", err)
		}
	}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parse(t, tt.input)
			c := compiler.New()
			err := c.Compile(program)
			if nil != err {
				t.Fatal(err)
			}
			vm := New(c.Bytecode())
			if err := vm.Run(); nil != err {
				t.Fatal(err)
			}
			e := vm.LastPopped()
			testExpectedObject(t, tt.want, e)
		})
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "1", 1},
		{"case_2", "2", 2},
		{"case_3", "1 + 2", 3},
		{"case_4", "1 - 2", -1},
		{"case_5", "1 * 2", 2},
		{"case_6", "4 / 2", 2},
		{"case_7", "50 / 2 * 2 + 10 - 5", 55},
		{"case_8", "5 + 5 + 5 + 5 - 10", 10},
		{"case_9", "2 * 2 * 2 * 2 * 2", 32},
		{"case_10", "5 * 2 + 10", 20},
		{"case_11", "5 + 2 * 10", 25},
		{"case_12", "5 * (2 + 10)", 60},
	}
	runVmTests(t, tests)
}
