package adsig

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"git.bode.fun/adsig/config"
)

type Signature struct {
	Name   string
	Fields map[string]string
	Files  []string
}

func (s Signature) ParseFiles() ([]*template.Template, error) {
	tmpls := make([]*template.Template, 0)

	supportedExtensions := []string{
		".rtf",
		".txt",
		".html",
		".htm",
	}

	for _, fPath := range s.Files {
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

func filterSignaturesByName(src []Signature, names []string) []Signature {
	signatures := make([]Signature, 0)

	for _, name := range names {
		for _, sig := range src {
			if name == sig.Name {
				signatures = append(signatures, sig)
			}
		}
	}

	return signatures
}

func SignaturesFromConfig(cnf config.Config) ([]Signature, error) {
	// Get templates
	templatesDir, err := getTemplatesFolder()
	if err != nil {
		return nil, err
	}

	signatures := make([]Signature, 0)

	for cnfTmplName, cnfTmpl := range cnf.Templates {
		sig := Signature{
			Name:   cnfTmplName,
			Fields: cnfTmpl.Fields,
			Files:  make([]string, 0),
		}

		sigTemplateFiles, err := getFilesForSignature(templatesDir, sig.Name)
		if err != nil {
			return nil, err
		}

		sig.Files = sigTemplateFiles

		signatures = append(signatures, sig)
	}

	return signatures, nil
}

func getFilesForSignature(templatesDir, signatureName string) ([]string, error) {
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
