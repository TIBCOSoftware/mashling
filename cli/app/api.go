/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"fmt"
	"os"
	"path"

	"bytes"

	"strconv"

	api "github.com/TIBCOSoftware/flogo-cli/app"
	config "github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	assets "github.com/TIBCOSoftware/mashling/cli/assets"
	"github.com/TIBCOSoftware/mashling/cli/env"
	"github.com/TIBCOSoftware/mashling/lib/model"
	"github.com/TIBCOSoftware/mashling/lib/types"
	"github.com/TIBCOSoftware/mashling/lib/util"
	"github.com/xeipuuv/gojsonschema"
)

// CreateMashling creates a gateway application from the specified json gateway descriptor
func CreateMashling(env env.Project, gatewayJson string, manifest io.Reader, appDir string, appName string, vendorDir string, customizeFunc func() error) error {

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

	err = CreateApp(SetupNewProjectEnv(), flogoJson, manifest, appDir, appName, vendorDir)
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

	if customizeFunc != nil {
		err = customizeFunc()
		if err != nil {
			return err
		}
	}
	options := &api.BuildOptions{SkipPrepare: false, PrepareOptions: &api.PrepareOptions{OptimizeImports: false, EmbedConfig: embed}}
	BuildApp(SetupExistingProjectEnv(appDir), options)
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
	/*
		//This creates an issue of pulling down the latest flogo packages on every build.
		//User can run 'create' instead of 'build' if there is an additional packages required
		//for the recipe. Thus, commenting out to make sure the build command executes only with
		//the local vendor folder.
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
	*/
	//END of workaround https://github.com/TIBCOSoftware/flogo-cli/issues/56

	embed := util.Flogo_App_Embed_Config_Property_Default

	envFlogoEmbed := os.Getenv(util.Flogo_App_Embed_Config_Property)
	if len(envFlogoEmbed) > 0 {
		embed, err = strconv.ParseBool(os.Getenv(util.Flogo_App_Embed_Config_Property))
	}

	options := &api.BuildOptions{SkipPrepare: false, PrepareOptions: &api.PrepareOptions{OptimizeImports: false, EmbedConfig: embed}}
	BuildApp(SetupExistingProjectEnv(appDir), options)

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

// PublishToMashery publishes to mashery
func PublishToMashery(user *ApiUser, appDir string, gatewayJSON string, host string, mock bool) error {
	// Get HTTP triggers from JSON
	swaggerDoc, err := generateSwagger(host, "", gatewayJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate swagger doc\n\n")
		return err
	}

	token, err := user.FetchOAuthToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to fetch the OAauth token\n\n")
		return err
	}

	// Delay to avoid hitting QPS limit
	delayMilli(500)

	tfSwaggerDoc, err := user.TransformSwagger(string(swaggerDoc), token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to transform swagger doc\n\n")
		return err
	}

	// Only need the value of 'document'. Including the rest will cause errors
	m := map[string]interface{}{}
	if err = json.Unmarshal([]byte(tfSwaggerDoc), &m); err != nil {
		panic(err)
	}

	var cleanedTfSwaggerDoc []byte

	if cleanedTfSwaggerDoc, err = json.Marshal(m["document"]); err != nil {
		panic(err)
	}

	if mock == false {
		s, err := user.CreateAPI(string(cleanedTfSwaggerDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create the api %s\n\n", s)
			fmt.Errorf("%v", err)
			return err
		}

		fmt.Println("Successfully published to mashery!")
	} else {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, cleanedTfSwaggerDoc, "", "\t")
		if err != nil {
			return err
		}

		fmt.Printf("%s", prettyJSON.Bytes())
		fmt.Println("\nMocked! Did not attempt to publish.\n")
	}

	return nil
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

//IsValidGateway validates the gateway schema instance returns bool and error
func IsValidGateway(gatewayJSON string) (bool, error) {

	isValidSchema := false

	suplliedSchemaVersion, err := getSchemaVersion(gatewayJSON)

	if err != nil {
		return isValidSchema, err
	}

	schemaPath, isValidSchema := GetSupportedSchema(suplliedSchemaVersion)

	//check whether CLI supports this schema version
	if !isValidSchema {
		fmt.Printf("Schema version [%v] not supported. Please upgrade mashling cli \n", suplliedSchemaVersion)
		return isValidSchema, nil
	}

	schema, err := assets.Asset(schemaPath)
	if err != nil {
		panic(err.Error())
	}
	schemaString := string(schema)
	schemaLoader := gojsonschema.NewStringLoader(schemaString)
	documentLoader := gojsonschema.NewStringLoader(gatewayJSON)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}

	if result.Valid() {
		isValidSchema = true
	} else {
		fmt.Printf("The gateway json is not valid. See errors:\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		isValidSchema = false
	}

	return isValidSchema, err

}

func getSchemaVersion(gatewayJSON string) (string, error) {

	suplliedSchema := ""
	gatewayDescriptor := &struct {
		MashlingSchema string `json:"mashling_schema"`
	}{}
	err := json.Unmarshal([]byte(gatewayJSON), gatewayDescriptor)

	if err != nil {
		return suplliedSchema, err
	}
	suplliedSchema = gatewayDescriptor.MashlingSchema

	return suplliedSchema, err
}

// CreateApp creates an application from the specified json application descriptor
func CreateApp(env env.Project, appJson string, manifest io.Reader, appDir string, appName string, vendorDir string) error {

	descriptor, err := api.ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	if appName != "" {
		// override the application name

		altJson := strings.Replace(appJson, `"`+descriptor.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := api.ParseAppDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == appName {
			appJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(appJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			appJson = string(updApp)
		}

		descriptor.Name = appName
	}

	env.Init(appDir)
	err = env.Create(false, vendorDir)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(path.Join(appDir, "flogo.json"), appJson)
	if err != nil {
		return err
	}

	deps := config.ExtractDependencies(descriptor)

	//if manifest exists, use it to set up the dependecies
	err = env.RestoreDependency(manifest)
	if err == nil {
		fmt.Println("Dependent libraries are restored.")
	} else {
		//todo allow ability to specify flogo-lib version
		env.InstallDependency("github.com/TIBCOSoftware/flogo-lib", "")

		for _, dep := range deps {
			path, version := splitVersion(dep.Ref)
			err = env.InstallDependency(path, version)
			/*
				if err != nil {
					return err
				}
			*/
		}
	}

	// create source files
	cmdPath := path.Join(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	os.MkdirAll(cmdPath, 0777)

	CreateMainGoFile(cmdPath, "")
	CreateImportsGoFile(cmdPath, deps)

	return nil
}

func BuildApp(env env.Project, options *api.BuildOptions) (err error) {

	if options == nil {
		options = &api.BuildOptions{}
	}

	if !options.SkipPrepare {
		err = PrepareApp(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	if !options.EmbedConfig {
		fgutil.CopyFile(path.Join(env.GetRootDir(), fileDescriptor), path.Join(env.GetBinDir(), fileDescriptor))
		if err != nil {
			return err
		}
	} else {
		os.Remove(path.Join(env.GetBinDir(), fileDescriptor))
	}

	return
}

// PrepareApp do all pre-build setup and pre-processing
func PrepareApp(env env.Project, options *api.PrepareOptions) (err error) {

	if options == nil {
		options = &api.PrepareOptions{}
	}

	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(env)
		if err != nil {
			return err
		}
	}

	//generate metadata
	err = generateGoMetadata(env)
	if err != nil {
		return err
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(path.Join(env.GetRootDir(), "flogo.json"))

	if err != nil {
		return err
	}
	descriptor, err := api.ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	//generate imports file
	var deps []*config.Dependency

	if options.OptimizeImports {

		deps = config.ExtractDependencies(descriptor)

	} else {
		deps, err = ListDependencies(env, 0)
	}

	cmdPath := path.Join(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	CreateImportsGoFile(cmdPath, deps)

	removeEmbeddedAppGoFile(cmdPath)
	removeShimGoFiles(cmdPath)

	if options.Shim != "" {

		removeMainGoFile(cmdPath) //todo maybe rename if it exists
		createShimSupportGoFile(cmdPath, appJson, options.EmbedConfig)

		fmt.Println("Shim:", options.Shim)

		for _, value := range descriptor.Triggers {

			fmt.Println("Id:", value.ID)
			if value.ID == options.Shim {
				triggerPath := path.Join(env.GetVendorSrcDir(), value.Ref, "trigger.json")

				mdJson, err := fgutil.LoadLocalFile(triggerPath)
				if err != nil {
					return err
				}
				metadata, err := api.ParseTriggerMetadata(mdJson)
				if err != nil {
					return err
				}

				if metadata.Shim != "" {

					//todo blow up if shim file not found
					shimFilePath := path.Join(env.GetVendorSrcDir(), value.Ref, dirShim, fileShimGo)
					fmt.Println("Shim File:", shimFilePath)
					fgutil.CopyFile(shimFilePath, path.Join(cmdPath, fileShimGo))

					if metadata.Shim == "plugin" {
						//look for Makefile and execute it
						makeFilePath := path.Join(env.GetVendorSrcDir(), value.Ref, dirShim, makeFile)
						fmt.Println("Make File:", makeFilePath)
						fgutil.CopyFile(makeFilePath, path.Join(cmdPath, makeFile))

						// Copy the vendor folder (Ugly workaround, this will go once our app is golang structure compliant)
						vendorDestDir := path.Join(cmdPath, "vendor")
						_, err = os.Stat(vendorDestDir)
						if err == nil {
							// We don't support existing vendor folders yet
							return fmt.Errorf("Unsupported vendor folder found for function build, please create an issue on https://github.com/TIBCOSoftware/flogo")
						}
						// Create vendor folder
						err = api.CopyDir(env.GetVendorSrcDir(), vendorDestDir)
						if err != nil {
							return err
						}
						defer os.RemoveAll(vendorDestDir)

						// Execute make
						cmd := exec.Command("make", "-C", cmdPath)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Env = append(os.Environ(),
							fmt.Sprintf("GOPATH=%s", env.GetRootDir()),
						)

						err = cmd.Run()
						if err != nil {
							return err
						}
					}
				}

				break
			}
		}

	} else if options.EmbedConfig {
		createEmbeddedAppGoFile(cmdPath, appJson)
	}

	return
}

func generateGoMetadata(env env.Project) error {
	//todo optimize metadata recreation to minimize compile times
	dependencies, err := ListDependencies(env, 0)

	if err != nil {
		return err
	}

	for _, dependency := range dependencies {
		createMetadata(env, dependency)
	}

	return nil
}

func createMetadata(env env.Project, dependency *config.Dependency) error {

	vendorSrc := env.GetVendorSrcDir()
	mdFilePath := path.Join(vendorSrc, dependency.Ref)
	mdGoFilePath := path.Join(vendorSrc, dependency.Ref)
	pkg := path.Base(mdFilePath)

	tplMetadata := tplMetadataGoFile

	switch dependency.ContribType {
	case config.ACTION:
		mdFilePath = path.Join(mdFilePath, "action.json")
		mdGoFilePath = path.Join(mdGoFilePath, "action_metadata.go")
	case config.TRIGGER:
		mdFilePath = path.Join(mdFilePath, "trigger.json")
		mdGoFilePath = path.Join(mdGoFilePath, "trigger_metadata.go")
		tplMetadata = tplTriggerMetadataGoFile
	case config.ACTIVITY:
		mdFilePath = path.Join(mdFilePath, "activity.json")
		mdGoFilePath = path.Join(mdGoFilePath, "activity_metadata.go")
		tplMetadata = tplActivityMetadataGoFile
	default:
		return nil
	}

	raw, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		return err
	}

	info := &struct {
		Package      string
		MetadataJSON string
	}{
		Package:      pkg,
		MetadataJSON: string(raw),
	}

	f, _ := os.Create(mdGoFilePath)
	fgutil.RenderTemplate(f, tplMetadata, info)
	f.Close()

	return nil
}

var tplMetadataGoFile = `package {{.Package}}

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

func getJsonMetadata() string {
	return jsonMetadata
}
`

var tplActivityMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(NewActivity(md))
}
`

var tplTriggerMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
`

// ListDependencies lists all installed dependencies
func ListDependencies(env env.Project, cType config.ContribType) ([]*config.Dependency, error) {

	vendorSrc := env.GetVendorSrcDir()
	var deps []*config.Dependency

	err := filepath.Walk(vendorSrc, func(filePath string, info os.FileInfo, _ error) error {

		if !info.IsDir() {

			switch info.Name() {
			case "action.json":
				if cType == 0 || cType == config.ACTION {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:action" {
						deps = append(deps, &config.Dependency{ContribType: config.ACTION, Ref: ref})
					}
				}
			case "trigger.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0 : len(filePath)-12]
				if _, err := os.Stat(fmt.Sprintf("%s/../trigger.json", dir)); err == nil {
					//old trigger.json, ignore
					return nil
				}
				if cType == 0 || cType == config.TRIGGER {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:trigger" {
						deps = append(deps, &config.Dependency{ContribType: config.TRIGGER, Ref: ref})
					}
				}
			case "activity.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0 : len(filePath)-13]
				if _, err := os.Stat(fmt.Sprintf("%s/../activity.json", dir)); err == nil {
					//old activity.json, ignore
					return nil
				}
				if cType == 0 || cType == config.ACTIVITY {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:activity" {
						deps = append(deps, &config.Dependency{ContribType: config.ACTIVITY, Ref: ref})
					}
				}
			case "flow-model.json":
				if cType == 0 || cType == config.FLOW_MODEL {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:flow-model" {
						deps = append(deps, &config.Dependency{ContribType: config.FLOW_MODEL, Ref: ref})
					}
				}
			}

		}

		return nil
	})

	return deps, err
}

func refPath(vendorSrc string, filePath string) string {

	startIdx := len(vendorSrc) + 1
	endIdx := strings.LastIndex(filePath, string(os.PathSeparator))

	return strings.Replace(filePath[startIdx:endIdx], string(os.PathSeparator), "/", -1)
}

func readDescriptor(path string, info os.FileInfo) (*config.Descriptor, error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	return api.ParseDescriptor(string(raw))
}
