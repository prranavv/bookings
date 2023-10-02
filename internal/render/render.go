package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/models"
)

var app *config.AppConfig
var pathtotemplate = "./templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders template using html template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	//get the template cache form the app config
	//get requested template from the cache
	t, ok := tc[tmpl]
	if !ok {
		log.Println(2)

		//log.Fatal("Could not get template from the template cache")
		return errors.New("cant get template from cache")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//mycache := make(map[string]*template.Template)
	mycache := map[string]*template.Template{}

	//get all of the files name
	pages, err := filepath.Glob(fmt.Sprintf("%s/*page.html", pathtotemplate))
	if err != nil {
		return mycache, err
	}
	//range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return mycache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathtotemplate))
		if err != nil {
			return mycache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathtotemplate))
			if err != nil {
				return mycache, err
			}
		}
		mycache[name] = ts
	}
	return mycache, nil
}