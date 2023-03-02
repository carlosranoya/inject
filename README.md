# Inject
## Dependency injection package for Go


***Inject*** is a runtime dependency injection package for Go, based on struct tags and built on top of Go reflect package.

## Some Features

- Dependency injection rules in external yaml files
- Custom struct factories, with singleton option
- Default struct values, using struct tags and/or external file

## Installation

Make sure you have Go installed ([download](https://go.dev/dl/)). Version 1.18 or higher is required.
Initialize your project by creating a folder and running ``` go mod init github.com/your/repo ``` inside the folder. Then install ***Inject*** with the ``` go get ```  command:

```sh
go get -u github.com/carlosranoya/inject
```

## Example of Use 1

Create a main.go file and choose a package name. In this example, I choose package main.

```sh
package main
```

Import ***Inject*** pakage and other dependencies.

```sh
import (
    "fmt"
    "github.com/carlosranoya/inject"
)
```

Create an interface that will be injected by some struct instance.

```sh
type TestInterface interface {
    Test()
}
```

Create a struct that implements the interface.

```sh
type TestStruct struct {
    Message string
}
func (t *TestStruct) Test() {
    fmt.Println(t.Message)
}
```

At the *main* or *init* function, use ***Inject*** to register the interface and the injected struct.

```sh
func main() {

    fmt.Println("testing inject")

    inject.AddInterface[TestInterface]()

    inject.AddInjectable[TestStruct]()
```

Use ***Inject*** to register an yaml file that describe the relationship between interfaces and structs (details about this file below).

```sh
    inject.ImportConfig("config_1.yaml")
```

Create a struct that contains the interface and inject the interface field with the struct defined at "config_1.yaml" file.

```sh
    type TestContainer struct {
        Tester TestInterface `inject:"struct"`
    }
    container := TestContainer{}
    inject.Inject(&container)
```
Note that is necessary a tag definition to the struct's field intended to be injected. Make sure you pass a pointer to the *Inject()* function.

Finally, call the interface function to test the dependency injection and close the main function.

```sh
    container.Tester.Test()
}
```

Let's see the config_1.yaml file. It contains a section for "structs", identified as "injectables":

```sh
injectables:
  - name: TestStruct
    package: main
    params: 
      Message: "This message was defined at config_1.yaml"
```
This is a list with structs that will be injected at interface fields. Note that optionally there's a params field where we define values that will be set to struct fields. In this example, the struct TestStruct owns a string field named Message.
Next, the interfaces section:
```sh
interfaces:
  - name: TestInterface
    injectable: TestStruct
    package: main
```
Here, we define the relationship between interfaces and injectable structs.
The "package" parameter is fundamental, and it corresponds to the name of the package where interfaces and structs where defined.
That's it.

## Example of Use 2

This is very similar of last example, but instead of using a container struct (TestContainer) to be filled with a struct, we will define an inline interface variable, and it'll be instanciated as the struct defined at congif_1.yaml file.

```sh
    var i TestInterface
    i, err = InstaciateInjected[TestInterface]()
    i.Test()
```



## Other Uses

Inject comes with another utilities.
It may be used to define standard values for structs.

For example:
```sh
type TestStruct struct {
	NonParametrizedField string
	StringField          string `inject:"stringField" value:"teste"`
	IntField             int    `inject:"intField" value:"12"`
	BooleanField         bool   `inject:"booleanField" value:"true"`
}
```
At this case, we defined 3 fields that will be initializated with values defined at the struct's tags. The name after "inject:" must be any unique, not empty string, and is not relevant for this example.
Next, create an instance using this ***Inject*** funcion:

```sh
testObject, err := inject.Instanciate[TestStruct]()
```
The names defined with ***Inject*** tags may be used to instanciate objects with map[string]any arguments.
```sh
type Developer struct {
   NonParametrizedField string
   StringField          string `inject:"name" value:"teste"`
   IntField             int    `inject:"age" value:"12"`
   BooleanField         bool   `inject:"likesGo" value:"false"`
}

var args map[string]any = map[string]any{
   "likesGo": true,
   "age":     27,
   "name": "Bob",
}

developer, err := inject.InstanciateWithArgs[Developer](args, false)
```


[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [go-install]: <https://go.dev/dl/>
   [git-repo-url]: <https://github.com/joemccann/dillinger.git>
   [john gruber]: <http://daringfireball.net>
   [df1]: <http://daringfireball.net/projects/markdown/>
   [markdown-it]: <https://github.com/markdown-it/markdown-it>
   [Ace Editor]: <http://ace.ajax.org>
   [node.js]: <http://nodejs.org>
   [Twitter Bootstrap]: <http://twitter.github.com/bootstrap/>
   [jQuery]: <http://jquery.com>
   [@tjholowaychuk]: <http://twitter.com/tjholowaychuk>
   [express]: <http://expressjs.com>
   [AngularJS]: <http://angularjs.org>
   [Gulp]: <http://gulpjs.com>

   [PlDb]: <https://github.com/joemccann/dillinger/tree/master/plugins/dropbox/README.md>
   [PlGh]: <https://github.com/joemccann/dillinger/tree/master/plugins/github/README.md>
   [PlGd]: <https://github.com/joemccann/dillinger/tree/master/plugins/googledrive/README.md>
   [PlOd]: <https://github.com/joemccann/dillinger/tree/master/plugins/onedrive/README.md>
   [PlMe]: <https://github.com/joemccann/dillinger/tree/master/plugins/medium/README.md>
   [PlGa]: <https://github.com/RahulHP/dillinger/blob/master/plugins/googleanalytics/README.md>