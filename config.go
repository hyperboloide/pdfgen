package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	Root, FakeUrl string
	Templates     = make(map[string]*Template)
)

func init() {
	viper.SetEnvPrefix("pdfgen")
	viper.AutomaticEnv()

	viper.SetDefault("port", "8888")
	viper.SetDefault("addr", "0.0.0.0")

	viper.BindEnv("templates")
}

func IsValidTemplateDir(path string) bool {
	if fi, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return fi.Mode().IsDir()
	}
}

func SelectDir(choices []string) *string {
	for _, path := range choices {
		if IsValidTemplateDir(path) {
			return &path
		}
	}
	return nil
}

func ConfigRead() {
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		log.Fatal("executable wkhtmltopdf could not be found in PATH")
	}

	choices := []string{
		"./templates",
		"/etc/pdfgen/templates",
		os.Getenv("HOME") + "/.templates",
	}

	if viper.GetString("templates") != "" && IsValidTemplateDir(viper.GetString("templates")) {
		Root = viper.GetString("templates")
	} else if pth := SelectDir(choices); pth != nil {
		Root = *pth
	} else {
		log.Fatal("template directory not found!")
	}

	var err error
	if Root, err = filepath.Abs(Root); err != nil {
		log.Fatalf("invalid templates dir '%s'", Root)
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
