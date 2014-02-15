package gohiera

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"regexp"
	"runtime"
	"strings"
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
:ubuntu
*/

var interpolationRegexp *regexp.Regexp = regexp.MustCompile(`%{[a-zA-Z0-9_:]+}`)

type HieraConfig struct {
	RawBackends   interface{}   `yaml:":backends"`
	RawHiearchy   interface{}   `yaml:":hierarchy"`
	Yaml          SubConfig     `yaml:":yaml"`
	Json          SubConfig     `yaml:":json"`
	MergeBehavior MergeBehavior `yaml:":merge_behavior"`
	Backends      []string      `yaml:"-"`
	Hiearchy      []string      `yaml:"-"`
}

type SubConfig struct {
	Datadir string `yaml:":datadir"`
}

type Hiera struct {
	config HieraConfig
}

type MergeBehavior string

const (
	Native MergeBehavior = "native"
	Deep                 = "deep"
	Deeper               = "deeper"
)

// Hiera can only operate in a context of an facter environment.  Even things
// such as the Datadir to act as the root of the yaml/json files can change
// depending upon the "facts".
func LookupValue(key string, facts map[string]string) {
}

// Expand any cases of yaml's string interpolation:
// e.g. "Foo %{variable}"
func ExpandString(str string, facts map[string]string) string {
	return interpolationRegexp.ReplaceAllStringFunc(str, func(interp string) string {
		key := interp[2 : len(interp)-1]
		// A key can start with "::" implying it is a facter fact
		key = strings.TrimPrefix(key, "::")
		if val, ok := facts[key]; !ok {
			return interp
		} else {
			return val
		}
	})
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

	for i, backend := range c.Backends {
		// Force data into consistent form
		backend = strings.ToLower(backend)
		c.Backends[i] = backend

		// Use already lowered version
		if backend != "json" && backend != "yaml" {
			return nil, fmt.Errorf("gohiera does not handle backend: '%s'", backend)
		}

		// Set default backend location (if not set)
		if backend == "json" && c.Json.Datadir == "" {
			if runtime.GOOS == "windows" {
				c.Json.Datadir = "%PROGRAMDATA%\\PuppetLabs\\Hiera\\var"
			} else {
				c.Json.Datadir = "/var/lib/hiera"
			}
		}

		// Set default backend location (if not set)
		if backend == "yaml" && c.Yaml.Datadir == "" {
			if runtime.GOOS == "windows" {
				c.Yaml.Datadir = "%PROGRAMDATA%\\PuppetLabs\\Hiera\\var"
			} else {
				c.Yaml.Datadir = "/var/lib/hiera"
			}
		}
	}

	c.Hiearchy, err = singletonToArray("Hiearchy", c.RawHiearchy)
	if err != nil {
		return nil, err
	}

	if err := validMergeBehavior(c.MergeBehavior); err != nil {
		return nil, err
	}

	h.config = c
	return &h, nil
}

func validMergeBehavior(behavior MergeBehavior) error {
	switch behavior {
	case "":
		// Alias for Native
	case Native:
	case Deep:
	case Deeper:
	default:
		return fmt.Errorf("%s is not a valid behavior", behavior)
	}
	return nil
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
