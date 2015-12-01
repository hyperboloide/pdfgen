package main

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	Addr, Port, Root string
	Templates        = make(map[string]*Template)
)

func init() {
	viper.SetEnvPrefix("pdfgen")

	viper.BindEnv("port")
	viper.SetDefault("port", "8888")

	viper.BindEnv("addr")
	viper.SetDefault("addr", "0.0.0.0")

	viper.BindEnv("templates")
	viper.SetDefault("templates", ".")
}

func configRead() {
	Addr = viper.GetString("addr")
	Port = viper.GetString("port")

	path := viper.GetString("templates")
	var err error
	if Root, err = filepath.Abs(path); err != nil {
		log.Fatalf("invalid templates dir '%s'", path)
	}

	if fi, err := ioutil.ReadDir(Root); err != nil {
		log.Fatalf("failed to read templates dir '%s'", Root)
	} else {
		for _, i := range fi {
			if i.IsDir() && i.Name()[0] != '.' {
				Templates[i.Name()] = NewTemplate(Root, i.Name())
			}
		}
	}
}
