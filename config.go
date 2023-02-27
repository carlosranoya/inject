package inject

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var config configData = configData{}

type componentPath struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
}

func (path *componentPath) GetPath() string {
	return path.Package + "." + path.Name
}

type factoryDescription struct {
	componentPath `yaml:",inline"`
	IsSingleton   bool `yaml:"is-singleton"`
}

type injectableDescription struct {
	componentPath `yaml:",inline"`
	Factory       string `yaml:"factory"`
	Params        any    `yaml:"params"`
	InjectMode    string `yaml:"mode"` // values: auto, interface, factory
}

type interfaceDescription struct {
	componentPath `yaml:",inline"`
	Injectable    string `yaml:"injectable"`
}

type configData struct {
	Factories   []factoryDescription    `yaml:"factories"`
	Injectables []injectableDescription `yaml:"injectables"`
	Interfaces  []interfaceDescription  `yaml:"interfaces"`
}

func GetInjectable(fieldType string) *injectableDescription {
	inter := config.GetInterface(fieldType)
	if inter == nil {
		inj := config.GetInjectable(fieldType)
		if inj != nil && inj.InjectMode == "auto" {
			return inj
		}
		return nil
	}
	if config.Injectables == nil {
		return nil
	}
	for _, inj := range config.Injectables {
		if inj.Name == inter.Injectable {
			return &inj
		}
	}
	inj := config.GetInjectable(fieldType)
	if inj != nil && inj.InjectMode == "auto" {
		return inj
	}
	return nil
}

func (data *configData) GetInterface(name string) *interfaceDescription {
	if data.Interfaces == nil {
		return nil
	}
	for _, i := range data.Interfaces {
		if i.Name == name || i.GetPath() == name {
			return &i
		}
	}
	return nil
}

func (data *configData) GetInjectable(name string) *injectableDescription {
	if data.Injectables == nil {
		return nil
	}
	for _, i := range data.Injectables {
		if i.Name == name || i.GetPath() == name {
			return &i
		}
	}
	return nil
}

func ImportConfig(filename string) {
	resetFactories()
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(content, &config)

	if err != nil {
		log.Fatal(err)
	}

}
