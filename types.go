package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// Instance ...
type Instance struct {
	Type      string
	ShortName string
	Host      string
	Instance  string
	Port      int
	URL       string
	Category  string
	Secure    bool
}

// Service ...
type Service struct {
	Instance    Instance
	LastChecked time.Time
	Checking    bool
	State       State
}

// State ...
type State int

const (
	// NOTCHECKEDYET ...
	NOTCHECKEDYET State = 0
	// UP ...
	UP State = 1
	// DOWN ...
	DOWN State = 2
	// RESTARING ...
	RESTARING State = 3
)

// Configuration ...
type Configuration struct {
	envConfig   *viper.Viper
	httpClient  *http.Client
	httpsClient *http.Client
}

func (svc Service) String() string {
	return fmt.Sprintf("%s:%d/%s, [%s]", svc.Instance.Host, svc.Instance.Port, svc.Instance.URL, svc.State)
}

func (svc Service) ShortString() string {
	return fmt.Sprintf("%s [%s],", svc.Instance.ShortName, svc.State)
}
