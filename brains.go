package main

// This is purely for loading Brain data
// State is managed by the BrainManager

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

//Brains is an array of Brain
type Brains struct {
	Brains map[string]Brain `json:"brains"`
}

// Brain is a template for starting and stopping a brain
type Brain struct {
	ID    string `yaml:"id",json:"id"`
	Name  string `yaml:"name",json:"name"`
	Start string `yaml:"start",json:"start"`
	Stop  string `yaml:"stop",json:"stop"`
}

// Initialise the brains
func (brains *Brains) Initialise() (err error) {
	// read the brains dir
	brains.Brains = make(map[string]Brain)
	dirs, err := ioutil.ReadDir("./brains")

	if err != nil {
		return fmt.Errorf("no ./brains")
	}

	for _, d := range dirs {
		if !strings.HasPrefix(".", d.Name()) {
			brainfilename := path.Join("./brains", d.Name(), "brain.yaml")
			brainfile, err := ioutil.ReadFile(brainfilename)
			if err != nil {
				return err
			}
			b := new(Brain)
			err = yaml.Unmarshal(brainfile, &b)
			brains.Brains[b.ID] = *b
		}
	}

	return nil
}
