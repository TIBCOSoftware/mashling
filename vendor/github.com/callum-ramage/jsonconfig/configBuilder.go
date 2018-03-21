//jsonconfig contains a set of useful structures for accessing json data from a
//configuration file. It uses a pre-processor that removes //comments from the file
//before parsing it.
//
//  package main
//
//  import (
//    "github.com/callum-ramage/jsonconfig"
//    "fmt"
//  )
//
//  func main() {
//    config, err := jsonconfig.LoadAbstract("./configs/ExampleConfig.conf", "")
//
//    if err != nil {
//      return
//    }
//
//    fmt.Println(config["example_string"].Str)
//    fmt.Println(config["example_array"].Arr[0].Str)
//    fmt.Println(config["example_object"].Obj["example_number"].Num)
//    fmt.Println(config["example_object"].Obj["example_number"].Int)
//    //Or
//    fmt.Println(config["example_string"].Str)
//    fmt.Println(config["example_array.0"].Str)
//    fmt.Println(config["example_object.example_number"].Num)
//    fmt.Println(config["example_object.example_number"].Int)
//  }
//
//Outputs
//
//  string value
//  array value 0
//  5.3
//  5
//
//Where ./configs/ExampleConfig.conf is
//
//  {
//    "example_string": "string value",
//    "example_array": [
//      "array value 0"
//    ],
//    "example_object": {
//      "example_number": 5.3
//    }
//  }
package jsonconfig

import (
  "encoding/json"
  "os"
  "strings"
  "strconv"
)

type Configuration map[string]JSONValue

//Is a convenience struct that makes working with abstract JSON data more tolerable.
//The internal values Arr and Obj can be nil, so should not be assumed to be safe.
type JSONValue struct{
  Value interface{}
  Arr []JSONValue
  Str string
  Int int
  Num float64
  Bool bool
  Obj Configuration
}

//Creates a JSONValue from the interface provided. It attempts to fill the values Arr, Str, Int, Num, and Obj
//by checking against the type of the value provided.
func NewJSONValue(value interface{}) JSONValue {
  outputValue := JSONValue{Value: value}
  outputValue.Arr = outputValue.Array()
  outputValue.Str = outputValue.String()
  outputValue.Num = outputValue.Number()
  outputValue.Int = outputValue.Integer()
  outputValue.Bool = outputValue.Boolean()
  outputValue.Obj = outputValue.Object()
  return outputValue
}

//Checks if the type of the json value is an array and if appropriate, casts it into
//an array of JSONValue.
func (key JSONValue) Array() []JSONValue {
  switch typedValue := key.Value.(type) {
    case []interface{}:
      typedArray := make([]JSONValue, len(typedValue))
      for i := range typedValue {
        typedArray[i] = NewJSONValue(typedValue[i])
      }
      return typedArray
    default:
      return nil
  }
}

//Checks if the type of the json value is a string and if appropriate, casts it into a string.
func (key JSONValue) String() string {
  switch typedValue := key.Value.(type) {
    case string:
      return typedValue
    default:
      return ""
  }
}

//Checks if the type of the json value is a float64 and if appropriate, casts it into an int.
func (key JSONValue) Integer() int {
  switch typedValue := key.Value.(type) {
    case float64:
      return int(typedValue)
    default:
      return 0
  }
}

//Checks if the type of the json value is a float64 and if appropriate, casts it into a float64.
func (key JSONValue) Number() float64 {
  switch typedValue := key.Value.(type) {
    case float64:
      return typedValue
    default:
      return 0
  }
}

//Checks if the type of the json value is a bool and if appropriate, casts it into a bool.
func (key JSONValue) Boolean() bool {
  switch typedValue := key.Value.(type) {
    case bool:
      return typedValue
    default:
      return false
  }
}

//Checks if the type of the json value is an object and if appropriate, casts it into a map of JSONValue.
func (key JSONValue) Object() Configuration {
  switch typedValue := key.Value.(type) {
    case map[string]interface{}:
      return convertMap(typedValue)
    default:
      return Configuration{}
  }
}

//Facilitates the collapse of a configuration (map[string]JSONValue) so that you can access sub
//levels with a "." delimeter. This function wont overwrite any keys that already exist, so if
//you have a structure of the form
//  {
//    "example": {
//      "collision": "ignored"
//    },
//    "example.collision": "used"
//  }
//The value "used" will be returned by config["example.collision"].
func (key JSONValue) collapse(path string, config Configuration) {
  if _, exists := config[path]; !exists {
    config[path] = key
  }
  for childKey, childValue := range key.Obj {
    childValue.collapse(childKey, key.Obj)
    childValue.collapse(path + "." + childKey, config)
  }
  for childKey, childValue := range key.Arr {
    childValue.collapse(path + "." + strconv.Itoa(childKey), config)
  }
}

//Flattens a configuration (map[string]JSONValue) so that you can access sub levels with a "." delimeter.
//This function wont overwrite any keys that already exist, so if you have a structure of the form
//  {
//    "example": {
//      "collision": "ignored"
//    },
//    "example.collision": "used"
//  }
//The value "used" will be returned by config["example.collision"].
func (config Configuration) Collapse() {
  for childKey, childValue := range config {
    childValue.collapse(childKey, config)
  }
}

//Takes a "." delimited path and recursively uses the path, returning when a matching structure is found.
//So this func will return used in the following example because example.collision gets matched before
//example: { collision }.
//  {
//    "example": {
//      "collision": "ignored"
//    },
//    "example.collision": "used"
//  }
//This func is only available because a collapsed config pollutes the parent nodes, stopping you from using
//  for key, value := range config {
//    fmt.Println(key)
//  }
//This func also can't handle arrays so for the following config
//  {
//    "array": [
//      {
//        "value": 4
//      }
//    ]
//  }
//config.Get("array.0.value").Num will return nothing while a collapsed config config["array.0.value"].Num
//would return 4.
func (config Configuration) Get(path string) JSONValue {
  if value, ok := config[path]; ok {
    return value
  } else {
    keys := strings.Split(path, ".")
    if subPath, subPathOk := config.Get(strings.Join(keys[0:len(keys) - 1], ".")).Obj[keys[len(keys) - 1]]; subPathOk {
      return subPath
    } else {
      return NewJSONValue(nil)
    }
  }
}

//Converts an abstract map of json data into a map of JSONValue.
func convertMap(from map[string]interface{}) Configuration {
  output := Configuration{}
  for mapKey, mapValue := range from {
    output[mapKey] = NewJSONValue(mapValue)
  }
  return output
}

//Attempts to parse the file as a json object, removing any //comments in the process.
func loadFileAsJSON(filename string) (Configuration, error) {
  file, err := os.Open(filename)
  if err != nil {
    return Configuration{}, err
  }

  untypedMap := map[string]interface{}{}
  dec := json.NewDecoder(NewJsonCommentStripper(file))
  if err = dec.Decode(&untypedMap); err != nil {
    return Configuration{}, err
  }

  return convertMap(untypedMap), nil
}

//Attempts to parse the string as a json object, removing any //comments in the process.
func loadStringAsJSON(jsonstr string) (Configuration, error) {
  untypedMap := map[string]interface{}{}
  dec := json.NewDecoder(NewJsonCommentStripper(strings.NewReader(jsonstr)))
  if err := dec.Decode(&untypedMap); err != nil {
    return Configuration{}, err
  }

  return convertMap(untypedMap), nil
}

//Carefully copies the other Configurations values into the calling config file.
//If the key already exists in the calling config file then the one in the other config
//is ignored unless the value in the other config and the calling config are both objects.
//If the value is an object then the process is repeated, treating this key as a config in both
//the calling config and other config.
func (config Configuration) MergeConfig(other Configuration) {
  for key, value := range other {
    if _, exists := config[key]; !exists {
      config[key] = value
    } else {
      switch value.Value.(type) {
        case map[string]interface{}:
          switch config[key].Value.(type) {
            case map[string]interface{}:
              config[key].Obj.MergeConfig(other[key].Obj)
          }
      }
    }
  }
  // config.Collapse()
}

//Loads the file containing a json object into an abstract map of JSONValue valueType.
//You can provide a default configuration by providing a partial example of the config
//file as a string. This call should be used over LoadAbstract if you wish to use range
//on a JSON object. The collapse performed by LoadAbstract pollutes the keys of parent objects.
func LoadAbstractNoCollapse(filename string, defaults string) (config Configuration, err error) {
  config, err = loadFileAsJSON(filename)
  if err != nil {
    return
  }

  if len(defaults) > 0 {
    defaultValues, err := loadStringAsJSON(defaults)
    if err != nil {
      return Configuration{}, err
    }
    config.MergeConfig(defaultValues)
  }
  return
}

//Loads the file containing a json object into an abstract map of JSONValue valueType.
//You can provide a default configuration by providing a partial example of the config
//file as a string.
func LoadAbstract(filename string, defaults string) (config Configuration, err error) {
  config, err = LoadAbstractNoCollapse(filename, defaults)
  config.Collapse()
  return
}

//Loads the file containing a json object into the provided data structure. You can
//provide default values by defining them in the provided data structure before handing
//it to this func.
func Load(filename string, config interface{}) error {
  file, err := os.Open(filename)
  if err != nil {
    return err
  }

  dec := json.NewDecoder(NewJsonCommentStripper(file))
  if err = dec.Decode(config); err != nil {
    return err
  }

  return nil
}
