package object

import (
	"reflect"
	"strings"
)

type Type int

const (
	TypeUndefined Type = iota
	TypeBool
	TypeFloat
	TypeFn
	TypeInt
	TypeMap
	TypeSlice
	TypeString
	TypeRune
	TypeStruct
	TypeMapPointer
	TypeRgStruct
	TypeFnLiteral
	TypeCall
	TypeClosure
)

var typeNames = map[Type]string{
	TypeUndefined:  "undefined",
	TypeBool:       "bool",
	TypeFloat:      "float",
	TypeFn:         "fn",
	TypeInt:        "int",
	TypeMap:        "map",
	TypeSlice:      "slice",
	TypeString:     "string",
	TypeRune:       "rune",
	TypeStruct:     "struct",
	TypeMapPointer: "*map",
	TypeRgStruct:   "struct",
	TypeFnLiteral:  "fn",
	TypeCall:       "call",
	TypeClosure:    "closure",
}

func (t Type) String() string {
	return typeNames[t]
}

type Object interface {
	object()
	Type() Type
	GetValue() any
}

var Null = &Empty{}

type Empty struct{}

func (n *Empty) object() {}

func (n *Empty) Type() Type {
	return TypeUndefined
}

func (n *Empty) GetValue() any {
	return nil
}

func New(val any) Object {
	return ToObject(reflect.ValueOf(val))
}

func ToObject(value reflect.Value) Object {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Int{Val: value.Int()}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &Int{Val: int64(value.Uint())}
	case reflect.Float32, reflect.Float64:
		return &Float{Val: value.Float()}
	case reflect.String:
		return &String{Val: []rune(value.String())}
	case reflect.Bool:
		return &Bool{Val: value.Bool()}
	case reflect.Slice, reflect.Array:
		slice := make([]any, value.Len())
		for i := 0; i < value.Len(); i++ {
			slice[i] = value.Index(i).Interface()
		}
		return &Slice{Val: slice}
	case reflect.Map:
		m := make(map[any]any)
		for _, k := range value.MapKeys() {
			m[k.Interface()] = value.MapIndex(k).Interface()
		}
		return &Map{Val: m}
	case reflect.Struct:
		return &InjectStruct{
			Value:           value,
			methodToArgType: make(map[string][]reflect.Type),
		}
	case reflect.Ptr:
		if value.Elem().Kind() == reflect.Struct {
			if strings.Contains(value.Elem().String(), "object.Struct") {
				a := value.Elem().Interface().(Struct)
				return &a
			}
			return &InjectStruct{
				Value:           value,
				methodToArgType: make(map[string][]reflect.Type),
			}
		}
		if value.Elem().Kind() == reflect.Map {
			return &MapPointer{Value: value.Elem()}
		}
		return ToObject(value.Elem())
	case reflect.Func:
		return &InjectFn{Value: value}
	case reflect.Interface:
		return ToObject(value.Elem())
		// TODO: 这里可以根据需要添加更多类型处理逻辑
	default:
		return Null
	}
}

func TransToBool(obj Object) bool {
	switch obj.Type() {
	case TypeBool:
		return obj.(*Bool).Val
	case TypeFloat:
		return obj.(*Float).Val > 0
	case TypeInt:
		return obj.(*Int).Val > 0
	case TypeMap:
		return len(obj.(*Map).Val) > 0
	case TypeSlice:
		return len(obj.(*Slice).Val) > 0
	case TypeString:
		return len(obj.(*String).Val) > 0
	case TypeRune:
		return obj.(*Rune).Val != 0
	default:
		return false
	}
}
