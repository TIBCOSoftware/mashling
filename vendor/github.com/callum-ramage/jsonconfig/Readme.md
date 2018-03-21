#jsonconfig#

This package provides a convenient mechanism for using a json structure as a configuration file with the added benefit of allowing //comments.

## [GoDoc](http://godoc.org/github.com/callum-ramage/jsonconfig) ##

##Installation and Usage##

To install simply call

	go get github.com/callum-ramage/jsonconfig

Usage of the package couldn't be simpler

	package main

	import (
		"github.com/callum-ramage/jsonconfig"
		"fmt"
	)

	func main() {
		config, err := jsonconfig.LoadAbstract("./configs/ExampleConfig.conf", "")

		if err != nil {
			return
		}

		fmt.Println(config["example_string"].Str)
		fmt.Println(config["example_array"].Arr[0].Str)
		fmt.Println(config["example_object"].Obj["example_number"].Num)
		fmt.Println(config["example_object"].Obj["example_number"].Int)
		//Or
		fmt.Println(config["example_string"].Str)
		fmt.Println(config["example_array.0"].Str)
		fmt.Println(config["example_object.example_number"].Num)
		fmt.Println(config["example_object.example_number"].Int)
	}

Outputs

	string value
	array value 0
	5.3
	5

Where `./configs/ExampleConfig.conf` is

	{
		"example_string": "string value",
		"example_array": [
			"array value 0"
		],
		"example_object": {
			"example_number": 5.3
		}
	}

For a more detailed example that includes defining default values, have a look at [jsonconfig_test.go](jsonconfig_test.go) or the [GoDoc](http://godoc.org/github.com/callum-ramage/jsonconfig)
