package inject

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Args map[string]any

func Instanciate[T any](Type reflect.Type) (*T, error) {
	return InstanciateWithArgs[T](Type, nil)
}

func InstanciateWithArgs[T any](Type reflect.Type, args Args) (*T, error) {
	if Type.Kind() != reflect.Pointer {
		Type = reflect.PointerTo(Type)
	}

	v := reflect.New(Type.Elem())

	err := injectWithValueAndArgs(v, args, nil)
	if err != nil {
		return nil, err
	}
	return v.Interface().(*T), nil
}

func injectWithValueAndArgs(v reflect.Value, args Args, positional []any) error {

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

	for i := 0; i < t.NumField(); i++ {

		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		f := t.Field(i)
		var fieldValue reflect.Value

		tag := f.Tag
		value := tag.Get("inject")
		if value == "" {
			continue
		}

		value = tag.Get("value")

		var dat map[string]interface{}
		var slice []interface{}

		if value != "" {

			switch f.Type.Kind() {

			case reflect.String:
				var str string
				if err := json.Unmarshal([]byte(value), &str); err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(str)

			case reflect.Int | reflect.Int16 | reflect.Int32 | reflect.Int64 | reflect.Int8 | reflect.Uint | reflect.Uint16 | reflect.Uint32 | reflect.Uint64 | reflect.Uint8:
				var number int64
				if err := json.Unmarshal([]byte(value), &number); err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(number)

			case reflect.Float64:
				var number float64
				if err := json.Unmarshal([]byte(value), &number); err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(number)

			case reflect.Bool:
				var boolean bool
				if err := json.Unmarshal([]byte(value), &boolean); err != nil {
					panic(err)
				}
				fieldValue = reflect.ValueOf(boolean)
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
			}

		}

		path := f.Type.PkgPath() + "." + f.Type.Name()
		if f.Type.PkgPath() == "" || f.Type.Name() == "" {
			path = fmt.Sprint(f.Type)
			if f.Type.Kind() == reflect.Pointer {
				path = path[1:]
			}
		}
		descriptor := GetInjectable(path)
		if descriptor != nil {

			it := getInjectableType(descriptor.GetPath())

			if it != nil {
				fieldValue = reflect.New(it)
			}

			v, ok := findDeeperStruct(fieldValue)
			if ok {
				error := injectWithValueAndArgs(v, dat, slice)
				if error != nil {
					return error
				}
			}

			if descriptor.Params != nil {
				argsValue := reflect.ValueOf(descriptor.Params)
				fillStructField(argsValue, fieldValue)
			}
		}

		if fieldValue.IsValid() {
			setFieldValue(f.Name, field, fieldValue, args, i, slice)
		}

	}

	return nil
}

func setFieldValue(fieldName string, field reflect.Value, fieldValue reflect.Value, args Args, index int, positional []any) {
	if args != nil {
		if v, ok := args[fieldName]; ok {
			field.Set(reflect.ValueOf(v))
			return
		}
	}
	if positional != nil && index >= 0 && index < len(positional) {
		field.Set(reflect.ValueOf(positional[index]))
		return
	}
	if field.Type().Kind() == reflect.Struct {
		field.Set(fieldValue.Elem())
	} else {
		field.Set(fieldValue)
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

func InjectWithArgs(obj any, args Args) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer {
		return errors.New("object must me a pointer to a struct")
	}

	return injectWithValueAndArgs(v, args, nil)
}

func Inject(obj any) error {
	return InjectWithArgs(obj, nil)
}

func fillStructField(data reflect.Value, fieldValue reflect.Value) {

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
