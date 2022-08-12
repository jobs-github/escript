package compiler

import (
	"errors"
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

var (
	errUnsupportedVisitor = errors.New("unsupported visitor")
)

func unsupportedOp(entry string, op *token.Token, node ast.Node) error {
	return fmt.Errorf("%v -> unsupported op %v(%v), (`%v`)", entry, op.Literal, token.ToString(op.Type), node.String())
}

// visitor : implement ast.Visitor
type visitor struct {
	c Compiler
}

func (this *visitor) DoProgram(v *ast.Program) error {
	for _, s := range v.Stmts {
		if err := s.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *visitor) DoConst(v *ast.ConstStmt) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoBlock(v *ast.BlockStmt) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoExpr(v *ast.ExpressionStmt) error {
	if err := v.Expr.Do(this); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoFunction(v *ast.FunctionStmt) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoPrefix(v *ast.PrefixExpr) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoInfix(v *ast.InfixExpr) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	if err := v.Right.Do(this); nil != err {
		return function.NewError(err)
	}

	switch v.Op.Type {
	case token.ADD:
		_, err := this.c.encode(code.OpAdd)
		if nil != err {
			return function.NewError(err)
		}
	default:
		return unsupportedOp(function.GetFunc(), v.Op, v.Right)
	}

	return nil
}

func (this *visitor) DoIdent(v *ast.Identifier) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoConditional(v *ast.ConditionalExpr) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoFn(v *ast.Function) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoCall(v *ast.Call) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoCallMember(v *ast.CallMember) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoObjectMember(v *ast.ObjectMember) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoIndex(v *ast.IndexExpr) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoNull(v *ast.Null) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoInteger(v *ast.Integer) error {
	obj := object.NewInteger(v.Value)
	idx := this.c.addConst(obj)
	_, err := this.c.encode(code.OpConst, idx)
	return err
}

func (this *visitor) DoBoolean(v *ast.Boolean) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoString(v *ast.String) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoArray(v *ast.Array) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoHash(v *ast.Hash) error {
	return function.NewError(errUnsupportedVisitor)
}
