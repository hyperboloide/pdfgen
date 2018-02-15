package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	Addr, Port, Root, FakeUrl string
	Templates                 = make(map[string]*Template)
)

func init() {
	viper.SetEnvPrefix("pdfgen")

	viper.BindEnv("port")
	viper.SetDefault("port", "8888")

	viper.BindEnv("addr")
	viper.SetDefault("addr", "0.0.0.0")

	viper.BindEnv("templates")
	viper.SetDefault("templates", "./templates")
}

func configRead() {
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		log.Fatal("executable wkhtmltopdf could not be found in PATH")
	}

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
				if t, err := NewTemplate(Root, i.Name()); err != nil {
					log.Fatal(err)
				} else {
					Templates[i.Name()] = t
				}
			}
		}
	}

	nb := len(Templates)
	switch nb {
	case 0:
		fmt.Println("No template found, exiting.")
		return
	case 1:
		fmt.Printf("1 template found in '%s':\n", Root)
	default:
		fmt.Printf("%d templates found in '%s':\n", nb, Root)
	}
	for k, _ := range Templates {
		fmt.Printf("  - %s\n", k)
	}
}
