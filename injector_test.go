package inject

import (
	"testing"
)

func TestInjectValuesAsMappedArgs(t *testing.T) {

	type testStruct struct {
		StringField  string `inject:"true"`
		IntField     int    `inject:"true"`
		BooleanField bool   `inject:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	var args map[string]any = map[string]any{
		"BooleanField": boolValue,
		"IntField":     intValue,
		"StringField":  stringValue,
	}

	testObject := testStruct{}

	err := InjectWithArgs(&testObject, args, false)

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if testObject.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. testStruct.BooleanField = %v, expected %v", testObject.BooleanField, boolValue)
	}

	if testObject.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.IntField = %v, expected %v", testObject.IntField, intValue)
	}

	if testObject.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.StringField = %v, expected %v", testObject.StringField, stringValue)
	}

}

func TestInjectValuesAsPositionalArgs(t *testing.T) {

	type testStruct struct {
		StringField  string `inject:"true"`
		IntField     int    `inject:"true"`
		BooleanField bool   `inject:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	var args []any = []any{
		stringValue,
		intValue,
		boolValue,
	}

	testObject := testStruct{}

	err := InjectWithPositionalArgs(&testObject, args)

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if testObject.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. testStruct.BooleanField = %v, expected %v", testObject.BooleanField, boolValue)
	}

	if testObject.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.IntField = %v, expected %v", testObject.IntField, intValue)
	}

	if testObject.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.StringField = %v, expected %v", testObject.StringField, stringValue)
	}

}

func TestIntanciate(t *testing.T) {

	type testStruct struct {
		StringField  string `inject:"true"`
		IntField     int    `inject:"true"`
		BooleanField bool   `inject:"true"`
	}

	object, err := Instanciate[testStruct]()

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	t.Logf("object instanciated: %v", object)

}

func TestInstanciateWithMappedArgs(t *testing.T) {

	type testStruct struct {
		StringField  string `inject:"true"`
		IntField     int    `inject:"true"`
		BooleanField bool   `inject:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	var args map[string]any = map[string]any{
		"BooleanField": boolValue,
		"IntField":     intValue,
		"StringField":  stringValue,
	}

	testObject, err := InstanciateWithArgs[testStruct](args, false)

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if testObject.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. testStruct.BooleanField = %v, expected %v", testObject.BooleanField, boolValue)
	}

	if testObject.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.IntField = %v, expected %v", testObject.IntField, intValue)
	}

	if testObject.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.StringField = %v, expected %v", testObject.StringField, stringValue)
	}

}

func TestInstanciateWithPositionalArgs(t *testing.T) {

	type testStruct struct {
		StringField  string `inject:"true"`
		IntField     int    `inject:"true"`
		BooleanField bool   `inject:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	var args []any = []any{
		stringValue,
		intValue,
		boolValue,
	}

	testObject, err := InstanciateWithPositionalArgs[testStruct](args)

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if testObject.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. testStruct.BooleanField = %v, expected %v", testObject.BooleanField, boolValue)
	}

	if testObject.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.IntField = %v, expected %v", testObject.IntField, intValue)
	}

	if testObject.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.StringField = %v, expected %v", testObject.StringField, stringValue)
	}

}

func TestInstanciateWithInitialValues(t *testing.T) {

	type testStruct struct {
		NonParametrizedField string
		StringField          string `inject:"stringField" value:"teste"`
		IntField             int    `inject:"intField" value:"12"`
		BooleanField         bool   `inject:"booleanField" value:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	testObject, err := Instanciate[testStruct]()

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if testObject.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. testStruct.BooleanField = %v, expected %v", testObject.BooleanField, boolValue)
	}

	if testObject.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.IntField = %v, expected %v", testObject.IntField, intValue)
	}

	if testObject.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. testObject.StringField = %v, expected %v", testObject.StringField, stringValue)
	}

	type CmposedStruct struct {
		object testStruct `inject:"object" value:"{\"stringField\": \"teste2\",\"intField\": 13,\"booleanField\": 0}"`
	}

	testObject2, err := InstanciateWithRemap[CmposedStruct]()

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if testObject2.object.IntField != 13 {
		t.Fatalf("tested object int field: injection failed. testObject2.Object.IntField = %v, expected %v", testObject2.object.IntField, 13)
	}

	if testObject2.object.BooleanField != false {
		t.Fatalf("tested object boolean field: injection failed. testObject2.Object.BooleanField = %v, expected %v", testObject2.object.BooleanField, false)
	}
}

func TestInstanciateCompoundStructWithInitialValues(t *testing.T) {

	type MemberStruct struct {
		NonParametrizedField string
		StringField          string `inject:"stringField" value:"teste"`
		IntField             int    `inject:"intField" value:"12"`
		BooleanField         bool   `inject:"booleanField" value:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	type Container struct {
		Id     int          `inject:"id" value:"100"`
		Object MemberStruct `inject:"object"`
	}

	container, err := Instanciate[Container]()

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if container.Object.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. container.Object.BooleanField = %v, expected %v", container.Object.BooleanField, boolValue)
	}

	if container.Object.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. container.Object.IntField = %v, expected %v", container.Object.IntField, intValue)
	}

	if container.Object.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. container.Object.StringField = %v, expected %v", container.Object.StringField, stringValue)
	}

}

func TestInstanciateCompoundStructWithPointer(t *testing.T) {

	type MemberStruct struct {
		NonParametrizedField string
		StringField          string `inject:"stringField" value:"teste"`
		IntField             int    `inject:"intField" value:"12"`
		BooleanField         bool   `inject:"booleanField" value:"true"`
	}

	intValue := 12
	stringValue := "teste"
	boolValue := true

	type Container struct {
		Object *MemberStruct `inject:"object"`
	}

	container, err := Instanciate[Container]()

	if err != nil {
		t.Fatalf("Error calling injectWithValueAndArgs: %v", err)
	}

	if container.Object.BooleanField != boolValue {
		t.Fatalf("tested object boolean field: injection failed. container.Object.BooleanField = %v, expected %v", container.Object.BooleanField, boolValue)
	}

	if container.Object.IntField != intValue {
		t.Fatalf("tested object boolean field: injection failed. container.Object.IntField = %v, expected %v", container.Object.IntField, intValue)
	}

	if container.Object.StringField != stringValue {
		t.Fatalf("tested object boolean field: injection failed. container.Object.StringField = %v, expected %v", container.Object.StringField, stringValue)
	}

}
