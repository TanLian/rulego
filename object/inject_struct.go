package object

import (
	"fmt"
	"reflect"
)

type InjectStruct struct {
	Value           reflect.Value
	methodToArgType map[string][]reflect.Type
}

func (s *InjectStruct) object() {}

func (s *InjectStruct) Type() Type {
	return TypeStruct
}

func (s *InjectStruct) GetValue() any {
	return s.Value.Interface()
}

func (s *InjectStruct) GetFieldValue(field string) Object {
	stt := s.Value
	if stt.Kind() == reflect.Pointer {
		stt = stt.Elem()
	}
	val := stt.FieldByName(field)
	if val.IsValid() {
		return ToObject(val)
	}
	return Null
}

func (s *InjectStruct) SetField(fieldName string, value any) {
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

func (s *InjectStruct) lazyLoadMethodType(name string) ([]reflect.Type, bool) {
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
			fmt.Printf("Nested InjectStruct in field: %s\n", fieldType.Name)
			ParseStruct(field.Interface()) // 递归调用
		}
	}
}
