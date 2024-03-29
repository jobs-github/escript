package ast

import (
	"bytes"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// ConstStmt : implement Statement
type ConstStmt struct {
	defaultNode
	Name  *Identifier
	Value Expression
}

func (this *ConstStmt) Do(v Visitor) error {
	return v.DoConst(this)
}

func (this *ConstStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeStmtConst,
		keyValue: map[string]interface{}{
			"name":  this.Name.Encode(),
			"value": this.Value.Encode(),
		},
	}
}
func (this *ConstStmt) Decode(b []byte) error {
	var err error
	this.Name, this.Value, err = decodeKv(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ConstStmt) statementNode() {}

func (this *ConstStmt) String() string {
	var out bytes.Buffer
	out.WriteString(token.Const)
	out.WriteString(" ")
	out.WriteString(this.Name.String())
	out.WriteString(" = ")
	if nil != this.Value {
		out.WriteString(this.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (this *ConstStmt) Eval(e object.Env) (object.Object, error) {
	return evalVar(this.Name, this.Value, e)
}
