package gates

import (
	"math"
	"reflect"
	"unsafe"
)

type Map map[string]Value

func (Map) IsString() bool   { return false }
func (Map) IsInt() bool      { return false }
func (Map) IsFloat() bool    { return false }
func (Map) IsBool() bool     { return false }
func (Map) IsFunction() bool { return false }

func (Map) ToString() string     { return "[object Map]" }
func (Map) ToInt() int64         { return 0 }
func (Map) ToFloat() float64     { return math.NaN() }
func (m Map) ToNumber() Number   { return Float(m.ToFloat()) }
func (Map) ToBool() bool         { return true }
func (Map) ToFunction() Function { return _EmptyFunction }

func (m Map) ToNative() interface{} {
	return toNative(nil, m)
}

func (m Map) toNative(seen map[unsafe.Pointer]interface{}) interface{} {
	if m == nil {
		return map[string]interface{}(nil)
	}
	addr := unsafe.Pointer(reflect.ValueOf(m).Pointer())
	if v, ok := seen[addr]; ok {
		return v
	}
	result := make(map[string]interface{}, len(m))
	seen[addr] = result
	for k, v := range m {
		result[k] = toNative(seen, v)
	}
	return result
}

func (m Map) Equals(other Value) bool {
	o, ok := other.(Map)
	if !ok {
		return false
	}
	return reflect.DeepEqual(m, o)
}

func (m Map) SameAs(other Value) bool { return false }

func (m Map) Get(r *Runtime, key Value) Value {
	if m == nil {
		return Null
	}
	return r.ToValue(m[key.ToString()])
}

func (m Map) Set(r *Runtime, key, value Value) {
	if m == nil {
		return
	}
	m[key.ToString()] = value
}
