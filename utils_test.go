package main

import (
	"strings"
	"testing"
	"time"
)

func Test_readInstancesFile(t *testing.T) {
	file := strings.NewReader(`
TYPE,SHORT NAME,HOST,INSTANCE,PORT,URL,CATEGORY,SECURE
APP,XL00,localhost,EStore03,8080,monitoring/healthcheck,ATG,FALSE
# ...
APP,XL00,localhost,EStore03,8081,monitoring/healthcheck,ATG,FALSE
APP,XL00,localhost,EStore03,8082,monitoring/healthcheck,ATG,FALSE
`)

	expectedNumberOfServicesInFile := 3
	services, err := readInstancesFile(file)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(services) != expectedNumberOfServicesInFile {
		t.Errorf("expecting %d services in file, got %d", expectedNumberOfServicesInFile, len(services))
	}

	type test struct {
		Type      string
		ShortName string
		Host      string
		Instance  string
		Port      int
		URL       string
		Category  string
		Secure    bool
	}

	tests := []test{
		test{
			Type:      "APP",
			ShortName: "XL00",
			Host:      "localhost",
			Instance:  "EStore03",
			Port:      8080,
			URL:       "monitoring/healthcheck",
			Category:  "ATG",
			Secure:    false,
		},
		test{
			Type:      "APP",
			ShortName: "XL00",
			Host:      "localhost",
			Instance:  "EStore03",
			Port:      8081,
			URL:       "monitoring/healthcheck",
			Category:  "ATG",
			Secure:    false,
		},
		test{
			Type:      "APP",
			ShortName: "XL00",
			Host:      "localhost",
			Instance:  "EStore03",
			Port:      8082,
			URL:       "monitoring/healthcheck",
			Category:  "ATG",
			Secure:    false,
		},
	}

	for i, tc := range tests {
		got := services[i]
		if got.Instance.Type != tc.Type {
			t.Errorf("got=[%s], expected=[%s]", got.Instance.Type, tc.Type)
		}

		if got.Instance.ShortName != tc.ShortName {
			t.Errorf("got=[%s], expected=[%s]", got.Instance.ShortName, tc.ShortName)
		}

		if got.Instance.Host != tc.Host {
			t.Errorf("got=[%s], expected=[%s]", got.Instance.Host, tc.Host)
		}

		if got.Instance.Instance != tc.Instance {
			t.Errorf("got=[%s], expected=[%s]", got.Instance.Instance, tc.Instance)
		}

		if got.Instance.Port != tc.Port {
			t.Errorf("got=[%d], expected=[%d]", got.Instance.Port, tc.Port)
		}

		if got.Instance.URL != tc.URL {
			t.Errorf("got=[%s], expected=[%s]", got.Instance.URL, tc.URL)
		}

		if got.Instance.Category != tc.Category {
			t.Errorf("got=[%s], expected=[%s]", got.Instance.Category, tc.Category)
		}

		if got.Instance.Secure != tc.Secure {
			t.Errorf("got=[%t], expected=[%t]", got.Instance.Secure, tc.Secure)
		}
	}

}

func Test_buildServiceURL(t *testing.T) {

	type test struct {
		svc         Service
		expectedURL string
	}

	tests := []test{
		test{
			svc: Service{Instance: Instance{
				Type:      "APP",
				ShortName: "XL00",
				Host:      "localhost",
				Instance:  "EStore03",
				Port:      8082,
				URL:       "monitoring/healthcheck",
				Category:  "ATG",
				Secure:    false,
			}},
			expectedURL: "http://localhost:8082/monitoring/healthcheck",
		},

		test{
			svc: Service{Instance: Instance{
				Type:      "APP",
				ShortName: "XL00",
				Host:      "pr-blazercamaro-xl20",
				Instance:  "EStore03",
				Port:      8083,
				URL:       "healthcheck",
				Category:  "ATG",
				Secure:    false,
			}},
			expectedURL: "http://pr-blazercamaro-xl20:8083/healthcheck",
		},

		test{
			svc: Service{Instance: Instance{
				Type:      "APP",
				ShortName: "XL00",
				Host:      "bravada",
				Instance:  "EStore03",
				Port:      8084,
				URL:       "healthcheck",
				Category:  "ATG",
				Secure:    true,
			}},
			expectedURL: "https://bravada:8084/healthcheck",
		},
	}

	for _, tc := range tests {
		got := buildServiceURL(&tc.svc)
		if got != tc.expectedURL {
			t.Errorf("got=[%s], expected=[%s]", got, tc.expectedURL)
		}
	}

}

func Test_toDurationSeconds(t *testing.T) {

	type test struct {
		seconds  int
		duration time.Duration
	}

	tests := []test{
		test{
			seconds:  5,
			duration: (5 * time.Second),
		},
	}

	for _, tc := range tests {
		got := toDurationSeconds(tc.seconds)
		if got != tc.duration {
			t.Errorf("got=[%v], expected=[%v]", got, tc.duration)
		}
	}

}
