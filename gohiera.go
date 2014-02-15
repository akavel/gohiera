package gohiera

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

/*
The format of the hiera.yaml file is defined at:
  http://docs.puppetlabs.com/hiera/1/configuring.html

The hiera config file *must* be a yaml file.
The keys *must* start with ':'
Backends can be either a []strings or a string
Hiearchy can be either a []strings or a string

Does not parse (or deal with):
:logger
*/

type HieraConfig struct {
	RawBackends interface{} `yaml:":backends"`
	RawHiearchy interface{} `yaml:":hierarchy"`
	Yaml        SubConfig   `yaml:":yaml"`
	Json        SubConfig   `yaml:":json"`
	Backends    []string    `yaml:"-"`
	Hiearchy    []string    `yaml:"-"`
}

type SubConfig struct {
	Datadir string `yaml:":datadir"`
}

type Hiera struct {
	config HieraConfig
}

func LoadHiera(configFile string) (*Hiera, error) {
	contents, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	h, err := HieraFromString(contents)

	return h, err
}

func HieraFromString(config []byte) (*Hiera, error) {
	var c HieraConfig
	var err error
	h := Hiera{}

	if err := goyaml.Unmarshal(config, &c); err != nil {
		return nil, err
	}

	c.Backends, err = singletonToArray("Backends", c.RawBackends)
	if err != nil {
		return nil, err
	}
	c.Hiearchy, err = singletonToArray("Hiearchy", c.RawHiearchy)
	if err != nil {
		return nil, err
	}

	h.config = c
	return &h, nil
}

// Puppet allows for some fields to be either a string or an array of strings.
// We are best served by always treating the standalone string as an array of
// strings which only contains one value.
// singletonToArray as such gets a interface{} and converts it to a []string
// If the value in interface{} can not be converted to a []string it will
// return an error
func singletonToArray(fieldName string, input interface{}) ([]string, error) {
	switch typedInput := input.(type) {
	case string:
		return []string{typedInput}, nil
	case []interface{}:
		var arr []string

		// Convert each item in the array of interface{} into strings:
		for _, value := range typedInput {
			str, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf(`unable to parse '%s' "%v"`, fieldName, value)
			}
			arr = append(arr, str)
		}
		return arr, nil
	default:
		return nil, fmt.Errorf(`Unable to parse '%s': "%v"`, fieldName, typedInput)
	}
}
