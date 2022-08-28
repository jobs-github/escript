package ast

import (
	"bytes"
	"strings"

	"github.com/jobs-github/escript/function"
)

// Array : implement Expression
type Array struct {
	Items ExpressionSlice
}

func (this *Array) Do(v Visitor) error {
	return v.DoArray(this)
}

func (this *Array) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprArray,
		keyValue: this.Items.encode(),
	}
}

func (this *Array) Decode(b []byte) error {
	var err error
	this.Items, err = decodeExprs(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *Array) expressionNode() {}

func (this *Array) String() string {
	var out bytes.Buffer
	items := []string{}
	for _, v := range this.Items {
		items = append(items, v.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")
	return out.String()
}
