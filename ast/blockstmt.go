package ast

import (
	"bytes"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// BlockStmt : implement Statement
type BlockStmt struct {
	Stmts StatementSlice
}

func (this *BlockStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeStmtBlock,
		keyValue: this.Stmts.encode(),
	}
}
func (this *BlockStmt) Decode(b []byte) error {
	var err error
	this.Stmts, err = decodeStmts(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *BlockStmt) statementNode() {}

func (this *BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range this.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}
func (this *BlockStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return this.Stmts.eval(true, env, insideLoop)
}
func (this *BlockStmt) walk(cb func(module string)) {
	for _, s := range this.Stmts {
		s.walk(cb)
	}
}
func (this *BlockStmt) doDefer(env object.Env) error { return nil }
