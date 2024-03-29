package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/object"
)

// String : implement Expression
type String struct {
	defaultNode
	Value string
}

func (this *String) Do(v Visitor) error {
	return v.DoString(this)
}

func (this *String) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprString,
		keyValue: this.Value,
	}
}
func (this *String) Decode(b []byte) error {
	return json.Unmarshal(b, &this.Value)
}
func (this *String) expressionNode() {}

func (this *String) String() string {
	return this.Value
}

func (this *String) Eval(e object.Env) (object.Object, error) {
	return object.NewString(this.Value), nil
}
