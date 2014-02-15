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
*/

type HieraConfig struct {
	Backends []string  `yaml:":backends"`
	Yaml     SubConfig `yaml:":yaml"`
	Json     SubConfig `yaml:":json"`
	Hiearchy []string  `yaml:":hierarchy"`
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
	h := Hiera{}

	if err := goyaml.Unmarshal(config, &c); err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", c)

	return &h, nil
}
