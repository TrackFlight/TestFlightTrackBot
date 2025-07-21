package internal

import (
	"path/filepath"
	"slices"
	"text/template"
)

type Templates struct {
	DB    *template.Template
	Model *template.Template
	Types *template.Template
	Store *template.Template
}

func ParseTemplates(schema Schema, templateDir string) *Templates {
	functions := template.FuncMap{
		"ToPascalCase":     ToPascalCase,
		"ToCamelCase":      ToCamelCase,
		"ToGoType":         ToGoType,
		"ToGoCase":         ToGoCase,
		"Singular":         Singular,
		"GetQueryOptions":  GetQueryOptions,
		"GetSprintfFormat": GetSprintfFormat,
		"IsBulkQuery":      IsBulkQuery,
		"GetParamsOrdered": GetParamsOrdered,
		"contains": func(s []string, e string) bool {
			return slices.Contains(s, e)
		},
		"FindCompatible": FindCompatible(schema.Tables),
	}
	tDB := parseTemplate("db.tpl", templateDir, functions)
	tModel := parseTemplate("model.tpl", templateDir, functions)
	tTypes := parseTemplate("types.tpl", templateDir, functions)
	tStore := parseTemplate("store.tpl", templateDir, functions)

	return &Templates{
		DB:    tDB,
		Model: tModel,
		Types: tTypes,
		Store: tStore,
	}
}

func parseTemplate(name, dir string, functions template.FuncMap) *template.Template {
	tpl, err := template.New(name).Funcs(functions).ParseFiles(filepath.Join(dir, name))
	checkErr(err)
	return tpl
}
