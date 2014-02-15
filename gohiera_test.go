package gohiera

import "testing"

var PluralConfig []byte
var SingletonConfig []byte
var InvalidHiearchyConfig []byte

func init() {
	// Example config from http://docs.puppetlabs.com/hiera/1/configuring.html
	PluralConfig = []byte(`
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
`)
	SingletonConfig = []byte(`
---
:backends: yaml
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy: common
`)
	InvalidHiearchyConfig = []byte(`
---
:backends: yaml
:yaml:
  :datadir: /etc/puppet/hieradata
:json:
  :datadir: /etc/puppet/hieradata
:hierarchy: 
  - 506
`)
}

func TestHieraFromStringPlural(t *testing.T) {
	res, err := HieraFromString(PluralConfig)

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
	_, err := HieraFromString(InvalidHiearchyConfig)

	if err == nil {
		t.Errorf("Shouldn't have been able to parse InvalidHiearchyConfig but did.")
		return
	}
}

func TestHieraFromStringSingle(t *testing.T) {
	res, err := HieraFromString(SingletonConfig)

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
