package inject

import (
	"reflect"
	"testing"
)

func TestCheckType(t *testing.T) {
	var i int = 10
	ok := checkType[int](reflect.TypeOf(i))
	if !ok {
		t.Fatalf("checkType[int](reflect.TypeOf(%v) = %v, want checkType[int](reflect.TypeOf(%v) = %v", i, ok, i, !ok)
	}

	ok = checkType[int8](reflect.TypeOf(i))
	if ok {
		t.Fatalf("checkType[int8](reflect.TypeOf(%v) = %v, want checkType[int8](reflect.TypeOf(%v) = %v", i, !ok, i, ok)
	}
}

func TestAddInterface(t *testing.T) {

	type TestInterface interface {
		Test()
	}

	var I *TestInterface
	AddInterface(I)

	l := len(interfaces)

	if l != 1 {
		t.Fatalf("number of items: %v, want 1", l)
	}

	for a, b := range interfaces {
		t.Logf("key: %v, value:%v", a, b)
	}

}

func TestAddWrappedInterface(t *testing.T) {

	type TestInterface interface {
		Test()
	}

	AddWrappedInterface(InterfaceWrapper[TestInterface]{})

	l := len(interfaces)

	if l != 1 {
		t.Fatalf("number of items: %v, want 1", l)
	}

	for a, b := range interfaces {
		t.Logf("key: %v, value:%v", a, b)
	}

}

func TestAddAndGetInjectables(t *testing.T) {

	type Injectable struct {
		A int
		B string
	}

	var I Injectable
	AddInjectable(I)

	l := len(injectables)

	if l != 1 {
		t.Fatalf("number of items: %v, want 1", l)
	}

	for a, b := range injectables {
		t.Logf("key: %v, value:%v", a, b)
	}

	T := getInjectableType("inject.Injectable")

	if T != reflect.TypeOf(I) {
		t.Fatalf("wrong type of struct %v, got %v", l, T)
	}

}

func TestAddAndGetFactories(t *testing.T) {

	type Factory struct {
		A int
		B string
	}

	F := Factory{1, "test F"}
	AddFactory(&F, true)

	type Factory2 struct {
		A int
		B string
	}

	G := Factory2{2, "test G"}
	AddFactory(&G, false)

	l := len(factories)

	if l != 2 {
		t.Fatalf("number of items: %v, want 1", l)
	}

}

func TestAddAndGetInstance(t *testing.T) {

	type Factory struct {
		A int
		B string
	}

	F := Factory{1, "test F"}
	AddFactory(&F, true)

	type Factory2 struct {
		A int
		B string
	}

	G := Factory2{2, "test G"}
	AddFactory(&G, false)

	H := GetInstance[Factory](nil)

	if reflect.TypeOf(*H) != reflect.TypeOf(F) {
		t.Fatalf("Differente types: %T and %T. Expected the same types.", *H, F)
	}

	if reflect.TypeOf(*H) == reflect.TypeOf(G) {
		t.Fatalf("Equal types: %T and %T. Expected different types.", *H, G)
	}

	//singleton test

	J := GetInstance[Factory](nil)

	if J != H {
		t.Fatalf("Factory id build by a singleton factory. %v and %v should be the same pointer.", J, H)
	}

	// Reset test

	resetFactories()

	J = GetInstance[Factory](nil)

	if J == H {
		t.Fatalf("Factory factory was reseted. %v and %v should be the different pointers.", J, H)
	}

}


