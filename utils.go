package inject

import (
	"fmt"
	"reflect"
)

func Describe(v any, tabs string) {
	describe(reflect.ValueOf(v), tabs)
}

func describe(v reflect.Value, tabs string) {

	switch v.Kind() {
	case reflect.Struct:
		fmt.Print(tabs + v.Type().Name() + " {")
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fmt.Printf("\n%s  %v: ", tabs, v.Type().Field(i).Name)
			describe(f, tabs+"  ")
		}
		fmt.Print("\n" + tabs + "}")
	case reflect.Interface:
		fmt.Print(tabs)
		fmt.Print(v)
		if v.IsNil() {
			return
		}
	case reflect.Pointer:
		fmt.Print("*")
		describe(v.Elem(), tabs)
	case reflect.Array | reflect.Slice:

		if v.IsNil() {
			fmt.Print(v)
			return
		}

		fmt.Println("[")
		for i := 0; i < v.Len(); i++ {
			sep := ","
			if i == v.Len()-1 {
				sep = ""
			}
			describe(v.Index(i), tabs+"  ")
			fmt.Println(sep)
		}
		fmt.Print(tabs + "]")

	case reflect.Map:

		if v.IsNil() {
			fmt.Print(v)
			return
		}
		fmt.Printf("%smap[%v] {", tabs, v.Kind())
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			f := v.MapIndex(key)
			fmt.Printf("\n%s  %v: ", tabs, key)
			describe(f, tabs+"  ")
		}
		fmt.Print("\n" + tabs + "}")

	case reflect.Func:
		fmt.Print(tabs)
		fmt.Print(v)
	case reflect.Chan:
		fmt.Print(tabs)
		fmt.Print(v)
		if v.IsNil() {
			return
		}
	case reflect.Invalid:
		fmt.Print("'undefined'")
		return
	case reflect.String:
		fmt.Print("\"", v, "\"")
	default:
		fmt.Print(v)
	}

}
