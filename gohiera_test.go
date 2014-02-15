package gohiera

import "testing"

var PluralConfig []byte
var SingletonConfig []byte

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
}

func TestHieraFromString(t *testing.T) {
	_, err := HieraFromString(PluralConfig)

	if err != nil {
		t.Errorf("Error reading hiera config: %v", err)
	} else {
		t.Log("Succcess parsing PluralConfig")
	}

	_, err = HieraFromString(SingletonConfig)

	if err != nil {
		t.Errorf("Error reading hiera config: %v", err)
	} else {
		t.Log("Success parsing SingletonConfig")
	}
}
