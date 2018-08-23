package main

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger"
	_ "github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry"
	_ "github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry/triggers"
	mAction "github.com/TIBCOSoftware/mashling/pkg/flogo/action"
)

func main() {
	app := flogo.NewApp()
	app.AddResource("mashling:HTTP", Data)

	trg := app.NewTrigger(&gorillamuxtrigger.RestTrigger{}, map[string]interface{}{"port": "8080"})
	h1 := trg.NewHandler(map[string]interface{}{"method": "GET", "path": "/magic", "autoIdReply": "false", "useReplyHandler": "true"})
	h1.NewAction(&mAction.MashlingAction{}, map[string]interface{}{"mashlingURI": "mashling:HTTP"})
	e, err := flogo.NewEngine(app)
	if err != nil {
		logger.Error(err)
		return
	}

	engine.RunEngine(e)
}

var Data = json.RawMessage(`{
    "pattern": "DefaultHttpPattern",
    "configuration": {
        "backendUrl": "https://petstore.swagger.io/v2/pet/32",
        "jwtKey": "qwertyuiopasdfghjklzxcvbnm123456",
        "useCircuitBreaker": true,
        "useJWT": true
    }
}`)
