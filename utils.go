package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
	"github.com/spf13/viper"
)

func readConfig(filename, configPath string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	return v, err
}

func readInstancesFile(r io.Reader) ([]Service, error) {
	scanner := bufio.NewScanner(r)
	services := make([]Service, 0)
	idx := 0

	for scanner.Scan() {
		if idx == 0 {
			idx++
			continue
		}
		instanceTextInfo := scanner.Text()
		if strings.HasPrefix(instanceTextInfo, "#") {
			continue
		}
		instance, err := parseInstanceInfo(instanceTextInfo)
		if err != nil {
			continue
		}
		svc := Service{Instance: instance, State: NOTCHECKEDYET}
		services = append(services, svc)
		idx++
	}
	return services, nil
}

func parseInstanceInfo(text string) (Instance, error) {
	instanceInfoTokens := strings.Split(text, ",")
	if len(instanceInfoTokens) != RequiredNumberFieldsInInstanceFile {
		return Instance{}, fmt.Errorf("expecting %d comma separated fields in `%s`", RequiredNumberFieldsInInstanceFile, text)
	}

	instance := Instance{}
	instance.Type = strings.TrimSpace(instanceInfoTokens[0])
	instance.ShortName = strings.TrimSpace(instanceInfoTokens[1])
	instance.Host = strings.TrimSpace(instanceInfoTokens[2])
	instance.Instance = strings.TrimSpace(instanceInfoTokens[3])
	port, err := strconv.Atoi(strings.TrimSpace(instanceInfoTokens[4]))
	if err != nil {
		instance.Port = -1
		return instance, err
	}

	instance.Port = port

	instance.URL = strings.TrimSpace(instanceInfoTokens[5])
	instance.Category = strings.TrimSpace(instanceInfoTokens[6])
	isSecure, err := strconv.ParseBool(strings.TrimSpace(instanceInfoTokens[7]))
	if err != nil {
		return Instance{}, err
	}
	instance.Secure = isSecure

	return instance, nil
}

func binaryPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func toDurationSeconds(every int) time.Duration {
	everyDuration := time.Duration(every) * time.Second
	return everyDuration
}

func buildServiceURL(service *Service) string {
	var url strings.Builder

	if service.Instance.Secure {
		url.WriteString("https://")
	} else {
		url.WriteString("http://")
	}

	url.WriteString(service.Instance.Host)
	url.WriteString(fmt.Sprintf(":%d/%s", service.Instance.Port, service.Instance.URL))

	return url.String()
}

func callService(service *Service, config *Configuration) error {
	var client *http.Client
	if service.Instance.Secure {
		client = config.httpsClient
	} else {
		client = config.httpClient
	}

	url := buildServiceURL(service)
	log.Printf("About to call -> [%s]\n", url)
	resp, err := client.Get(url)
	if err != nil {
		service.State = DOWN
		// Error calling the service, finish here
		service.Checking = false
		return err
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= 299 {
		service.State = UP
	} else {
		service.State = DOWN
	}
	serviceResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		service.Checking = false
		return err
	}
	defer resp.Body.Close()
	log.Printf("Service response {%s} -> [%s]\n", *service, serviceResponse)
	service.Checking = false
	service.LastChecked = time.Now()

	return nil
}

func check(service *Service, config *Configuration) {
	// Is the service being checked?
	if service.Checking {
		return
	}
	// If not check if we have already checked the service a few seconds ago ...
	// We need to give the service a little breath ...
	every := time.Duration(config.envConfig.GetInt(CheckServicesEverySeconds)) * time.Second
	if time.Since(service.LastChecked) < every {
		return
	}
	// If OK then call the service and update the properties ...
	err := callService(service, config)
	if err != nil {
		log.Printf("Error calling service [%s], error: %s\n", *service, err)
		service.State = DOWN
		service.LastChecked = time.Now()
	}
}

func checkServices(services *[]Service, config *Configuration) {
	for i := 0; i < len(*services); i++ {
		go check(&(*services)[i], config)
	}
}

func createHTTPClient(config *viper.Viper) *http.Client {
	maxIdleConnections := config.GetInt(MaxIdleConnections)
	requestTimeout := config.GetInt(RequestTimeout)
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}

	return client
}

func createHTTPSClient(config *viper.Viper) *http.Client {
	maxIdleConnections := config.GetInt(MaxIdleConnections)
	requestTimeout := config.GetInt(RequestTimeout)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}

	return client
}

func updateScreen(services *[]Service, config *Configuration, textEdit *tui.TextEdit, ui tui.UI) {
	printEvery := config.envConfig.GetInt(PrintStatusEverySeconds)
	checkServicesEvery := config.envConfig.GetInt(CheckServicesEverySeconds)

	printStatusTicker := time.NewTicker(toDurationSeconds(printEvery))
	checkServicesTicker := time.NewTicker(toDurationSeconds(checkServicesEvery))

	for {
		select {
		case <-printStatusTicker.C:
			printServiceStatusInScreen(textEdit, ui, services)
		case <-checkServicesTicker.C:
			go checkServices(services, config)
		}
	}
}

func groupByInstance(services *[]Service) *map[string][]Service {
	byInstance := make(map[string][]Service, 0)

	for _, svc := range *services {
		byInstance[svc.Instance.Instance] = append(byInstance[svc.Instance.Instance], svc)
	}

	return &byInstance
}
