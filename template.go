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
	Header  *pongo2.Template
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
	if t.Header != nil {
		params = append(params, "--header-html", fmt.Sprintf("%s/header", url))
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

func loadFileIfExists(root, htmlPth string) (*pongo2.Template, error) {
	jpth := filepath.Join(root, htmlPth)
	if _, err := os.Stat(jpth); os.IsNotExist(err) {
		return nil, nil
	}
	return pongo2.FromFile(jpth)
}

// NewTemplate creates and initialize a template from a path
func NewTemplate(root, path string) (*Template, error) {
	root = filepath.Join(root, path)
	t := &Template{RootDir: root}
	var err error

	if t.Index, err = loadFileIfExists(root, "index.html"); err != nil {
		return nil, err
	} else if t.Index == nil {
		return nil, fmt.Errorf("file 'index.html' not found in directory %s", root)
	} else if t.Header, err = loadFileIfExists(root, "header.html"); err != nil {
		return nil, err
	} else if t.Footer, err = loadFileIfExists(root, "footer.html"); err != nil {
		return nil, err
	}
	return t, nil
}
