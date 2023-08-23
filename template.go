package adsig

import (
	"errors"
	"os"
	"path/filepath"

	"git.bode.fun/adsig/config"
)

type Template struct {
	Name   string
	Fields map[string]string
	Files  []string
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
	// Add template files to template
	signatureDir := filepath.Join(templatesDir, signatureName)

	if signatureDirInfo, err := os.Stat(signatureDir); err != nil || !signatureDirInfo.IsDir() {
		// FIXME: Remove dynamic error
		return nil, errors.New("main: template folder is not present or can not be opened")
	}

	filePaths := make([]string, 0)

	signatureFileNames := []string{
		"signature.html",
		"signature.rtf",
		"signature.txt",
	}

	for _, fileName := range signatureFileNames {
		signaturePath := filepath.Join(signatureDir, fileName)
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
