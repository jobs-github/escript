package ast

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/jobs-github/escript/function"
)

const (
	Suffix     = ".es"
	SuffixJson = ".json"
)

const (
	keyType  = "type"
	keyValue = "value"
)

func loadAst(
	baseDir string,
	module string,
	suffix string,
	decode func(code string) (Node, error),
) (Node, error) {
	path := filepath.Join(baseDir, fmt.Sprintf("%v%v", module, suffix))
	if suffix == SuffixJson {
		b, err := function.LoadFile(path)
		if nil != err {
			return nil, function.NewError(err)
		}
		return Decode(b)
	}
	b, err := function.LoadFile(path)
	if nil != err {
		return nil, function.NewError(err)
	}
	return decode(function.BytesToString(b))
}

func LoadAst(
	baseDir string,
	suffix string,
	decode func(code string) (Node, error),
) func(module string) (Node, error) {
	return func(module string) (Node, error) {
		return loadAst(baseDir, module, suffix, decode)
	}
}

func Decode(b []byte) (Node, error) {
	var root JsonNode
	if err := json.Unmarshal(b, &root); nil != err {
		return nil, function.NewError(err)
	}
	return root.decode()
}

type JsonNode struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}

func (this *JsonNode) decodeExpr() (Expression, error) {
	expr := this.newExpr()
	if nil == expr {
		return nil, fmt.Errorf("unknown expr type: %v", this.Type)
	}
	if err := expr.Decode(this.Value); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}

func (this *JsonNode) decodeBlockStmt() (*BlockStmt, error) {
	stmt := this.newBlockStmt()
	if nil == stmt {
		return nil, fmt.Errorf("not block stmt: %v", this.Type)
	}
	if err := stmt.Decode(this.Value); nil != err {
		return nil, function.NewError(err)
	}
	return stmt, nil
}

func (this *JsonNode) decode() (Node, error) {
	node := this.newNode()
	if nil == node {
		return nil, fmt.Errorf("unknown type: %v", this.Type)
	}
	if err := node.Decode(this.Value); nil != err {
		return nil, function.NewError(err)
	}
	return node, nil
}

func (this *JsonNode) decodeStmt() (Statement, error) {
	stmt := this.newStatement()
	if nil == stmt {
		return nil, fmt.Errorf("not stmt: %v", this.Type)
	}
	if err := stmt.Decode(this.Value); nil != err {
		return nil, function.NewError(err)
	}
	return stmt, nil
}

func (this *JsonNode) decodeIdent() (*Identifier, error) {
	return decodeIdent(this.Value)
}

func (this *JsonNode) decodeFn() (*Function, error) {
	return decodeFn(this.Value)
}

func (this *JsonNode) newBlockStmt() *BlockStmt {
	if this.Type == typeStmtBlock {
		return NewBlock()
	}
	return nil
}

func (this *JsonNode) newStatement() Statement {
	fn, ok := stmtFactory[this.Type]
	if !ok {
		return nil
	}
	return fn()
}

func (this *JsonNode) newExpr() Expression {
	fn, ok := exprFactory[this.Type]
	if !ok {
		return nil
	}
	return fn()
}

func (this *JsonNode) newNode() Node {
	if typeNodeProgram == this.Type {
		return &Program{}
	}
	if r := this.newStatement(); nil != r {
		return r
	}
	if r := this.newExpr(); nil != r {
		return r
	}
	return nil
}

type JsonMap map[string]JsonNode

func newJsonMap(b []byte) (JsonMap, error) {
	var m JsonMap
	if err := json.Unmarshal(b, &m); nil != err {
		return nil, function.NewError(err)
	}
	return m, nil
}

func (this *JsonMap) decodeExprMap() (ExpressionMap, error) {
	m := ExpressionMap{}
	for k, v := range *this {
		n, err := v.decodeExpr()
		if nil != err {
			return nil, function.NewError(err)
		}
		s := NewString()
		s.Value = k
		m[s] = n
	}
	return m, nil
}

type JsonNodes []JsonNode

func newJsonNodes(b []byte) (JsonNodes, error) {
	var arr JsonNodes
	if err := json.Unmarshal(b, &arr); nil != err {
		return nil, function.NewError(err)
	}
	return arr, nil
}

func (this *JsonNodes) decodeStmts() (StatementSlice, error) {
	stmts := StatementSlice{}
	for _, v := range *this {
		n, err := v.decodeStmt()
		if nil != err {
			return nil, function.NewError(err)
		}
		stmts = append(stmts, n)
	}
	return stmts, nil
}

func (this *JsonNodes) decodeExprs() (ExpressionSlice, error) {
	exprs := ExpressionSlice{}
	for _, v := range *this {
		n, err := v.decodeExpr()
		if nil != err {
			return nil, function.NewError(err)
		}
		exprs = append(exprs, n)
	}
	return exprs, nil
}

func (this *JsonNodes) decodeIdents() (IdentifierSlice, error) {
	idents := IdentifierSlice{}
	for _, v := range *this {
		n, err := v.decodeIdent()
		if nil != err {
			return nil, function.NewError(err)
		}
		idents = append(idents, n)
	}
	return idents, nil
}

func decodeStmts(b []byte) (StatementSlice, error) {
	arr, err := newJsonNodes(b)
	if nil != err {
		return nil, function.NewError(err)
	}
	return arr.decodeStmts()
}

func decodeExprs(b []byte) (ExpressionSlice, error) {
	arr, err := newJsonNodes(b)
	if nil != err {
		return nil, function.NewError(err)
	}
	return arr.decodeExprs()
}

func decodeIdents(b []byte) (IdentifierSlice, error) {
	arr, err := newJsonNodes(b)
	if nil != err {
		return nil, function.NewError(err)
	}
	return arr.decodeIdents()
}

func decodeExprMap(b []byte) (ExpressionMap, error) {
	m, err := newJsonMap(b)
	if nil != err {
		return nil, function.NewError(err)
	}
	return m.decodeExprMap()
}

func decodeIdent(b []byte) (*Identifier, error) {
	v := NewIdent()
	if err := v.Decode(b); nil != err {
		return nil, function.NewError(err)
	}
	return v, nil
}

func decodeFn(b []byte) (*Function, error) {
	v := NewFn()
	if err := v.Decode(b); nil != err {
		return nil, function.NewError(err)
	}
	return v, nil
}

func decodeStmt(b []byte) (Statement, error) {
	var v JsonNode
	if err := json.Unmarshal(b, &v); nil != err {
		return nil, function.NewError(err)
	}
	stmt, err := v.decodeStmt()
	if nil != err {
		return nil, function.NewError(err)
	}
	return stmt, nil
}

func decodeExpr(b []byte) (Expression, error) {
	var v JsonNode
	if err := json.Unmarshal(b, &v); nil != err {
		return nil, function.NewError(err)
	}
	expr, err := v.decodeExpr()
	if nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}

func decodeKv(b []byte) (*Identifier, Expression, error) {
	var v struct {
		Name  JsonNode `json:"name"`
		Value JsonNode `json:"value"`
	}
	if err := json.Unmarshal(b, &v); nil != err {
		return nil, nil, function.NewError(err)
	}
	name, err := v.Name.decodeIdent()
	if nil != err {
		return nil, nil, function.NewError(err)
	}
	value, err := v.Value.decodeExpr()
	if nil != err {
		return nil, nil, function.NewError(err)
	}
	return name, value, nil
}
