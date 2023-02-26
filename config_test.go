package inject

import (
	"fmt"
	"testing"
)

type iMessageTester interface {
	Test()
}

type iMessagePrinter interface {
	Print()
	GetMessage() string
}

const messagePrinterA_MSG string = "This message is from messagePrinterA"

type messagePrinterA struct {
}

func (mp *messagePrinterA) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *messagePrinterA) GetMessage() string {
	return messagePrinterA_MSG
}

type messagePrinterB struct {
	Message string
}

func (mp *messagePrinterB) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *messagePrinterB) GetMessage() string {
	return mp.Message
}

type messagePrinterC struct {
	Message string
	Count   int
}

func (mp *messagePrinterC) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *messagePrinterC) GetMessage() string {
	var result string
	for i := 0; i < mp.Count; i++ {
		result += fmt.Sprintf("%s: %d\n", mp.Message, i)
	}
	return result
}

type messagePrinterD struct {
	Message1 string
	Message2 any
	Flag     bool
	Message3 any
	Value    float64

	SubPrinter messagePrinterB `inject:"true"`
}

func (mp *messagePrinterD) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *messagePrinterD) GetMessage() string {
	var result string
	result += fmt.Sprintf("%v\n%v\n%v\n%v\n%v\n", mp.Message1, mp.Message2, mp.Flag, mp.Message3, mp.Value)
	result += fmt.Sprintln("Message from messagePrinterD.SubPrinter")
	result += mp.SubPrinter.GetMessage()
	return result
}

type messagePrinterE struct {
	Message1 string `inject:"true" value:"\"direct injection on field Message1\""`
	Message2 any
	Flag     bool
	Message3 any
	Value    float64

	SubPrinter *messagePrinterB `inject:"true"`
}

func (mp *messagePrinterE) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *messagePrinterE) GetMessage() string {
	result := fmt.Sprintf("%v\n%v\n%v\n%v\n%v\n", mp.Message1, mp.Message2, mp.Flag, mp.Message3, mp.Value)
	if mp.SubPrinter != nil {
		result += fmt.Sprintln("Message from messagePrinterD.SubPrinter")
		result += mp.SubPrinter.GetMessage()
	} else {
		result += fmt.Sprintln("SubPrinter at messagePrinterE = nil")
	}
	return result
}

type printerContainer struct {
	Printer iMessagePrinter `inject:"struct"`
}

func init_instances() {

	// interface registration (mode 2)
	AddWrappedInterface(InterfaceWrapper[iMessagePrinter]{})

	// injectable structs registration
	AddInjectable(messagePrinterA{})
	AddInjectable(messagePrinterB{})
	AddInjectable(messagePrinterC{})
	AddInjectable(messagePrinterD{})
	AddInjectable(messagePrinterE{})

	// injected structs registration
	AddFactory(&printerContainer{}, false)

}

func TestImportConfig(t *testing.T) {

	t.Log("initializing TestImportConfig")
	init_instances()

	file := "test_files/injection-config.local.yaml"
	t.Logf("\n\nimporting config file %s\n", file)
	ImportConfig(file)

	description := config.GetInterface("inject.iMessagePrinter")
	t.Log(description)

	if description == nil {
		t.Fatalf("Interface name. Expected not nil, got %v", description)
	}

	const iName = "iMessagePrinter"

	if description.Name != iName {
		t.Fatalf("Interface name. Expected %s, got %s", iName, description.Name)
	}

	const injName = "messagePrinterA"

	if description.Injectable != injName {
		t.Fatalf("Interface Injectable. Expected %s, got %s", injName, description.Injectable)
	}

	const iPkg = "inject"

	if description.Package != iPkg {
		t.Fatalf("Interface Package. Expected %s, got %s", iPkg, description.Package)
	}

	t.Log("TestImportConfig ok")

}

func TestGetInstanceWithImportedCongif(t *testing.T) {

	init_instances()
	pc := GetInstance[printerContainer](nil)

	if pc.Printer != nil {
		t.Fatalf("TestGetInstanceWithImportedCongif(). Expected %v, got %v", nil, pc.Printer)
	}

	file := "test_files/injection-config.local.yaml"
	ImportConfig(file)

	printerContainer := printerContainer{}
	Inject(&printerContainer)

	if printerContainer.Printer == nil {
		t.Fatalf("TestGetInstanceWithImportedCongif(). Expected %v, got %v", nil, pc.Printer)
	}

	message := printerContainer.Printer.GetMessage()

	if message != messagePrinterA_MSG {
		t.Fatalf("TestGetInstanceWithImportedCongif(). Expected %v, got %v", message, messagePrinterA_MSG)
	}

	Inject(pc)

	if pc.Printer == nil {
		t.Fatalf("TestGetInstanceWithImportedCongif(). Expected %v, got %v", nil, pc.Printer)
	}

	message = pc.Printer.GetMessage()

	if message != messagePrinterA_MSG {
		t.Fatalf("TestGetInstanceWithImportedCongif(). Expected %v, got %v", message, messagePrinterA_MSG)
	}

}
