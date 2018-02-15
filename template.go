package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/flosch/pongo2"
)

// Template represents the html template that will be rendered as a PDF.
type Template struct {
	RootDir string
	Index   *pongo2.Template
	Footer  *pongo2.Template
}

// BuildParams creates the params for wkhtmltopdf
func (t *Template) BuildParams(url string) []string {
	params := []string{
		fmt.Sprintf("%s/main", url),
		// "--disable-smart-shrinking",
	}

	if t.Footer != nil {
		params = append(params, "--footer-html", fmt.Sprintf("%s/footer", url))
	}
	params = append(params, "-")
	return params
}

// WritePDF executes wkhtmltopdf with the correct params and
// writes the output to the provided io.Writer
func (t *Template) WritePDF(baseURL string, w io.Writer) error {
	params := t.BuildParams(baseURL)
	cmd := exec.Command("wkhtmltopdf", params...)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errChan := make(chan error)
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		_, err := io.Copy(w, output)
		errChan <- err
	}()

	if err := <-errChan; err != nil {
		return err
	}
	return cmd.Wait()
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// NewTemplate creates and initialize a template from a path
func NewTemplate(root, path string) (*Template, error) {
	t := &Template{
		RootDir: filepath.Join(root, path),
	}

	indexPath := filepath.Join(t.RootDir, "index.html")
	if !fileExists(indexPath) {
		return nil, fmt.Errorf("template %s not found", indexPath)
	} else if tmpl, err := pongo2.FromFile(indexPath); err != nil {
		return nil, err
	} else {
		t.Index = tmpl
	}

	footerPath := filepath.Join(t.RootDir, "footer.html")
	if !fileExists(footerPath) {
		return t, nil
	} else if tmpl, err := pongo2.FromFile(footerPath); err != nil {
		return nil, err
	} else {
		t.Footer = tmpl
	}
	return t, nil
}
