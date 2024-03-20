package tsq

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/valyala/fastjson"
)

type JSONValue struct {
	top     *JSONValue
	path    []string
	node    *fastjson.Value
	spaces2 string
}

var _ AnyValue = new(JSONValue)

func newJSONRoot(v *fastjson.Value) *JSONValue {
	if v == nil {
		return nil
	}
	return &JSONValue{
		top:     nil,
		path:    nil,
		node:    v,
		spaces2: "",
	}
}

func JSONNull() *JSONValue {
	return newJSONRoot(new(fastjson.Arena).NewNull())
}

func JSONString(val string) *JSONValue {
	return newJSONRoot(new(fastjson.Arena).NewString(val))
}

func JSONNumber[N AnyRealNumber](val N) *JSONValue {
	return newJSONRoot(new(fastjson.Arena).NewNumberString(fmt.Sprint(val)))
}

func JSONInt64(val int64) *JSONValue {
	return newJSONRoot(new(fastjson.Arena).NewNumberInt(int(val)))
}

func JSONFloat64(val float64) *JSONValue {
	return newJSONRoot(new(fastjson.Arena).NewNumberFloat64(val))
}

func JSONBool(val bool) *JSONValue {
	a := new(fastjson.Arena)
	var v *fastjson.Value
	if val {
		v = a.NewTrue()
	} else {
		v = a.NewFalse()
	}
	return newJSONRoot(v)
}

func JSONArray(vals ...*JSONValue) *JSONValue {
	arr := new(fastjson.Arena).NewArray()
	for i, v := range vals {
		arr.SetArrayItem(i, v.node)
	}
	return newJSONRoot(arr)
}

func JSONObject(vals map[string]*JSONValue) *JSONValue {
	obj := new(fastjson.Arena).NewObject()
	for k, v := range vals {
		obj.Set(k, v.node)
	}
	return newJSONRoot(obj)
}

func ParseJSON(data string) *JSONValue {
	v, err := fastjson.Parse(data)
	if err != nil {
		panic(fmt.Errorf("could not parse JSON: %w", err))
	}
	return newJSONRoot(v)
}

func getAsOrFail[A any](obj *JSONValue, getValue func(*fastjson.Value) (A, error), path []string) A {
	node := obj.node.Get(path...)
	var val A
	var err error
	if node == nil {
		err = fmt.Errorf("value is missing")
	} else {
		val, err = getValue(node)
	}
	if err != nil {
		panic(unexpectedTypeError(path, err, obj))
	}
	return val
}

func newChild(parent *JSONValue, node *fastjson.Value, path []string) *JSONValue {
	return &JSONValue{
		top:     parent.top,
		path:    append(parent.path, path...),
		node:    node,
		spaces2: "",
	}
}

func (o *JSONValue) TopJSON() *JSONValue {
	if o.top == nil {
		return o
	}
	return o.top
}

func (o *JSONValue) Top() AnyValue {
	return o.TopJSON()
}

func (o *JSONValue) PathFromTop() []string {
	return o.path
}

func (o *JSONValue) IsTop() bool {
	return o.top == nil
}

func (o *JSONValue) IsEqual(other AnyValue) bool {
	if other, ok := other.(*JSONValue); ok {
		// checking the type is required to mutate the type from rawString to String for DeepEqual to work
		return o.node.Type() == other.node.Type() &&
			reflect.DeepEqual(o.path, other.path) &&
			reflect.DeepEqual(o.node, other.node)
	}
	return false
}

func (o *JSONValue) MustGetJSON(path ...string) *JSONValue {
	sub := getAsOrFail(o, func(v *fastjson.Value) (*fastjson.Value, error) { return v, nil }, path)
	return newChild(o, sub, path)
}

func (o *JSONValue) MustGet(path ...string) AnyValue {
	return o.MustGetJSON(path...)
}

func (o *JSONValue) MustGetObject(path ...string) map[string]AnyValue {
	obj := getAsOrFail(o, (*fastjson.Value).Object, path)
	result := make(map[string]AnyValue, obj.Len())
	obj.Visit(func(key []byte, v *fastjson.Value) {
		result[string(key)] = newChild(o, v, append(path, string(key)))
	})
	return result
}

func (o *JSONValue) MustGetArray(path ...string) []AnyValue {
	arr := getAsOrFail(o, (*fastjson.Value).Array, path)
	result := make([]AnyValue, len(arr))
	for i, v := range arr {
		result[i] = newChild(o, v, append(path, strconv.Itoa(i)))
	}
	return result
}

func (o *JSONValue) MustGetString(path ...string) string {
	return string(getAsOrFail(o, func(v *fastjson.Value) (string, error) {
		b, err := v.StringBytes()
		return string(b), err
	}, path))
}

func (o *JSONValue) MustGetInt64(path ...string) int64 {
	return getAsOrFail(o, (*fastjson.Value).Int64, path)
}

func (o *JSONValue) MustGetFloat64(path ...string) float64 {
	return getAsOrFail(o, (*fastjson.Value).Float64, path)
}

func (o *JSONValue) MustGetBool(path ...string) bool {
	return getAsOrFail(o, (*fastjson.Value).Bool, path)
}

func (o *JSONValue) Exists(path ...string) bool {
	return o.node.Get(path...) != nil
}

func (o *JSONValue) IsNull(path ...string) bool {
	node := o.node.Get(path...)
	return node != nil && node.Type() == fastjson.TypeNull
}

func (o *JSONValue) Pretty() string {
	if len(o.spaces2) == 0 {
		out := new(strings.Builder)
		writeIndentedJSONValueTo(out, o.node, "", "  ")
		o.spaces2 = out.String()
	}
	return o.spaces2
}

// MarshalIndentJSON extends json.MarshalIndent to support [fastjson.Value]s.
func MarshalIndentJSON(v any, prefix, indent string) ([]byte, error) {
	out := new(strings.Builder)
	switch v := v.(type) {
	case *fastjson.Value:
		writeIndentedJSONValueTo(out, v, prefix, indent)
	case *fastjson.Object:
		writeIndentedJSONObjectTo(out, v, prefix, indent)
	case []*fastjson.Value:
		writeIndentedJSONArrayTo(out, v, prefix, indent)
	default:
		marshalled, err := json.MarshalIndent(v, prefix, indent)
		if err != nil {
			return nil, err
		}
		out.Write(marshalled)
	}
	return []byte(out.String()), nil
}

func writeIndentedJSONValueTo(out io.Writer, node *fastjson.Value, prefix, indent string) {
	switch node.Type() {
	case fastjson.TypeObject:
		v := node.GetObject()
		writeIndentedJSONObjectTo(out, v, prefix, indent)
		return
	case fastjson.TypeArray:
		v := node.GetArray()
		writeIndentedJSONArrayTo(out, v, prefix, indent)
		return
	default:
		v := node.String()
		out.Write([]byte(v))
		return
	}
}

func writeIndentedJSONObjectTo(out io.Writer, obj *fastjson.Object, prefix, indent string) {
	prefixChild := prefix + indent
	out.Write([]byte{'{'})
	i := 0
	stop := obj.Len()
	obj.Visit(func(key []byte, v *fastjson.Value) {
		i += 1
		out.Write([]byte{'\n'})
		out.Write([]byte(prefixChild))
		out.Write([]byte{'"'})
		out.Write(key)
		out.Write([]byte("\": "))
		writeIndentedJSONValueTo(out, v, prefixChild, indent)
		if i == stop {
			out.Write([]byte{'\n'})
			out.Write([]byte(prefix))
		} else {
			out.Write([]byte{','})
		}
	})
	out.Write([]byte{'}'})
}

func writeIndentedJSONArrayTo(out io.Writer, arr []*fastjson.Value, prefix, indent string) {
	prefixChild := prefix + indent
	out.Write([]byte{'['})
	stop := len(arr) - 1
	for i, v := range arr {
		out.Write([]byte{'\n'})
		out.Write([]byte(prefixChild))
		writeIndentedJSONValueTo(out, v, prefixChild, indent)
		if i == stop {
			out.Write([]byte{'\n'})
			out.Write([]byte(prefix))
		} else {
			out.Write([]byte{','})
		}
	}
	out.Write([]byte{']'})
}
