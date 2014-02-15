package gohiera

import "testing"

func TestHieraFromStringPlural(t *testing.T) {
	config := []byte(`
---
:backends:
  - yaml
  - json
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy:
  - "%{::clientcert}"
  - "%{::custom_location}"
  - common
:merge_behavior: native
`)
	res, err := HieraFromString(config)

	if err != nil {
		t.Errorf("Error reading hiera config: %v", err)
		return
	} else {
		t.Log("Succcess parsing")
	}
	t.Logf("%#v", res.config)

	if len(res.config.Hiearchy) != 3 {
		t.Errorf("Error decoding hiearchies")
	}
}

func TestInvalidHiearchyConfig(t *testing.T) {
	config := []byte(`
---
:backends: yaml
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy:
  - 506
`)
	_, err := HieraFromString(config)

	if err == nil {
		t.Errorf("Failed to detect error in InvalidHiearchyConfig")
		return
	}
}

func TestHieraFromStringSingle(t *testing.T) {
	config := []byte(`
---
:backends: yaml
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy: common
`)
	res, err := HieraFromString(config)

	if err != nil {
		t.Errorf("Error reading hiera config: %v", err)
		return
	} else {
		t.Log("Success parsing SingletonConfig")
	}
	t.Logf("%#v", res.config)

	if len(res.config.Hiearchy) != 1 {
		t.Errorf("Error decoding hiearchies")
	}
}

func TestInvalidBackend(t *testing.T) {
	config := []byte(`
---
:backends: foobar
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy:
  - common
`)
	res, err := HieraFromString(config)

	if err == nil {
		t.Errorf("Failed to detect error in InvalidBackend: '%v'", res.config.Backends[0])
	}
}

func TestInvalidMergeBehavior(t *testing.T) {
	config := []byte(`
---
:backends: yaml
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy:
  - common
:merge_behavior: foobar
`)
	res, err := HieraFromString(config)

	if err == nil {
		t.Errorf("Failed to detect error in InvalidMergeBehavior: '%v'", res.config.MergeBehavior)
	}
}

func TestExpandString(t *testing.T) {
	facts := map[string]string{"environment": "production"}

	tests := [][]string{
		[]string{"/%{environment}/", "/production/"},
		[]string{"/%{missingkey}/%{environment}/", "/%{missingkey}/production/"},
		[]string{"/%{::missingkey}/%{::environment}/", "/%{::missingkey}/production/"},
		[]string{"%{}", "%{}"},
		[]string{"Test %{ a partial match", "Test %{ a partial match"},
		[]string{"A normal String", "A normal String"},
		[]string{"", ""},
	}

	for _, test := range tests {
		start, finish := test[0], test[1]
		if out := ExpandString(start, facts); out != finish {
			t.Errorf("Invalid expansion: '%s' expected '%s'", out, finish)
		}
	}
}
