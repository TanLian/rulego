package object

import (
	"fmt"
	"reflect"
)

type Struct struct {
	Value           reflect.Value
	methodToArgType map[string][]reflect.Type
}

func (s *Struct) object() {}

func (s *Struct) Type() Type {
	return TypeStruct
}

func (s *Struct) GetValue() any {
	return s.Value.Interface()
}

func (s *Struct) Call(name string, args []Object) (Object, bool) {
	m := s.Value.MethodByName(name)
	if !m.IsValid() {
		return Null, false
	}

	types, ok := s.lazyLoadMethodType(name)
	if !ok {
		return Null, false
	}

	// 参数数量不匹配
	if len(types) != len(args) {
		return Null, false
	}

	var in []reflect.Value
	for i, v := range args {
		in = append(in, convert(v.GetValue(), types[i]))
	}

	out := m.Call(in)
	if len(out) == 0 {
		return Null, true
	}
	if len(out) == 1 {
		return ToObject(out[0]), true
	}
	res := &Slice{}
	for i := 0; i < len(out); i++ {
		res.Val = append(res.Val, ToObject(out[i]))
	}
	return res, true
}

func (s *Struct) SetField(fieldName string, value any) {
	f := s.Value.Elem().FieldByName(fieldName)
	if !f.IsValid() || !f.CanSet() {
		panic(fmt.Sprintf("cannot set field %s", fieldName))
	}

	val := reflect.ValueOf(value)
	if f.Type() != val.Type() {
		val = val.Convert(f.Type())
	}
	f.Set(val)
}

func (s *Struct) lazyLoadMethodType(name string) ([]reflect.Type, bool) {
	if types, ok := s.methodToArgType[name]; ok {
		return types, true
	}

	m, ok := s.Value.Type().MethodByName(name)
	if !ok {
		return nil, false
	}

	var types []reflect.Type
	for i := 1; i < m.Type.NumIn(); i++ {
		types = append(types, m.Type.In(i))
	}
	s.methodToArgType[name] = types
	return types, true
}

// convert is used to convert a value of any type into a value of a specified type.
func convert(value any, t reflect.Type) reflect.Value {
	rflValue := reflect.ValueOf(value)
	if rflValue.Type().ConvertibleTo(t) {
		return rflValue.Convert(t)
	}

	switch rflValue.Kind() {
	case reflect.Int64: // int64 => float64
		if t.Kind() == reflect.Float64 {
			return reflect.ValueOf(float64(rflValue.Int())).Convert(t)
		}
	case reflect.Float64: // float64 => int64
		if t.Kind() == reflect.Int64 {
			return reflect.ValueOf(int(rflValue.Float())).Convert(t)
		}
	case reflect.Slice: // []any => []T
		if t.Kind() != reflect.Slice {
			panic(fmt.Sprintf("Element %v cannot be converted to target type %v", value, t))
		}

		input := value.([]any)
		result := reflect.MakeSlice(reflect.SliceOf(t.Elem()), len(input), len(input))
		for i, v := range input {
			vv := reflect.ValueOf(v)
			if !vv.Type().ConvertibleTo(t.Elem()) {
				panic(fmt.Sprintf("Element %v cannot be converted to target type %v", value, t))
			}
			result.Index(i).Set(vv.Convert(t.Elem()))
		}
		return result
	case reflect.Map: // map[any]any => map[K]V
		keyType := t.Key()
		valueType := t.Elem()
		result := reflect.MakeMap(reflect.MapOf(keyType, valueType))
		for k, v := range value.(map[any]any) {
			key := reflect.ValueOf(k)
			val := reflect.ValueOf(v)
			if !key.Type().ConvertibleTo(keyType) || !val.Type().ConvertibleTo(valueType) {
				panic(fmt.Sprintf("Key %v or value %v cannot be converted to target types", k, v))
			}
			result.SetMapIndex(key.Convert(keyType), val.Convert(valueType))
		}
		return result
	default:
	}
	return rflValue
}

func ParseStruct(s interface{}) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// 确保传入的是结构体类型
	if typ.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct")
		return
	}

	fmt.Printf("Parsing struct: %s\n", typ.Name())

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 打印字段名称、类型和值
		fmt.Printf("Field: %s, Type: %s, Value: %v\n", fieldType.Name, fieldType.Type, field.Interface())

		// 如果字段是结构体，递归解析
		if field.Kind() == reflect.Struct {
			fmt.Printf("Nested Struct in field: %s\n", fieldType.Name)
			ParseStruct(field.Interface()) // 递归调用
		}
	}
}
