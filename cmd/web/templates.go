package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "pages/*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// parse layouts
		ts, err = ts.ParseGlob(filepath.Join(dir, "layouts/*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// parse partials
		ts, err = ts.ParseGlob(filepath.Join(dir, "partials/*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
