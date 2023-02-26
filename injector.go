package inject

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type Args map[string]any

func Instanciate[T any]() (*T, error) {
	return InstanciateWithArgs[T](nil, false)
}

func InstanciateWithRemap[T any]() (*T, error) {
	return InstanciateWithArgs[T](nil, true)
}

func InstanciateWithArgs[T any](args Args, remap bool) (*T, error) {
	var t T
	Type := reflect.TypeOf(t)
	if Type.Kind() != reflect.Pointer {
		Type = reflect.PointerTo(Type)
	}

	v := reflect.New(Type.Elem())

	err := injectWithValueAndArgs(v, args, nil, remap)
	if err != nil {
		return nil, err
	}
	return v.Interface().(*T), nil
}

func InstanciateWithPositionalArgs[T any](args []any) (*T, error) {
	var t T
	Type := reflect.TypeOf(t)
	if Type.Kind() != reflect.Pointer {
		Type = reflect.PointerTo(Type)
	}

	v := reflect.New(Type.Elem())

	err := injectWithValueAndArgs(v, nil, args, false)
	if err != nil {
		return nil, err
	}
	return v.Interface().(*T), nil
}

func injectWithValueAndArgs(v reflect.Value, args Args, positional []any, doRemap bool) error {

	t := v.Type()

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("object must me a pointer to a struct")
	}

	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	var remap map[string]string
	var reverse map[string]string
	if doRemap {
		remap = make(map[string]string)
		reverse = make(map[string]string)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag
			injectFieldName := tag.Get("inject")
			if injectFieldName == "" {
				continue
			}
			if f.Name != injectFieldName {
				remap[f.Name] = injectFieldName
				reverse[injectFieldName] = f.Name
			}
		}
		if len(remap) == 0 {
			remap = nil
		}
	}

	for i := 0; i < t.NumField(); i++ {

		field := v.Field(i)

		f := t.Field(i)
		var fieldValue reflect.Value

		tag := f.Tag
		injectFieldName := tag.Get("inject")
		if injectFieldName == "" {
			continue
		}

		// TODO: implement option modes for injection with value of tag "inject"

		value := tag.Get("value")

		var dat = make(map[string]interface{})
		var slice []interface{}

		if value != "" {

			switch f.Type.Kind() {

			case reflect.String:
				fieldValue = reflect.ValueOf(value)

			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				value = strings.Trim(value, " ")
				var number, err = strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(number)

			case reflect.Float64:
				value = strings.Trim(value, " ")
				var number, err = strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(number)

			case reflect.Bool:
				value = strings.Trim(value, " ")
				value = strings.ToLower(value)
				b := !(value == "false" || value == "0" || value == "nil" || value == "none" || value == "null" || value == "")
				fieldValue = reflect.ValueOf(b)
			case reflect.Array | reflect.Slice:
				if err := json.Unmarshal([]byte(value), &slice); err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(slice)
			case reflect.Map:
				if err := json.Unmarshal([]byte(value), &dat); err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(dat)

			case reflect.Struct:

				if err := json.Unmarshal([]byte(value), &dat); err != nil {
					panic(err)
				}
				rf := v.Field(i)
				tt := fmt.Sprintf("%v |  %v", rf.Kind().String(), rf.Type())
				fmt.Println(tt)
				rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
				injectWithValueAndArgs(rf, dat, nil, doRemap)
				//InjectWithArgs(obj.Interface(), dat, true)

			default:
				// TODO: warning message - not supported types
			}

		}

		//path := f.Type.PkgPath() + "." + f.Type.Name()
		path := f.Type.String()
		if f.Type.PkgPath() == "" || f.Type.Name() == "" {
			path = fmt.Sprint(f.Type)
			if f.Type.Kind() == reflect.Pointer {
				path = path[1:]
			}
		}

		if injectFieldName == "struct" {
			descriptor := GetInjectable(path)
			if descriptor != nil {

				it := getInjectableType(descriptor.GetPath())

				if it != nil {
					fieldValue = reflect.New(it)
				}

				v, ok := findDeeperStruct(fieldValue)
				if ok {
					error := injectWithValueAndArgs(v, dat, slice, doRemap)
					if error != nil {
						return error
					}
				}

				if descriptor.Params != nil {
					argsValue := reflect.ValueOf(descriptor.Params)
					fillStructField(argsValue, fieldValue, reverse)
				}
			}
		}

		k := f.Type.Kind()
		if fieldValue.IsValid() ||
			(k != reflect.Struct && k != reflect.Pointer && k != reflect.UnsafePointer && k != reflect.Func) {
			if len(slice) > 0 {
				setFieldValue(f.Name, field, fieldValue, args, i, slice, remap)
			} else {
				setFieldValue(f.Name, field, fieldValue, args, i, positional, remap)
			}
		}

	}

	return nil
}

func setFieldValue(fieldName string, field reflect.Value, fieldValue reflect.Value, args Args, index int, positional []any, remap map[string]string) {
	if args != nil {
		if remap != nil {
			if v, ok := remap[fieldName]; ok {
				fieldName = v
			}
		}
		if v, ok := args[fieldName]; ok {
			switch field.Type().Kind() {
			case reflect.Int:
				x := v.(float64)
				n := int64(x)
				field.SetInt(n)
			case reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				field.SetInt(v.(int64))
			case reflect.Float32, reflect.Float64:
				field.SetFloat(v.(float64))
			case reflect.Bool:
				if v != nil {
					x, ok := v.(float64)
					if ok {
						field.SetBool(x != 0)
						break
					}
					s, ok := v.(string)
					if ok {
						field.SetBool(s != "" && s != "false")
					}
					n, ok := v.(int64)
					if ok {
						field.SetBool(n != 0)
						break
					}
					b, ok := v.(bool)
					if ok {
						field.SetBool(b)
						break
					}
				}
				field.SetBool(v != nil)
			default:
				field.Set(reflect.ValueOf(v))
			}
			return
		}
	}
	if positional != nil && index >= 0 && index < len(positional) {
		field.Set(reflect.ValueOf(positional[index]))
		return
	}
	if fieldValue.IsValid() {
		if field.Type().Kind() == reflect.Struct {
			field.Set(fieldValue.Elem())
		} else {
			field.Set(fieldValue)
		}
	}
}

func findDeeperStruct(v reflect.Value) (reflect.Value, bool) {
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		return findDeeperStruct(v.Elem())
	} else if v.Kind() == reflect.Struct {
		return v, true
	}
	return v, false
}

func InjectWithArgs(obj any, args Args, doRemap bool) error {
	v := reflect.ValueOf(obj)
	k := v.Kind()
	if k != reflect.Pointer {
		return errors.New("object must me a pointer to a struct")
	}
	return injectWithValueAndArgs(v, args, nil, doRemap)
}

func InjectWithPositionalArgs(obj any, args []any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer {
		fmt.Printf("--- kind:%v", v.Kind())
		return errors.New("object must me a pointer to a struct.")
	}

	return injectWithValueAndArgs(v, nil, args, false)
}

func Inject(obj any) error {
	return InjectWithArgs(obj, nil, false)
}

func fillStructField(data reflect.Value, fieldValue reflect.Value, remap map[string]string) {

	kind := data.Kind()
	elem := fieldValue.Elem()
	switch kind {
	case reflect.Array | reflect.Slice:
		for i := 0; i < data.Len(); i++ {

			elemField := elem.Field(i)
			v := data.Index(i)
			if elemField.CanSet() {
				elemField.Set(v.Elem())
			}

		}
	case reflect.Map:
		iter := data.MapRange()
		for iter.Next() {

			k := iter.Key().Elem().String()
			if remap != nil {
				k = remap[k]
			}
			v := iter.Value()
			elemField := elem.FieldByName(k)

			if elemField.CanSet() {
				elemField.Set(v.Elem())
			}

		}
	default:
		if elem.NumField() > 0 {
			elemField := elem.Field(0)
			if elemField.CanSet() {
				elemField.Set(data)
			}
		}
	}
}
