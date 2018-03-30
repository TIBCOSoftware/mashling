package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

const (
	registerURI   = "/v1/agent/service/register"
	deRegisterURI = "/v1/agent/service/deregister/"
)

type ServiceDefinition struct {
	Name    string `json:"Name"`
	Port    string `json:"port"`
	Address string `json:"address"`
}

// getLocalIP gets the public ip address of the system
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

//RegisterWithConsul registers suplied gateway.json services with consul
func RegisterWithConsul(consulServices []ServiceDefinition, consulToken string, consulDefDir string, consulAddress string) error {
	var err error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	var localIP = getLocalIP()

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	for _, content := range consulServices {

		port, _ := strconv.Atoi(content.Port)

		checkMap := map[string]interface{}{
			"tcp":      localIP + ":" + content.Port,
			"interval": "30s",
			"timeout":  "1s",
		}

		contentMap := map[string]interface{}{
			"Name":    content.Name,
			"Address": localIP,
			"Port":    port,
			"check":   checkMap,
		}

		contentPayload, err := json.MarshalIndent(&contentMap, "", "    ")
		if err != nil {
			return err
		}

		fullURI := "http://" + consulAddress + registerURI

		if len(consulDefDir) != 0 {

			err := os.Chdir(consulDefDir)
			if err != nil {
				return err
			}

			file, err := os.Create(content.Name + ".json")
			defer file.Close()

			if err != nil {
				return err
			}

			serviceMap := map[string]interface{}{
				"service": contentMap,
			}

			serviceContentPayload, err := json.MarshalIndent(&serviceMap, "", "    ")
			if err != nil {
				return err
			}
			_, dataErr := file.Write(serviceContentPayload)
			if dataErr != nil {
				return dataErr
			}

		} else {

			statusCode, err := callConsulService(fullURI, []byte(contentPayload), consulToken)

			if err != nil {
				return err
			}

			if statusCode != http.StatusOK {
				return fmt.Errorf("registration failed : status code %v", statusCode)
			}
		}
	}

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	fmt.Println("===================================")
	fmt.Println("Successfully registered with consul")
	fmt.Println("===================================")
	return nil
}

//DeregisterFromConsul removes suplied gateway.json services from consul
func DeregisterFromConsul(consulServices []ServiceDefinition, consulToken string, consulDefDir string, consulAddress string) error {
	var err error

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	for _, content := range consulServices {

		fullURI := "http://" + consulAddress + deRegisterURI + content.Name

		if len(consulDefDir) != 0 {

			err := os.Chdir(consulDefDir)
			if err != nil {
				return err
			}

			err = os.Remove(content.Name + ".json")
			if err != nil {
				return err
			}

		} else {
			statusCode, cErr := callConsulService(fullURI, []byte(""), consulToken)

			if cErr != nil {
				return cErr
			}

			if statusCode != http.StatusOK {
				return fmt.Errorf("deregistration failed : status code %v", statusCode)
			}
		}
	}

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	fmt.Println("======================================")
	fmt.Println("Successfully de-registered with consul")
	fmt.Println("======================================")
	return nil
}

/**
callConsulService Performs PUT API call on consul agent
**/
func callConsulService(uri string, payload []byte, consulToken string) (int, error) {

	client := &http.Client{}
	r, _ := http.NewRequest("PUT", uri, bytes.NewReader([]byte(payload)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Consul-Token", consulToken)

	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, err
}

/**
reloadConsul used to reload consul services
**/
func reloadConsul(consulSecurityToken string) error {

	command := exec.Command("consul", "reload", "-token="+consulSecurityToken)
	err := command.Run()

	if err != nil {
		return fmt.Errorf("command error output [%v]", err)
	}

	return nil
}
