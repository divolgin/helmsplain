package values

import (
	"html/template"

	"github.com/Masterminds/sprig/v3"
)

// Helm does not expose its func map unfortunately, so this is copied over and replaced with dummy functions

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	// Add some extra functionality
	extra := template.FuncMap{
		"toToml":        dummyFunc,
		"toYaml":        dummyFunc,
		"fromYaml":      dummyFunc,
		"fromYamlArray": dummyFunc,
		"toJson":        dummyFunc,
		"fromJson":      dummyFunc,
		"fromJsonArray": dummyFunc,
		"include":       dummyFunc,
		"tpl":           dummyFunc,
		"required":      dummyFunc,
		"lookup":        dummyFunc,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func dummyFunc(...interface{}) string {
	return "this is not helm"
}
