package mashery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/TIBCOSoftware/mashling/internal/pkg/logger"
)

type ApiUser struct {
	Username     string
	Password     string
	ApiKey       string
	ApiSecretKey string
	Uuid         string
	Portal       string
	Noop         bool
}

const (
	masheryUri   = "https://api.mashery.com"
	restUri      = "/v3/rest/"
	transformUri = "transform"
	accessToken  = "access_token"
)

func shortDelay() {
	time.Sleep(time.Duration(500) * time.Millisecond)
}

// PublishToMashery publishes to mashery
func PublishToMashery(user *ApiUser, appDir string, swaggerDoc []byte, host string, mock bool, iodocs bool, testplan bool, apiTemplateJSON []byte) error {
	// Get HTTP triggers from JSON
	token, err := user.FetchOAuthToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to fetch the OAauth token\n\n")
		return err
	}

	// Delay to avoid hitting QPS limit
	shortDelay()

	mApi, err := TransformSwagger(user, string(swaggerDoc), "swagger2", "masheryapi", token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to transform swagger to mashery api\n\n")
		return err
	}

	shortDelay()

	mIodoc, err := TransformSwagger(user, string(swaggerDoc), "swagger2", "iodocsv1", token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to transform swagger to mashery iodocs\n\n")
		return err
	}

	shortDelay()

	templApi, templEndpoint, templPackage, templPlan := BuildMasheryTemplates(string(apiTemplateJSON))
	if mock == false {

		mApi = UpdateApiWithDefaults(mApi, templApi, templEndpoint)

		apiId, apiName, endpoints, updated := CreateOrUpdateApi(user, token, MapToByteArray(mApi), mApi)

		if iodocs == true {

			cleanedTfIodocSwaggerDoc := UpdateIodocsDataWithApi(MapToByteArray(mIodoc), apiId)

			CreateOrUpdateIodocs(user, token, cleanedTfIodocSwaggerDoc, apiId, updated)
			shortDelay()
		}

		var key string
		if testplan == true {

			packagePlanDoc := CreatePackagePlanDataFromApi(apiId, apiName, endpoints)
			packagePlanDoc = UpdatePackageWithDefaults(packagePlanDoc, templPackage, templPlan)
			var marshalledDoc []byte
			marshalledDoc, err = json.Marshal(packagePlanDoc)
			if err != nil {
				panic(err)
			}

			shortDelay()

			p := CreateOrUpdatePackage(user, token, marshalledDoc, apiName, updated)

			shortDelay()

			key = CreateApplicationAndKey(user, token, p, apiName)

		}
		fmt.Println("==================================================================")
		fmt.Printf("Successfully published to mashery= API %s (id=%s)\n", apiName, apiId)
		fmt.Println("==================================================================")
		fmt.Println("API Control Center Link: https://" + strings.Replace(user.Portal, "api", "admin", -1) + "/control-center/api-definitions/" + apiId)
		if testplan == true {
			fmt.Println("==================================================================")
			fmt.Println("Example Curls:")
			for _, endpoint := range endpoints {
				ep := endpoint.(map[string]interface{})
				fmt.Println(GenerateExampleCall(ep, key))
			}
		}
	} else {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, MapToByteArray(mApi), "", "\t")
		if err != nil {
			return err
		}

		//fmt.Printf("%s", prettyJSON.Bytes())
		fmt.Println("Mocked! Did not attempt to publish.")
	}

	return nil
}

func UpdateApiWithDefaults(mApi map[string]interface{}, templApi map[string]interface{}, templEndpoint map[string]interface{}) map[string]interface{} {
	var m1 map[string]interface{}
	json.Unmarshal(MapToByteArray(mApi), &m1)
	merged := merge(m1, templApi, 0)
	m_d := m1["endpoints"].([]interface{})

	items := []map[string]interface{}{}

	for _, d_item := range m_d {
		merged := merge(d_item.(map[string]interface{}), templEndpoint, 0)
		items = append(items, merged)
	}

	merged["endpoints"] = items
	return merged

}

func UpdatePackageWithDefaults(mApi map[string]interface{}, templPackage map[string]interface{}, templPlan map[string]interface{}) map[string]interface{} {
	var m1 map[string]interface{}
	json.Unmarshal(MapToByteArray(mApi), &m1)
	merged := merge(m1, templPackage, 0)
	m_d := m1["plans"].([]interface{})

	items := []map[string]interface{}{}

	for _, d_item := range m_d {
		merged := merge(d_item.(map[string]interface{}), templPlan, 0)
		items = append(items, merged)
	}

	merged["plans"] = items
	return merged

}

func BuildMasheryTemplates(apiTemplateJSON string) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}) {
	apiTemplate := map[string]interface{}{}
	endpointTemplate := map[string]interface{}{}
	packageTemplate := map[string]interface{}{}
	planTemplate := map[string]interface{}{}

	if apiTemplateJSON != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(apiTemplateJSON), &m); err != nil {
			panic(err)
		}
		apiTemplate = m["api"].(map[string]interface{})
		endpointTemplate = apiTemplate["endpoint"].(map[string]interface{})
		delete(apiTemplate, "endpoint")
		packageTemplate = m["package"].(map[string]interface{})
		planTemplate = packageTemplate["plan"].(map[string]interface{})
		delete(packageTemplate, "plan")

	} else {
		apiTemplate["qpsLimitOverall"] = 0
		endpointTemplate["requestAuthenticationType"] = "apiKeyAndSecret_SHA256"
		packageTemplate["sharedSecretLength"] = 10
		planTemplate["selfServiceKeyProvisioningEnabled"] = false

	}

	return apiTemplate, endpointTemplate, packageTemplate, planTemplate
}
func TransformSwagger(user *ApiUser, swaggerDoc string, sourceFormat string, targetFormat string, oauthToken string) (map[string]interface{}, error) {
	tfSwaggerDoc, err := user.TransformSwagger(string(swaggerDoc), sourceFormat, targetFormat, oauthToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to transform swagger doc\n\n")
	}

	// Only need the value of 'document'. Including the rest will cause errors
	var m map[string]interface{}
	if err = json.Unmarshal([]byte(tfSwaggerDoc), &m); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to process swagger doc\n\n")
	}

	return m, err
}

func MapToByteArray(mapToConvert map[string]interface{}) []byte {
	var convertedByteArray []byte
	var err error

	if val, ok := mapToConvert["document"]; ok {
		mapToConvert = val.(map[string]interface{})
	}

	if convertedByteArray, err = json.Marshal(mapToConvert); err != nil {
		panic(err)
	}

	return convertedByteArray
}

func CreateOrUpdateApi(user *ApiUser, token string, cleanedTfApiSwaggerDoc []byte, mApi map[string]interface{}) (string, string, []interface{}, bool) {
	updated := false

	masheryObject := "services"
	masheryObjectProperties := "id,name,endpoints.id,endpoints.name,endpoints.inboundSslRequired,endpoints.outboundRequestTargetPath,endpoints.outboundTransportProtocol,endpoints.publicDomains,endpoints.requestAuthenticationType,endpoints.requestPathAlias,endpoints.requestProtocol,endpoints.supportedHttpMethods,endoints.systemDomains,endpoints.trafficManagerDomain"
	var apiId string
	var apiName string
	var endpoints [](interface{})

	api, err := user.Read(masheryObject, "name:"+mApi["name"].(string), masheryObjectProperties, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to fetch api\n\n")
		panic(err)
	}

	shortDelay()

	var f [](interface{})
	if err = json.Unmarshal([]byte(api), &f); err != nil {
		panic(err)
	}
	if len(f) == 0 {
		s, err := user.Create(masheryObject, masheryObjectProperties, string(cleanedTfApiSwaggerDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create the api %s\n\n", s)
			panic(err)
		}
		apiId, apiName, endpoints = GetApiDetails(s)

	} else {
		m := f[0].(map[string]interface{})
		var m1 map[string]interface{}
		json.Unmarshal(cleanedTfApiSwaggerDoc, &m1)
		merged := merge(m, m1, 0)
		var mergedDoc []byte
		if mergedDoc, err = json.Marshal(merged); err != nil {
			panic(err)
		}
		serviceId := merged["id"].(string)
		s, err := user.Update(masheryObject+"/"+serviceId, masheryObjectProperties, string(mergedDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to update the api %s\n\n", s)
			panic(err)
		}
		apiId, apiName, endpoints = GetApiDetails(s)

		updated = true
	}

	return apiId, apiName, endpoints, updated
}

func merge(dst, src map[string]interface{}, depth int) map[string]interface{} {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			if reflect.ValueOf(dstVal).Kind() == reflect.Map {
				srcMap, srcMapOk := mapify(srcVal)
				dstMap, dstMapOk := mapify(dstVal)
				if srcMapOk && dstMapOk {
					srcVal = merge(dstMap, srcMap, depth+1)
				}
			} else if (key == "endpoints" || key == "plans") && reflect.ValueOf(dstVal).Kind() == reflect.Slice {
				m_d := dstVal.([]interface{})
				m_s := srcVal.([]interface{})
				items := []map[string]interface{}{}

				for _, d_item := range m_d {
					i_d := d_item.(map[string]interface{})
					var i_s map[string]interface{}
					for _, s_item := range m_s {
						i_s = s_item.(map[string]interface{})
						if i_s["requestPathAlias"] == i_d["requestPathAlias"] {
							i_s2 := merge(i_d, i_s, depth+1)
							items = append(items, i_s2)
						}
					}
				}

				for _, s_item := range m_s {
					i_s := s_item.(map[string]interface{})
					if !MatchingEndpoint(i_s, m_d) {
						items = append(items, i_s)
					}
				}
				srcVal = items
			}
		}

		dst[key] = srcVal
	}
	return dst
}

func MatchingEndpoint(ep map[string]interface{}, epList []interface{}) bool {
	var i_d map[string]interface{}
	for _, d_item := range epList {
		i_d = d_item.(map[string]interface{})
		if i_d["requestPathAlias"] == ep["requestPathAlias"] {
			return true
		}
	}
	return false
}

func mapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}
		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}
		return m, true
	}
	return map[string]interface{}{}, false
}

func CreateOrUpdateIodocs(user *ApiUser, token string, cleanedTfIodocSwaggerDoc []byte, apiId string, updated bool) {
	masheryObject := "iodocs/services"
	masheryObjectProperties := "id"

	item, err := user.Read(masheryObject, "serviceId:"+apiId, masheryObjectProperties, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to fetch iodocs\n\n")
		panic(err)
	}

	var f [](interface{})
	if err = json.Unmarshal([]byte(item), &f); err != nil {
		panic(err)
	}

	shortDelay()

	if len(f) == 0 {
		s, err := user.Create(masheryObject, masheryObjectProperties, string(cleanedTfIodocSwaggerDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create the iodocs %s\n\n", s)
		}
	} else {
		s, err := user.Update(masheryObject+"/"+apiId, masheryObjectProperties, string(cleanedTfIodocSwaggerDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create the iodocs %s\n\n", s)
		}
	}
}

func CreateOrUpdatePackage(user *ApiUser, token string, packagePlanDoc []byte, apiName string, updated bool) string {
	var p string
	masheryObject := "packages"
	masheryObjectProperties := "id,name,plans.id,plans.name"

	item, err := user.Read(masheryObject, "name:"+apiName, masheryObjectProperties, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to fetch package\n\n")
		panic(err)
	}

	var f [](interface{})
	if err = json.Unmarshal([]byte(item), &f); err != nil {
		panic(err)
	}

	if len(f) == 0 {
		p, err = user.Create(masheryObject, masheryObjectProperties, string(packagePlanDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create the package %s\n\n", p)
			panic(err)
		}
	} else {

		m := f[0].(map[string]interface{})

		var m1 map[string]interface{}
		json.Unmarshal(packagePlanDoc, &m1)
		merged := merge(m, m1, 0)
		var mergedDoc []byte
		if mergedDoc, err = json.Marshal(merged); err != nil {
			panic(err)
		}
		packageId := merged["id"].(string)
		p, err = user.Update(masheryObject+"/"+packageId, masheryObjectProperties, string(mergedDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to update the package %s\n\n", p)
			panic(err)
		}
	}
	return p
}

func GetApiDetails(api string) (string, string, []interface{}) {
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(api), &m); err != nil {
		panic(err)
	}
	return m["id"].(string), m["name"].(string), m["endpoints"].([]interface{}) // getting the api id and name
}

func GetPackagePlanDetails(packagePlan string) (string, string) {
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(packagePlan), &m); err != nil {
		panic(err)
	}
	plans := m["plans"].([]interface{})
	plan := plans[0].(map[string]interface{})
	return m["id"].(string), plan["id"].(string) // getting the package id and plan id
}

func UpdateIodocsDataWithApi(ioDoc []byte, apiId string) []byte {
	// need to create a different json representation for an IOdocs post body
	m1 := map[string]interface{}{}
	if err := json.Unmarshal([]byte(string(ioDoc)), &m1); err != nil {
		panic(err)
	}

	var cleanedTfIodocSwaggerDoc []byte

	m := map[string]interface{}{}
	m["definition"] = m1
	m["serviceId"] = apiId
	cleanedTfIodocSwaggerDoc, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return cleanedTfIodocSwaggerDoc
}

func CreatePackagePlanDataFromApi(apiId string, apiName string, endpoints []interface{}) map[string]interface{} {
	pack := map[string]interface{}{}
	pack["name"] = apiName
	pack["sharedSecretLength"] = 10

	plan := map[string]interface{}{}
	plan["name"] = apiName
	plan["selfServiceKeyProvisioningEnabled"] = false
	plan["numKeysBeforeReview"] = 1

	service := map[string]interface{}{}
	service["id"] = apiId

	service["endpoints"] = endpoints

	planServices := []map[string]interface{}{}
	planServices = append(planServices, service)

	plan["services"] = planServices

	plans := []map[string]interface{}{}
	plans = append(plans, plan)
	pack["plans"] = plans

	return pack
}

func CreateApplicationAndKey(user *ApiUser, token string, packagePlan string, apiName string) string {
	var key string
	member, err := user.Read("members", "username:"+user.Username, "id,username,applications,packageKeys", token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to fetch api\n\n")
		panic(err)
	}

	var f [](interface{})
	if err = json.Unmarshal([]byte(member), &f); err != nil {
		panic(err)
	}

	var f_app interface{}
	testApplication := map[string]interface{}{}
	m := f[0].(map[string]interface{})
	var f2 [](interface{})
	f2 = m["applications"].([](interface{}))
	for _, application := range f2 {
		if application.(map[string]interface{})["name"] == "Test Application: "+apiName {
			testApplication = application.(map[string]interface{})
			packageKeys, err := user.Read("applications/"+testApplication["id"].(string)+"/packageKeys", "", "id,apikey,secret", token)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Unable to fetch packagekeys\n\n")
				panic(err)
			}

			var f [](interface{})
			if err = json.Unmarshal([]byte(packageKeys), &f); err != nil {
				panic(err)
			}
			if len(f) > 0 {
				pk := f[0].(map[string]interface{})

				testKeyDoc, err := json.Marshal(pk)
				if err != nil {
					panic(err)
				}
				key = string(testKeyDoc)
			}
			f_app = testApplication
		}
	}

	if len(testApplication) == 0 {
		testApplication["name"] = "Test Application: " + apiName
		testApplication["username"] = user.Username
		testApplication["is_packaged"] = true
		var testApplicationDoc []byte

		testApplicationDoc, err = json.Marshal(testApplication)
		if err != nil {
			panic(err)
		}
		application, err := user.Create("members/"+m["id"].(string)+"/applications", "id,name", string(testApplicationDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create application\n\n")
			panic(err)
		}

		if err = json.Unmarshal([]byte(application), &f_app); err != nil {
			panic(err)
		}

	}

	if key == "" {
		packageId, planId := GetPackagePlanDetails(packagePlan)
		keyToCreate := map[string]interface{}{}
		keyPackage := map[string]interface{}{}
		keyPackage["id"] = packageId
		keyPlan := map[string]interface{}{}
		keyPlan["id"] = planId
		keyToCreate["package"] = keyPackage
		keyToCreate["plan"] = keyPlan
		var testKeyDoc []byte

		testKeyDoc, err = json.Marshal(keyToCreate)
		if err != nil {
			panic(err)
		}
		key, err = user.Create("applications/"+f_app.(map[string]interface{})["id"].(string)+"/packageKeys", "", string(testKeyDoc), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to create key\n\n")
			panic(err)
		}
	}

	return key

}

func GenerateExampleCall(endpoint map[string]interface{}, key string) string {
	var exampleCall string

	public_domains := endpoint["publicDomains"].([]interface{})
	pd_map := public_domains[0].(map[string]interface{})
	var pk map[string]interface{}
	if err := json.Unmarshal([]byte(key), &pk); err != nil {
		panic(err)
	}
	protocol := "https"
	if !endpoint["inboundSslRequired"].(bool) {
		protocol = "http"
	}
	sig := ""
	if endpoint["requestAuthenticationType"] == "apiKeyAndSecret_SHA256" {
		sig = "&sig='$(php -r \"echo hash('sha256', '" + pk["apikey"].(string) + "'.'" + pk["secret"].(string) + "'.time());\")"
	}
	exampleCall = "curl -i -v -k -X " + strings.ToUpper(endpoint["supportedHttpMethods"].([]interface{})[0].(string)) + " '" + protocol + "://" + pd_map["address"].(string) + endpoint["requestPathAlias"].(string) + "?api_key=" + pk["apikey"].(string) + sig
	return exampleCall
}

type Responder func(*http.Request) (*http.Response, error)

type NopTransport struct {
	responders map[string]Responder
}

var DefaultNopTransport = &NopTransport{}

func debug(data []byte, err error) {
	if err == nil {
		logger.Debugf("%s\n\n", data)
	} else {
		logger.Errorf("%s\n\n", err)
		os.Exit(1)
	}
}

func init() {
	DefaultNopTransport.responders = make(map[string]Responder)
}

func (n *NopTransport) RegisterResponder(method, url string, responder Responder) {
	n.responders[method+" "+url] = responder
}

func (n *NopTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()

	// Scan through the responders
	for k, r := range n.responders {
		if k != key {
			continue
		}
		return r(req)
	}

	return nil, errors.New("No responder found")
}

func RegisterResponder(method, url string, responder Responder) {
	DefaultNopTransport.RegisterResponder(method, url, responder)
}

func newHttp(nop bool) *http.Client {
	client := &http.Client{}
	if nop {
		client.Transport = DefaultNopTransport
	}

	return client
}

func setContentType(r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
}

func setOauthToken(r *http.Request, oauthToken string) {
	r.Header.Add("Authorization", "Bearer "+oauthToken)
}

func readBody(body io.Reader) ([]byte, error) {
	bodyText, err := ioutil.ReadAll(body)
	if err != nil {
		return bodyText, err
	}
	return bodyText, nil
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) CreateAPI(tfSwaggerDoc string, oauthToken string) (string, error) {
	return user.CreateUpdate("POST", "services", "", tfSwaggerDoc, oauthToken)
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) Create(resource string, fields string, content string, oauthToken string) (string, error) {
	return user.CreateUpdate("POST", resource, fields, content, oauthToken)
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) CreateUpdate(method string, resource string, fields string, content string, oauthToken string) (string, error) {
	fullUri := masheryUri + restUri + resource
	if fields != "" {
		fullUri = fullUri + "?fields=" + fields
	}
	client := newHttp(user.Noop)
	r, _ := http.NewRequest(method, fullUri, bytes.NewReader([]byte(content)))
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	s := string(bodyText)
	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("Unable to create the api: status code %v", resp.StatusCode)
	}

	return s, err
}

// Read fetch data
func (user *ApiUser) Read(resource string, filter string, fields string, oauthToken string) (string, error) {

	fullUri := masheryUri + restUri + resource
	if fields != "" && filter == "" {
		fullUri = fullUri + "?fields=" + fields
	} else if fields == "" && filter != "" {
		fullUri = fullUri + "?filter=" + filter
	} else {
		fullUri = fullUri + "?fields=" + fields + "&filter=" + filter
	}

	client := newHttp(user.Noop)

	r, _ := http.NewRequest("GET", masheryUri+restUri+resource+"?filter="+filter+"&fields="+fields, nil)
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	s := string(bodyText)
	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("Unable to create the api: status code %v", resp.StatusCode)
	}

	return s, err
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) Update(resource string, fields string, content string, oauthToken string) (string, error) {
	return user.CreateUpdate(http.MethodPut, resource, fields, content, oauthToken)
}

// TransformSwagger sends the swagger doc to Mashery API to be
// transformed into the target format.
func (user *ApiUser) TransformSwagger(swaggerDoc string, sourceFormat string, targetFormat string, oauthToken string) (string, error) {
	// New client
	client := newHttp(user.Noop)

	v := url.Values{}
	v.Set("sourceFormat", sourceFormat)
	v.Add("targetFormat", targetFormat)
	v.Add("publicDomain", user.Portal)

	r, _ := http.NewRequest("POST", masheryUri+restUri+transformUri+"?"+v.Encode(), bytes.NewReader([]byte(swaggerDoc)))
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	if bodyText, err := readBody(resp.Body); err == nil {
		if resp.StatusCode != http.StatusOK {
			return string(bodyText), fmt.Errorf("Unable to transform the swagger doc: status code %v", resp.StatusCode)
		}
		return string(bodyText), nil
	} else {
		return string(bodyText), err
	}
}

// FetchOAuthToken exchanges the creds for an OAuth token
func (user *ApiUser) FetchOAuthToken() (string, error) {
	// New client
	client := newHttp(user.Noop)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", user.Username)
	data.Set("password", user.Password)
	data.Set("scope", user.Uuid)

	r, _ := http.NewRequest("POST", masheryUri+"/v3/token", strings.NewReader(data.Encode()))
	r.SetBasicAuth(user.ApiKey, user.ApiSecretKey)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "*/*")

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	if bodyText, err := readBody(resp.Body); err == nil {
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("Unable to get the OAuth token: status code (%v), message (%v)", resp.StatusCode, string(bodyText))
		}

		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(string(bodyText)), &dat); err != nil {
			return "", errors.New("Unable to unmarshal JSON")
		}

		accessToken, ok := dat[accessToken].(string)
		if !ok {
			return "", errors.New("Invalid json. Expected a field with access_token")
		}

		return accessToken, nil
	} else {
		return string(bodyText), err
	}
}
