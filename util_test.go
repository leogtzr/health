package main

import (
	"testing"
)

func TestShouldParseTextInfo(t *testing.T) {
	type test struct {
		text         string
		instanceType string
		shortName    string
		host         string
		instance     string
		port         int
		url          string
		category     string
		secure       bool
	}

	tests := []test{
		test{
			text:         `APP,XL00,pr-camaro1-xl00,EStore01,8443,monitoring/healthcheck.jsp,ATG,TRUE`,
			instanceType: "APP",
			shortName:    "XL00",
			host:         "pr-camaro1-xl00",
			instance:     "EStore01",
			port:         8443,
			url:          "monitoring/healthcheck.jsp",
			category:     "ATG",
			secure:       true,
		},
	}

	for _, tc := range tests {
		got, err := parseInstanceInfo(tc.text)
		if err != nil {
			t.Errorf(err.Error())
		}
		if got.Type != tc.instanceType {
			t.Errorf("got=[%s], expected=[%s]", got.Type, tc.instanceType)
		}
		if got.ShortName != tc.shortName {
			t.Errorf("got=[%s], expected=[%s]", got.ShortName, tc.shortName)
		}
		if got.Host != tc.host {
			t.Errorf("got=[%s], expected=[%s]", got.Host, tc.host)
		}
		if got.Instance != tc.instance {
			t.Errorf("got=[%s], expected=[%s]", got.Instance, tc.instance)
		}
		if got.Port != tc.port {
			t.Errorf("got=[%d], expected=[%d]", got.Port, tc.port)
		}
		if got.URL != tc.url {
			t.Errorf("got=[%s], expected=[%s]", got.URL, tc.url)
		}
		if got.Category != tc.category {
			t.Errorf("got=[%s], expected=[%s]", got.Category, tc.category)
		}
		if got.Secure != tc.secure {
			t.Errorf("got=[%t], expected=[%t]", got.Secure, tc.secure)
		}

	}

}

func TestShouldParseTextInfo2(t *testing.T) {
	type test struct {
		text              string
		shouldParsingFail bool
	}

	tests := []test{
		test{
			text:              `APP,XL00,pr-camaro1-xl00,EStore01,8443,monitoring/healthcheck.jsp,ATG,TRUE`,
			shouldParsingFail: false,
		},
		test{
			text:              `APP,XL00,pr-camaro1-xl00,EStore01,8443,monitoring/healthcheck.jsp,ATG`,
			shouldParsingFail: true,
		},
		test{
			text:              `APP,XL00,pr-camaro1-xl00,EStore01,XXX,monitoring/healthcheck.jsp,ATG,TRUE`,
			shouldParsingFail: true,
		},
	}

	for _, tc := range tests {
		_, err := parseInstanceInfo(tc.text)
		hasFailed := err != nil
		if hasFailed != tc.shouldParsingFail {
			t.Errorf("parsing for '%s' should have failed", tc.text)
		}
	}
}
