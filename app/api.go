package app

import (
	"encoding/json"
	"strings"

	"fmt"
	"os"
	"path"

	"bytes"

	"strconv"

	api "github.com/TIBCOSoftware/flogo-cli/app"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling-lib/model"
	"github.com/TIBCOSoftware/mashling-lib/types"
	"github.com/TIBCOSoftware/mashling-lib/util"
	"github.com/xeipuuv/gojsonschema"
	assets "github.com/TIBCOSoftware/mashling-cli/assets"
)

// CreateMashling creates a gateway application from the specified json gateway descriptor
func CreateMashling(env env.Project, gatewayJson string, appDir string, appName string, vendorDir string) error {

	descriptor, err := model.ParseGatewayDescriptor(gatewayJson)
	if err != nil {
		return err
	}

	if appName != "" {
		altJson := strings.Replace(gatewayJson, `"`+descriptor.Gateway.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := model.ParseGatewayDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Gateway.Name == appName {
			gatewayJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(gatewayJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			gatewayJson = string(updApp)
		}

		descriptor.Gateway.Name = appName
	} else {
		appName = descriptor.Gateway.Name
	}

	flogoAppTriggers := []*ftrigger.Config{}
	flogoAppActions := []*faction.Config{}

	//1. load the configuration, if provided.
	configNamedMap := make(map[string]types.Config)
	for _, config := range descriptor.Gateway.Configurations {
		configNamedMap[config.Name] = config
	}

	triggerNamedMap := make(map[string]types.Trigger)
	for _, trigger := range descriptor.Gateway.Triggers {
		triggerNamedMap[trigger.Name] = trigger
	}

	handlerNamedMap := make(map[string]types.EventHandler)
	for _, evtHandler := range descriptor.Gateway.EventHandlers {
		handlerNamedMap[evtHandler.Name] = evtHandler
	}

	createdHandlers := make(map[string]bool)

	//new map to maintain existing trigger and its settings, to be used in comparing one trigger definition with another
	//createdTriggersSettingsMap := make(map[*ftrigger.Config]map[string]interface{})
	createdTriggersMap := make(map[string]*ftrigger.Config)

	//translate the gateway model to the flogo model
	for _, link := range descriptor.Gateway.EventLinks {
		triggerNames := link.Triggers

		for _, triggerName := range triggerNames {
			dispatches := link.Dispatches

			//create trigger sections for flogo
			/**
			TODO handle condition parsing and setting the condition in the trigger.
			//get the condition if available
			condition := path.If
			//create the trigger using the condition
			.......
			*/
			flogoTrigger, isNew, err := model.CreateFlogoTrigger(configNamedMap, triggerNamedMap[triggerName], handlerNamedMap, dispatches, createdTriggersMap)
			if err != nil {
				return err
			}

			if *isNew {
				//looks like a new trigger has been added
				flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
			} else {
				//looks like an existing trigger with matching settings is found and modified with a new handler
				for index, v := range flogoAppTriggers {
					if v.Name == flogoTrigger.Name {
						// Found the old trigger entry in the list!
						//remove it..
						flogoAppTriggers = append(flogoAppTriggers[:index], flogoAppTriggers[index+1:]...)
						//...and add the modified trigger to the list
						flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
					}
				}
			}

			//create unique handler actions
			for _, dispatch := range dispatches {
				handlerName := dispatch.Handler

				if !createdHandlers[handlerName] {
					//not already created, so create it
					flogoAction, err := model.CreateFlogoFlowAction(handlerNamedMap[handlerName])
					if err != nil {
						return err
					}

					flogoAppActions = append(flogoAppActions, flogoAction)
					createdHandlers[handlerName] = true
				}
			}
		}

	}

	flogoApp := app.Config{
		Name:        descriptor.Gateway.Name,
		Type:        util.Flogo_App_Type,
		Version:     descriptor.Gateway.Version,
		Description: descriptor.Gateway.Description,
		Triggers:    flogoAppTriggers,
		Actions:     flogoAppActions,
	}

	//create flogo PP JSON
	bytes, err := json.MarshalIndent(flogoApp, "", "\t")
	if err != nil {
		return nil
	}

	flogoJson := string(bytes)

	err = api.CreateApp(SetupNewProjectEnv(), flogoJson, appDir, appName, vendorDir)
	if err != nil {
		return err
	}

	fmt.Println("Generated mashling Artifacts.")
	fmt.Println("Building mashling Artifacts.")

	embed := util.Flogo_App_Embed_Config_Property_Default

	envFlogoEmbed := os.Getenv(util.Flogo_App_Embed_Config_Property)
	if len(envFlogoEmbed) > 0 {
		embed, err = strconv.ParseBool(os.Getenv(util.Flogo_App_Embed_Config_Property))
	}

	options := &api.BuildOptions{SkipPrepare: false, PrepareOptions: &api.PrepareOptions{OptimizeImports: false, EmbedConfig: embed}}
	api.BuildApp(SetupExistingProjectEnv(appDir), options)
	//delete flogo.json file from the app dir
	fgutil.DeleteFilesWithPrefix(appDir, "flogo")
	//create the mashling json descriptor file
	err = fgutil.CreateFileFromString(path.Join(appDir, util.Gateway_Definition_File_Name), gatewayJson)
	if err != nil {
		return err
	}

	fmt.Println("Mashling gateway successfully built!")

	return nil
}

// TranslateGatewayJSON2FlogoJSON tanslates mashling json to flogo json
func TranslateGatewayJSON2FlogoJSON(gatewayJSON string) (string, error) {
	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		return "", err
	}

	flogoAppTriggers := []*ftrigger.Config{}
	flogoAppActions := []*faction.Config{}

	//1. load the configuration, if provided.
	configNamedMap := make(map[string]types.Config)
	for _, config := range descriptor.Gateway.Configurations {
		configNamedMap[config.Name] = config
	}

	triggerNamedMap := make(map[string]types.Trigger)
	for _, trigger := range descriptor.Gateway.Triggers {
		triggerNamedMap[trigger.Name] = trigger
	}

	handlerNamedMap := make(map[string]types.EventHandler)
	for _, evtHandler := range descriptor.Gateway.EventHandlers {
		handlerNamedMap[evtHandler.Name] = evtHandler
	}

	createdHandlers := make(map[string]bool)

	//new map to maintain existing trigger and its settings, to be used in comparing one trigger definition with another
	createdTriggersMap := make(map[string]*ftrigger.Config)

	//translate the gateway model to the flogo model
	for _, link := range descriptor.Gateway.EventLinks {
		triggerNames := link.Triggers

		for _, triggerName := range triggerNames {
			dispatches := link.Dispatches

			//create trigger sections for flogo
			/**
			TODO handle condition parsing and setting the condition in the trigger.
			//get the condition if available
			condition := path.If
			//create the trigger using the condition
			.......
			*/
			flogoTrigger, isNew, err := model.CreateFlogoTrigger(configNamedMap, triggerNamedMap[triggerName], handlerNamedMap, dispatches, createdTriggersMap)
			if err != nil {
				return "", err
			}

			//	check if the trigger is a new trigger or a "logically" same trigger.
			if *isNew {
				//looks like a new trigger has been added
				flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
			} else {
				//looks like an existing trigger with matching settings is found and is now modified with a new handler
				for index, v := range flogoAppTriggers {
					if v.Name == flogoTrigger.Name {
						// Found the old trigger entry in the list!
						//remove it..
						flogoAppTriggers = append(flogoAppTriggers[:index], flogoAppTriggers[index+1:]...)
						//...and add the modified trigger to the list
						flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
					}
				}
			}

			//create unique handler actions
			for _, dispatch := range dispatches {
				handlerName := dispatch.Handler

				if !createdHandlers[handlerName] {
					//not already created, so create it
					flogoAction, err := model.CreateFlogoFlowAction(handlerNamedMap[handlerName])
					if err != nil {
						return "", err
					}

					flogoAppActions = append(flogoAppActions, flogoAction)
					createdHandlers[handlerName] = true
				}
			}
		}

	}

	flogoApp := app.Config{
		Name:        descriptor.Gateway.Name,
		Type:        util.Flogo_App_Type,
		Version:     descriptor.Gateway.Version,
		Description: descriptor.Gateway.Description,
		Triggers:    flogoAppTriggers,
		Actions:     flogoAppActions,
	}

	//create flogo PP JSON
	bytes, err := json.MarshalIndent(flogoApp, "", "\t")
	if err != nil {
		return "", nil
	}

	flogoJSON := string(bytes)

	return flogoJSON, nil
}

// BuildMashling Builds mashling gateway
func BuildMashling(appDir string, gatewayJSON string) error {

	//create flogo.json from gateway descriptor
	flogoJSON, err := TranslateGatewayJSON2FlogoJSON(gatewayJSON)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Error while processing gateway descriptor.\n\n")
		return err
	}
	err = fgutil.CreateFileFromString(path.Join(appDir, "flogo.json"), flogoJSON)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Error while creating flogo.json.\n\n")
		return err
	}

	//Install dependencies explicitly, as api.BuildApp() doesn't install newly added dependencies.
	//Workaround for https://github.com/TIBCOSoftware/flogo-cli/issues/56
	fmt.Println("Installing dependencies...")
	env := SetupExistingProjectEnv(appDir)
	flogoAppDescriptor, err := api.ParseAppDescriptor(flogoJSON)
	deps := api.ExtractDependencies(flogoAppDescriptor)

	for _, dep := range deps {
		path, version := splitVersion(dep.Ref)
		err = env.InstallDependency(path, version)
		if err != nil {
			return err
		}
	}
	//END of workaround https://github.com/TIBCOSoftware/flogo-cli/issues/56

	options := &api.BuildOptions{SkipPrepare: false, PrepareOptions: &api.PrepareOptions{OptimizeImports: false, EmbedConfig: false}}
	api.BuildApp(SetupExistingProjectEnv(appDir), options)

	//delete flogo.json file from the app dir
	fgutil.DeleteFilesWithPrefix(appDir, "flogo")

	fmt.Println("Mashling gateway successfully built!")

	return nil

}

func ListComponents(env env.Project, cType ComponentType) ([]*Component, error) {

	var components []*Component

	rootDir := env.GetRootDir()
	mashlingDescriptorFile := rootDir + "/" + util.Gateway_Definition_File_Name
	mashlingJson, err1 := fgutil.LoadLocalFile(mashlingDescriptorFile)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", mashlingDescriptorFile, err1.Error())
		os.Exit(2)
	}

	microgateway, err := model.ParseGatewayDescriptor(mashlingJson)

	if cType == 2 || cType == TRIGGER {
		if microgateway.Gateway.Triggers != nil {
			for _, trigger := range microgateway.Gateway.Triggers {
				components = append(components, &Component{Name: trigger.Name, Type: TRIGGER, Ref: trigger.Type})
			}
		}
	}

	if cType == 3 || cType == HANDLER {
		if microgateway.Gateway.EventHandlers != nil {
			for _, handler := range microgateway.Gateway.EventHandlers {
				cType.String()
				components = append(components, &Component{Name: handler.Name, Type: HANDLER, Ref: handler.Reference})
			}
		}
	}

	return components, err
}

func ListLinks(env env.Project, cType ComponentType) ([]*types.EventLink, error) {

	rootDir := env.GetRootDir()
	var links []*types.EventLink

	mashlingDescriptorFile := rootDir + "/" + util.Gateway_Definition_File_Name
	mashlingJson, err1 := fgutil.LoadLocalFile(mashlingDescriptorFile)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", mashlingDescriptorFile, err1.Error())
		os.Exit(2)
	}

	microgateway, err := model.ParseGatewayDescriptor(mashlingJson)

	if cType == 1 || cType == LINK {
		if microgateway.Gateway.EventLinks != nil {
			for _, link := range microgateway.Gateway.EventLinks {
				links = append(links, &link)
			}
		}
	}

	return links, err
}

// GetGatewayDetails returns gateway details i.e all Triggers, Handlers & Links
func GetGatewayDetails(env env.Project, cType ComponentType) (string, error) {
	gwInfoBuffer := bytes.NewBufferString("")
	rootDir := env.GetRootDir()
	mashlingDescriptorFile := rootDir + "/" + util.Gateway_Definition_File_Name
	mashlingJSON, err := fgutil.LoadLocalFile(mashlingDescriptorFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", mashlingDescriptorFile, err.Error())
		os.Exit(2)
	}
	microgateway, err := model.ParseGatewayDescriptor(mashlingJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Error while parsing gateway description json - %s\n\n", err.Error())
		os.Exit(2)
	}

	tmpString := ""
	//Triggers info
	if cType == TRIGGER || cType == ALL {
		tmpString = fmt.Sprintf("Triggers: %d\n", len(microgateway.Gateway.Triggers))
		gwInfoBuffer.WriteString(tmpString)
		for _, trigger := range microgateway.Gateway.Triggers {
			tmpString = "\t" + trigger.Name + "  " + trigger.Type + "\n"
			gwInfoBuffer.WriteString(tmpString)
		}
	}

	//Handlers info
	if cType == HANDLER || cType == ALL {
		tmpString = fmt.Sprintf("Handlers: %d\n", len(microgateway.Gateway.EventHandlers))
		gwInfoBuffer.WriteString(tmpString)
		for _, handler := range microgateway.Gateway.EventHandlers {
			tmpString = "\t" + handler.Name + "  " + handler.Reference + "\n"
			gwInfoBuffer.WriteString(tmpString)
		}
	}

	//Links info
	if cType == LINK || cType == ALL {
		links := microgateway.Gateway.EventLinks
		tmpString = fmt.Sprintf("Links: %d\n", len(links))
		gwInfoBuffer.WriteString(tmpString)
		//loop through links
		for _, link := range links {
			gwInfoBuffer.WriteString("\tTrigger: ")
			for _, trigger := range link.Triggers {
				gwInfoBuffer.WriteString(trigger + " ")
			}
			gwInfoBuffer.WriteString("\n")
			gwInfoBuffer.WriteString("\tHandlers:\n")
			for _, dispatcher := range link.Dispatches {
				gwInfoBuffer.WriteString("\t\t" + dispatcher.Handler + "\n")
			}
			gwInfoBuffer.WriteString("\n")
		}
	}

	unLinkedTriggers := 0
	tmpBuf := bytes.NewBufferString("")
	if cType == ALL {
		//Unlinked triggers
		for _, trigger := range microgateway.Gateway.Triggers {
			triggerFound := false
			for _, link := range microgateway.Gateway.EventLinks {
				for _, trigger2 := range link.Triggers {
					if trigger.Name == trigger2 {
						triggerFound = true
						break
					}
				}
				if triggerFound {
					break
				}
			}
			if !triggerFound {
				unLinkedTriggers++
				tmpBuf.WriteString("\t" + trigger.Name + "  " + trigger.Type + "\n")
			}
		}
		if unLinkedTriggers != 0 {
			gwInfoBuffer.WriteString(fmt.Sprintf("Unlinked Triggers: %d", unLinkedTriggers) + "\n")
			gwInfoBuffer.WriteString(tmpBuf.String())
		}

		//Unliked handlers
		unlinkedHandlers := 0
		tmpBuf.Reset()
		for _, handler := range microgateway.Gateway.EventHandlers {
			handlerFound := false
			for _, link := range microgateway.Gateway.EventLinks {
				for _, dispatch := range link.Dispatches {
					if handler.Name == dispatch.Handler {
						handlerFound = true
						break
					}
				}
				if handlerFound {
					break
				}
			}
			if !handlerFound {
				unlinkedHandlers++
				tmpBuf.WriteString("\t" + handler.Name + "  " + handler.Reference + "\n")
			}
		}
		if unlinkedHandlers != 0 {
			gwInfoBuffer.WriteString(fmt.Sprintf("Unlinked Handlers: %d", unlinkedHandlers) + "\n")
			gwInfoBuffer.WriteString(tmpBuf.String())
		}
	}

	return gwInfoBuffer.String(), nil
}

//ValidateGateway validates the gatewway schema instance
func ValidateGateway(gatewayJson string) error {

	schema, err := assets.Asset("schema/gateway_schema.json")
	if err != nil {
		panic(err.Error())
	}
	schemaString := string(schema)
	schemaLoader := gojsonschema.NewStringLoader(schemaString)
	documentLoader := gojsonschema.NewStringLoader(gatewayJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The gateway json is valid\n")
	} else {
		fmt.Printf("The gateway json not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	return err

}
