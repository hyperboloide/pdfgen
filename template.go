package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/flosch/pongo2"
)

type Template struct {
	RootDir string
	Url     string
	Index   *pongo2.Template
	Footer  *pongo2.Template
}

func (t *Template) BuildParams(url string) []string {
	params := []string{fmt.Sprintf("%s/main", url)}

	if t.Footer != nil {
		params = append(params, "--footer-html", fmt.Sprintf("%s/footer", url))
	}
	params = append(params, "-")
	return params
}

func (t *Template) WritePDF(baseUrl string, w io.Writer) error {
	params := t.BuildParams(baseUrl)
	cmd := exec.Command("wkhtmltopdf", params...)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		_, err := io.Copy(w, output)
		errChan <- err
	}()
	if err := cmd.Start(); err != nil {
		return err
	} else if err := cmd.Wait(); err != nil {
		return err
	}
	return <-errChan
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewTemplate(root, path string) *Template {
	t := &Template{
		Url:     path,
		RootDir: filepath.Join(root, path),
	}

	// fi, err := ioutil.ReadDir(t.RootDir)
	// if err != nil {
	// 	log.Fatalf("failed to open template dir '%s'", t.RootDir)
	// }

	indexPath := filepath.Join(t.RootDir, "index.html")
	if !fileExists(indexPath) {
		log.Fatal("template index '%s' not found.")
	} else if tmpl, err := pongo2.FromFile(indexPath); err != nil {
		log.Fatal(err)
	} else {
		t.Index = tmpl
	}

	footerPath := filepath.Join(t.RootDir, "footer.html")
	if !fileExists(footerPath) {
		return t
	} else if tmpl, err := pongo2.FromFile(footerPath); err != nil {
		log.Fatal(err)
	} else {
		t.Footer = tmpl
	}
	return t
}
