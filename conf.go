package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type myData struct {
	Conf struct {
		Languages []string `yaml:"languages"`
	}
}

func readConf(filename string) (*myData, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &myData{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
