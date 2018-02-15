package main

import (
	"fmt"
	"io"
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

func (t *Template) WritePDF(baseUrl string, w io.Writer) error {
	params := t.BuildParams(baseUrl)
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

func NewTemplate(root, path string) (*Template, error) {
	t := &Template{
		Url:     path,
		RootDir: filepath.Join(root, path),
	}

	indexPath := filepath.Join(t.RootDir, "index.html")
	if !fileExists(indexPath) {
		return nil, fmt.Errorf("template %s not found.", indexPath)
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
