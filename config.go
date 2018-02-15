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
	// Templates is a map containing all templates,
	// the keys are used to match urls with directories.
	Templates map[string]*Template
)

func init() {
	viper.SetEnvPrefix("pdfgen")
	viper.AutomaticEnv()
	viper.SetDefault("port", "8888")
	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("templates", "")
}

// IsValidDir returns true if the path exists and is a directory.
func IsValidDir(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return fi.Mode().IsDir()
}

// SelectDir selects the first valid dir.
func SelectDir(choices []string) *string {
	for _, path := range choices {
		if IsValidDir(path) {
			return &path
		}
	}
	return nil
}

// BuildTemplates will iterate all dirs one by one to create the corresponding
// templates.
func BuildTemplates(rootDir string) error {
	fi, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return fmt.Errorf("failed to read templates dir '%s'", rootDir)
	}
	Templates = make(map[string]*Template)
	for _, i := range fi {
		if i.IsDir() && i.Name()[0] != '.' {
			t, err := NewTemplate(rootDir, i.Name())
			if err != nil {
				return err
			}
			Templates[i.Name()] = t
		}
	}
	return nil
}

// FindRoot choose the templates root.
func FindRoot() (string, error) {
	rootDir := viper.GetString("templates")
	choices := []string{
		"/etc/pdfgen/templates",
		os.Getenv("HOME") + "/.templates",
	}

	if rootDir != "" {
		if !IsValidDir(rootDir) {
			return "", errors.New("invalid template directory")
		}
		return rootDir, nil
	}
	if pth := SelectDir(choices); pth != nil {
		return *pth, nil
	}
	return "", errors.New("template directory not found")
}

// ConfigRead validates the configuration.
func ConfigRead() error {
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		return errors.New("executable wkhtmltopdf could not be found in PATH")
	} else if rootDir, err := FindRoot(); err != nil {
		return err
	} else if rootDir, err = filepath.Abs(rootDir); err != nil {
		return fmt.Errorf("invalid templates dir '%s'", rootDir)
	} else if err := BuildTemplates(rootDir); err != nil {
		return err
	} else {
		nb := len(Templates)
		switch nb {
		case 0:
			return errors.New("No template found")
		case 1:
			fmt.Printf("1 template found in '%s':\n", rootDir)
		default:
			fmt.Printf("%d templates found in '%s':\n", nb, rootDir)
		}
		for k := range Templates {
			fmt.Printf("  - %s\n", k)
		}
		return nil
	}
}
