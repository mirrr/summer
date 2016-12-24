package main

import (
	"os"
	"path/filepath"
	"text/template"
)

type (
	obj map[string]interface{}
	arr []interface{}
)

func writeFiles(path string, filenames []string, t *template.Template, data arr) error {
	for i, filename := range filenames {
		if err := writeFile(path+filename, t, filename, data[i]); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filename string, t *template.Template, tname string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := t.ExecuteTemplate(fo, tname, data); err != nil {
		return err
	}
	fo.Close()
	return nil
}

func stripSlashes(s string) string {
	if len(s) > 0 {
		if s[len(s)-1] == '/' {
			s = s[:len(s)-1]
		}
	}
	if len(s) > 0 {
		if s[0] == '/' {
			s = s[1:]
		}
	}
	return s
}
