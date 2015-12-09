package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Template struct {
	rootDir string
	url     string

	template *template.Template

	Grayscale       bool   `json:"grayscale"`
	PageSize        string `json:"page_size"`
	JavascriptDelay int    `json:"javascript_delay"`
	Footer          bool   `json:"footer"`
	Header          bool   `json:"header"`
}

func (t *Template) BuildParams(url string) []string {
	params := []string{}
	if t.Grayscale {
		params = append(params, "--grayscale")
	}
	if t.PageSize != "" {
		params = append(params, "--page-size", t.PageSize)
	}
	if t.JavascriptDelay > 0 {
		params = append(params, "--javascript-delay", strconv.Itoa(t.JavascriptDelay))
	}
	if t.Footer {
		params = append(params, "--footer-html", fmt.Sprintf("%s/footer", url))
	}
	if t.Header {
		params = append(params, "--header-html", fmt.Sprintf("%s/header", url))
	}
	return params
}

func (t *Template) Gen(sessionId string, w io.Writer) error {
	url := fmt.Sprintf("%s/%s/%s", FakeUrl, t.url, sessionId)
	params := t.BuildParams(url)
	destDir, err := ioutil.TempDir("", "pdfgen")
	if err != nil {
		log.Fatalf("failed to create tmp directory")
	} else {
		defer os.RemoveAll(destDir)
	}
	output := filepath.Join(destDir, "output.pdf")
	params = append(
		params,
		fmt.Sprintf("%s/index", url),
		output)

	cmd := exec.Command("wkhtmltopdf", params...)
	cmd.Dir = destDir
	if o, err := cmd.CombinedOutput(); err != nil {
		if s, exists := Sessions.Get(sessionId); !exists {
			log.Printf("failed to generate pdf for template: '%s', session empty", t.url)
		} else {
			log.Printf("failed to generate pdf for template: '%s' with data: %+v", t.url, s)
		}
		log.Printf("wkhtmltopdf says: \n%s", string(o[:]))
		return err
	}

	r, err := os.OpenFile(output, os.O_RDONLY, 0400)
	if err != nil {
		log.Printf("failed to open generated pdf for template: '%s'", t.url)
		return err
	} else {
		defer r.Close()
	}
	if _, err := io.Copy(w, r); err != nil {
		log.Printf("failed to send generated pdf for template: '%s'", t.url)
		return err
	}
	return nil
}

func NewTemplate(root, path string) *Template {
	t := &Template{
		url:     path,
		rootDir: filepath.Join(root, path),
	}
	fi, err := ioutil.ReadDir(t.rootDir)
	if err != nil {
		log.Fatalf("failed to open template dir '%s'", t.rootDir)
	}

	for _, f := range fi {
		if f.Name() == "config.json" {
			configPath := filepath.Join(t.rootDir, f.Name())
			if b, err := ioutil.ReadFile(configPath); err != nil {
				log.Fatalf("failed to read template config '%s'", configPath)
			} else if err := json.Unmarshal(b, t); err != nil {
				log.Fatalf("invalid JSON for template config '%s'", configPath)
			}
			break
		}
	}
	t.template, err = template.ParseGlob(filepath.Join(t.rootDir, "**.html"))
	if err != nil {
		log.Println(err)
		log.Fatalf("no valid html file found in dir '%s'", t.rootDir)
	}
	if t.template.Lookup("index") == nil {
		log.Fatalf("template '%s' must have an index", t.url)
	}
	if t.Footer && t.template.Lookup("footer") == nil {
		log.Fatalf("template '%s' do not define a footer", t.url)
	}
	if t.Header && t.template.Lookup("header") == nil {
		log.Fatalf("template '%s' do not define a header", t.url)
	}

	return t
}
