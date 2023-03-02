package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// LoopExpr : implement Expression
type LoopExpr struct {
	defaultNode
	Cnt  Expression // loop count
	Body Expression // Function or Identifier
}

func (this *LoopExpr) expressionNode() {}

func (this *LoopExpr) Do(v Visitor) error {
	return v.DoLoop(this)
}

func (this *LoopExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"cnt":  this.Cnt.Encode(),
		"body": this.Body.Encode(),
	}
	return m
}

func (this *LoopExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprLoop,
		keyValue: this.value(),
	}
}
func (this *LoopExpr) Decode(b []byte) error {
	var v struct {
		Cnt  JsonNode `json:"cnt"`
		Body JsonNode `json:"body"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Cnt, err = v.Cnt.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Body, err = v.Body.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *LoopExpr) String() string {
	return ""
}
func (this *LoopExpr) Eval(e object.Env) (object.Object, error) {
	v, err := this.Cnt.Eval(e)
	cnt, err := object.ToInteger(v)
	if nil != err {
		return object.Nil, err
	}
	fn, err := this.Body.Eval(e)
	if !object.Callable(fn) {
		return object.Nil, err
	}
	var state *object.State
	for i := int64(0); i < cnt; i++ {
		v, err := fn.Call(object.Objects{object.NewInteger(i)})
		if nil != err {
			return object.Nil, err
		}
		if s, err := v.AsState(); nil != err {
			return object.Nil, err
		} else {
			state = s
		}
		if state.Quit {
			break
		}
	}
	return state, nil
}
