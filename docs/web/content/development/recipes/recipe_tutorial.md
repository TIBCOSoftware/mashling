---
title: Build recipe example
weight: 4210
pre: "<i class=\"fa fa-asterisk\" aria-hidden=\"true\"></i> "
---

In this example we will create a mashling gateway recipe that conditionally invokes a PetStore API.
The recipe will either register or retrieve pet data. The scenario assumes that only pets with a category name of `dog` or `cat` are supported. If an attempt is made to register any other kind of pet the gateway responds with an error message without hitting the PetStore backend.


Open an editor and define event triggers like the following:

```json
{
  "mashling_schema": "1.0",
  "gateway": {
    "name": "MyProxy",
    "version": "1.0.0",
    "description": "This is a simple proxy.",
    "triggers": [
      {
        "name": "MyProxy",
        "description": "Animals rest trigger - PUT animal details",
        "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
        "settings": {
          "port": "9096"
        },
        "handlers": [
          {
            "dispatch": "Retrieve",
            "settings": {
              "autoIdReply": "false",
              "method": "GET",
              "path": "/pets/{petId}",
              "useReplyHandler": "false"
            }
          },
          {
            "dispatch": "Register",
            "settings": {
              "autoIdReply": "false",
              "method": "PUT",
              "path": "/pets",
              "useReplyHandler": "false"
            }
          }
        ]
      }
    ],
```

There are two handlers created for this HTTP trigger. One for retrieving the pet data through an HTTP GET operation and the other for registering a pet through an HTTP PUT operation. For the details of the handler configuration, see the documentation at https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/gorillamuxtrigger

Add a dispatches section:

```json
    "dispatches": [
      {
        "name": "Retrieve",
        "routes": [
          {
            "steps": [
              {
                "service": "GetPet",
                "input": {
                  "inputs.pathParams": "${payload.pathParams}"
                }
              }
            ],
            "responses": [
              {
                "if": "JSON.parse(GetPet.response.outputs.code) != 200",
                "error": true,
                "output": {
                  "error": "Pet is not available."
                }
              },
              {
                "if": "JSON.parse(GetPet.response.outputs.code) == 200",
                "error": false,
                "output": {
                  "result": "${GetPet.response.outputs.data}"
                }
              }
            ]
          }
        ]
      },
      {
        "name": "Register",
        "routes": [
          {
            "if": "payload.content.category.name == 'DOG' || payload.content.category.name == 'CAT'",
            "steps": [
              {
                "service": "PutPet",
                "input": {
                  "inputs.content": "${payload.content}"
                }
              }
            ],
            "responses": [
              {
                "if": "PutPet.response.outputs.code != 200",
                "error": true,
                "output": {
                  "error": "Pet is not added."
                }
              },
              {
                "if": "PutPet.response.outputs.code == 200",
                "error": false,
                "output": {
                  "success": "${PutPet.response.outputs.data}"
                }
              }
            ]
          },
          {
            "steps": [
              {
                "service": "InvalidAnimal",
                "input": {
                  "parameters.content": "${payload.content}"
                }
              }
            ],
            "responses": [
              {
                "error": false,
                "output": {
                  "error": "${InvalidAnimal.response.result.msg}"
                }
              }
            ]
          }
        ]
      }
    ],
```

Both `Retrieve` and `Register` routes have their own steps and their own response handling defined. Note that `Register` has a step which is executed only when a dog or cat is registered.
If no condition is defined for a set of steps it is considered the default route. Only one route is executed. Thus, in the `Register` route, the InvalidAnimal service is executed only when the request is for an animal other than a dog or a cat.

Add the services section below, which defines each service reference used in the dispatch section above:

```json
    "services": [
      {
        "name": "GetPet",
        "description": "Make GET calls against a remote HTTP service using a Flogo flow.",
        "type": "flogoFlow",
        "settings": {
          "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestGetActivity.json"
        }
      },
      {
        "name": "PutPet",
        "description": "Make PUT calls against a remote HTTP service using a Flogo flow.",
        "type": "flogoFlow",
        "settings": {
          "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
        }
      },
      {
        "name": "InvalidAnimal",
        "description": "Message for an invalid animal",
        "type": "js",
        "settings": {
          "script": "result.msg = 'Unsupported animal: ' + parameters.content.category.name;"
        }
      }
    ]
  }
}
```

The first two services invoke a flogo flow to execute the pet store backend API. InvalidAnimal service executes a javascript service to produce a message that the request is invalid. This avoids invoking the backend API if the request is invalid.

Save the recipe file created and validate it with mashling-cli in a terminal:
```
mashling-cli -validate recipe_tutorial.json
```

Start the recipe with mashling-gateway:
```
mashling-gateway -c recipe_tutorial.json
```

In another terminal, register a pet with

```
curl -X PUT "http://localhost:9096/pets" -H "Content-Type: application/json" -d '{"category":{"id":1,"name":"DOG"},"id":16,"name":"Olive","photoUrls":["unavailable"],"status":"sold","tags":[{"id":76543,"name":"Olive"}]}'
```

Now, retrieve the pet data with

```
curl -request GET http://localhost:9096/pets/16
```

The pet data just registered should be returned.

Try registering a pet with unsupported type

```
curl -X PUT "http://localhost:9096/pets" -H "Content-Type: application/json" -d '{"category":{"id":3,"name":"BIRD"},"id":17,"name":"Birdie","photoUrls":["unavailable"],"status":"sold","tags":[{"id":87654,"name":"Birdie"}]}'
```

An error message is returned like the following:
```
{
 "error": "Unsupported animal: BIRD"
}
```
