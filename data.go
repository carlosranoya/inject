package inject

import (
	"fmt"
	"reflect"
)

type InterfaceWrapper[T any] struct {
	pointer *T
}

type injectedFactory[T any] struct {
	IsSingleton bool
	instance    *T
}

type iResetable interface {
	Reset()
}

func (factory *injectedFactory[T]) GetInstance() *T {
	return factory.getInstanceWithArgs(nil)
}
func (factory *injectedFactory[T]) getInstanceWithArgs(args Args) *T {
	if factory.IsSingleton {
		if factory.instance == nil {
			instance, err := InstanciateWithArgs[T](args, false)
			if err == nil {
				factory.instance = instance
			}
		}
		return factory.instance
	}
	instance, err := InstanciateWithArgs[T](args, false)
	if err == nil {
		return instance
	}
	return nil
}
func (factory *injectedFactory[T]) Reset() {
	factory.instance = nil
}

func ResetData() {
	factories = make(map[reflect.Type]iResetable)

	interfaces = make(map[string]reflect.Type)

	injectables = make(map[string]reflect.Type)
}

var factories map[reflect.Type]iResetable = make(map[reflect.Type]iResetable)

var interfaces map[string]reflect.Type = make(map[string]reflect.Type)

var injectables map[string]reflect.Type = make(map[string]reflect.Type)

func resetFactories() {
	for _, v := range factories {
		v.Reset()
	}
}

func AddFactory[T any](obj *T, IsSingleton bool) error {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	factory := injectedFactory[T]{IsSingleton: IsSingleton}
	factories[t] = &factory
	return nil
}

func checkType[T any](Type reflect.Type) bool {
	v := reflect.New(Type)
	_, ok := v.Interface().(*T)
	return ok
}

func GetInstance[T any](args Args) *T {
	var f iResetable
	for k, v := range factories {
		ok := checkType[T](k)
		if ok {
			f = v
			break
		}
	}
	if f == nil {
		return nil
	}
	factory, ok := f.(*injectedFactory[T])
	if ok {
		return factory.getInstanceWithArgs(args)
	}
	return nil
}

func AddInterface(pointer any) {
	t := reflect.TypeOf(pointer).Elem()
	name := fmt.Sprintf("%v", reflect.TypeOf(pointer).Elem())
	interfaces[name] = t
}

func AddWrappedInterface[T any](wrapper InterfaceWrapper[T]) {
	t := reflect.TypeOf(wrapper.pointer).Elem()
	name := fmt.Sprintf("%v", reflect.TypeOf(wrapper.pointer).Elem())
	interfaces[name] = t
}

func AddInjectable(obj any) {
	t := reflect.TypeOf(obj)
	name := fmt.Sprintf("%v", t)
	injectables[name] = t
}

func getInjectableType(name string) reflect.Type {
	return injectables[name]
}
