package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	Root, FakeUrl string
	Templates     map[string]*Template
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

func ConfigRead() error {
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		return errors.New("executable wkhtmltopdf could not be found in PATH")
	}

	choices := []string{
		"/etc/pdfgen/templates",
		os.Getenv("HOME") + "/.templates",
	}

	if viper.GetString("templates") != "" {
		Root = viper.GetString("templates")
		if !IsValidTemplateDir(Root) {
			return errors.New("invalid template directory")
		}
	} else if pth := SelectDir(choices); pth != nil {
		Root = *pth
	} else {
		return errors.New("template directory not found")
	}

	var err error
	if Root, err = filepath.Abs(Root); err != nil {
		return fmt.Errorf("invalid templates dir '%s'", Root)
	}

	if fi, err := ioutil.ReadDir(Root); err != nil {
		return fmt.Errorf("failed to read templates dir '%s'", Root)
	} else {
		Templates = make(map[string]*Template)
		for _, i := range fi {
			if i.IsDir() && i.Name()[0] != '.' {
				if t, err := NewTemplate(Root, i.Name()); err != nil {
					return err
				} else {
					Templates[i.Name()] = t
				}
			}
		}
	}

	nb := len(Templates)
	switch nb {
	case 0:
		return errors.New("No template found")
	case 1:
		fmt.Printf("1 template found in '%s':\n", Root)
	default:
		fmt.Printf("%d templates found in '%s':\n", nb, Root)
	}
	for k, _ := range Templates {
		fmt.Printf("  - %s\n", k)
	}
	return nil
}
