package adsig

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"git.bode.fun/adsig/config"
)

type Template struct {
	Name   string
	Fields map[string]string
	Files  []string
}

func (t Template) ParseFiles() ([]*template.Template, error) {
	tmpls := make([]*template.Template, 0)

	supportedExtensions := []string{
		".rtf",
		".txt",
		".html",
		".htm",
	}

	for _, fPath := range t.Files {
		templateName := filepath.Base(fPath)
		templateExt := filepath.Ext(fPath)

		isSupportedExt := false

		for _, ext := range supportedExtensions {
			if templateExt == ext {
				isSupportedExt = true
			}
		}

		if !isSupportedExt {
			continue
		}

		tmpl, err := template.New(templateName).Delims("[[", "]]").ParseFiles(fPath)
		if err != nil {
			return nil, err
		}

		tmpls = append(tmpls, tmpl)
	}

	return tmpls, nil
}

func filterTemplatesByName(src []Template, names []string) []Template {
	templates := make([]Template, 0)

	for _, name := range names {
		for _, tmpl := range src {
			if name == tmpl.Name {
				templates = append(templates, tmpl)
			}
		}
	}

	return templates
}

func templatesFromConfig(cnf config.Config) ([]Template, error) {
	// Get templates
	templatesDir, err := getTemplatesFolder()
	if err != nil {
		return nil, err
	}

	templates := make([]Template, 0)

	for cnfTmplName, cnfTmpl := range cnf.Templates {
		tmpl := Template{
			Name:   cnfTmplName,
			Fields: cnfTmpl.Fields,
			Files:  make([]string, 0),
		}

		tmplFiles, err := getFilesForTemplate(templatesDir, tmpl.Name)
		if err != nil {
			return nil, err
		}

		tmpl.Files = tmplFiles

		templates = append(templates, tmpl)
	}

	return templates, nil
}

func getFilesForTemplate(templatesDir, signatureName string) ([]string, error) {
	filePaths := make([]string, 0)

	signatureExtensions := []string{
		".html",
		".rtf",
		".txt",
	}

	for _, ext := range signatureExtensions {
		signaturePath := filepath.Join(templatesDir, signatureName+ext)
		if _, err := os.Stat(signaturePath); err != nil {
			return nil, err
		}

		filePaths = append(filePaths, signaturePath)
	}

	return filePaths, nil
}

func getTemplatesFolder() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	templateDir := filepath.Join(wd, "templates")

	fileInfo, err := os.Stat(templateDir)
	if err != nil || !fileInfo.IsDir() {
		// FIXME: Remove dynamic error
		return "", errors.New("main: templates folder is not present or can not be opened")
	}

	return templateDir, nil
}
