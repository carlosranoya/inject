package inject

import (
	"fmt"
	"testing"
)

type IMessageTester interface {
	Test()
}

type IMessagePrinter interface {
	Print()
	GetMessage() string
}

const MessagePrinterA_MSG string = "This message is from MessagePrinterA"

type MessagePrinterA struct {
}

func (mp *MessagePrinterA) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *MessagePrinterA) GetMessage() string {
	return MessagePrinterA_MSG
}

type MessagePrinterB struct {
	Message string
}

func (mp *MessagePrinterB) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *MessagePrinterB) GetMessage() string {
	return mp.Message
}

type MessagePrinterC struct {
	Message string
	Count   int
}

func (mp *MessagePrinterC) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *MessagePrinterC) GetMessage() string {
	var result string
	for i := 0; i < mp.Count; i++ {
		result += fmt.Sprintf("%s: %d\n", mp.Message, i)
	}
	return result
}

type MessagePrinterD struct {
	Message1 string
	Message2 any
	Flag     bool
	Message3 any
	Value    float64

	SubPrinter MessagePrinterB `inject:"true"`
}

func (mp *MessagePrinterD) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *MessagePrinterD) GetMessage() string {
	var result string
	result += fmt.Sprintf("%v\n%v\n%v\n%v\n%v\n", mp.Message1, mp.Message2, mp.Flag, mp.Message3, mp.Value)
	result += fmt.Sprintln("Message from MessagePrinterD.SubPrinter")
	result += mp.SubPrinter.GetMessage()
	return result
}

type MessagePrinterE struct {
	Message1 string `inject:"true" value:"\"direct injection on field Message1\""`
	Message2 any
	Flag     bool
	Message3 any
	Value    float64

	SubPrinter *MessagePrinterB `inject:"true"`
}

func (mp *MessagePrinterE) Print() {
	fmt.Println(mp.GetMessage())
}
func (mp *MessagePrinterE) GetMessage() string {
	result := fmt.Sprintf("%v\n%v\n%v\n%v\n%v\n", mp.Message1, mp.Message2, mp.Flag, mp.Message3, mp.Value)
	if mp.SubPrinter != nil {
		result += fmt.Sprintln("Message from MessagePrinterD.SubPrinter")
		result += mp.SubPrinter.GetMessage()
	} else {
		result += fmt.Sprintln("SubPrinter at MessagePrinterE = nil")
	}
	return result
}

type PrinterContainer struct {
	Printer IMessagePrinter `inject:"struct"`
}

func init_instances() {

	// interface registration (mode 2)
	AddWrappedInterface(InterfaceWrapper[IMessagePrinter]{})

	// injectable structs registration
	AddInjectable(MessagePrinterA{})
	AddInjectable(MessagePrinterB{})
	AddInjectable(MessagePrinterC{})
	AddInjectable(MessagePrinterD{})
	AddInjectable(MessagePrinterE{})

	// injected structs registration
	AddFactory(&PrinterContainer{}, false)

}

func TestImportConfig(t *testing.T) {

	init_instances()

	file := "test_files/injection-config.local.yaml"
	t.Logf("\n\nimporting config file %s\n", file)
	ImportConfig(file)

	description := config.GetInterface("inject.IMessagePrinter")
	t.Log(description)

	if description == nil {
		t.Fatalf("Interface name. Expected not nil, got %v", description)
	}

	const iName = "IMessagePrinter"

	if description.Name != iName {
		t.Fatalf("Interface name. Expected %s, got %s", iName, description.Name)
	}

	const injName = "MessagePrinterA"

	if description.Injectable != injName {
		t.Fatalf("Interface Injectable. Expected %s, got %s", injName, description.Injectable)
	}

	const iPkg = "inject"

	if description.Package != iPkg {
		t.Fatalf("Interface Package. Expected %s, got %s", iPkg, description.Package)
	}

}
